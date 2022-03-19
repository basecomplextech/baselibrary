package alloc

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

// Arena allocates objects in internal byte blocks.
// Arena is goroutine-safe, however, accessing a freed arena panics.
type Arena interface {
	// Alloc allocates raw data and returns a pointer to it.
	Alloc(size int) unsafe.Pointer

	// Free frees the arena and releases its memory.
	Free()

	// Allocated returns returns the number of allocated bytes in the arena.
	Allocated() int64

	// Size returns the total arena size in bytes.
	Size() int64
}

type ArenaStats struct {
}

// NewArena returns a new arena with a global heap.
func NewArena() Arena {
	return newArenaHeap(globalHeap)
}

// internal

type arena struct {
	heap *heap

	spinlock int32
	freed    bool
	total    int // total allocated size

	block  *block   // current block
	blocks []*block // all blocks
}

// newArena returns a new arena with a new heap, for tests only.
func newArena() *arena {
	return &arena{heap: newHeap()}
}

// newArenaHeap returns a new arena with a given heap.
func newArenaHeap(heap *heap) *arena {
	return &arena{heap: heap}
}

// Alloc allocates raw data and returns a pointer to it.
func (a *arena) Alloc(size int) unsafe.Pointer {
	a.lock()
	defer a.unlock()

	if a.freed {
		panic("allocation in a freed arena")
	}

	if a.block != nil {
		free := a.block.free()
		if free >= size {
			return a.block.alloc(size)
		}
	}

	// double last block size
	// limit it to maxBlockSize
	blockSize := 0
	if a.block != nil {
		blockSize = a.block.size() * 2
	}
	if blockSize > maxBlockSize {
		blockSize = maxBlockSize
	}
	if size > blockSize {
		blockSize = size
	}

	a.block, _ = a.heap.allocBlock(blockSize)
	a.blocks = append(a.blocks, a.block)
	a.total += a.block.size()

	return a.block.alloc(size)
}

// Free frees the arena and releases its memory.
func (a *arena) Free() {
	a.lock()
	defer a.unlock()

	if a.freed {
		return
	}

	a.freed = true
	a.total = 0

	blocks := a.blocks
	a.block = nil
	a.blocks = nil

	a.heap.freeBlocks(blocks...)
}

// Allocated returns returns the number of allocated bytes in the arena.
func (a *arena) Allocated() int64 {
	a.lock()
	defer a.unlock()

	total := int64(0)
	for _, block := range a.blocks {
		total += int64(block.allocated())
	}
	return total
}

// Size returns arena size in bytes.
func (a *arena) Size() int64 {
	a.lock()
	defer a.unlock()

	return int64(a.total)
}

// private

func (a *arena) lock() {
	for {
		if atomic.CompareAndSwapInt32(&a.spinlock, 0, 1) {
			return
		}
		runtime.Gosched()
	}
}

func (a *arena) unlock() {
	atomic.StoreInt32(&a.spinlock, 0)
}
