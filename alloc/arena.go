package alloc

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

// Arena allocates objects in internal byte blocks.
// It is goroutine-safe, but operations on a free arena panic.
type Arena struct {
	heap *heap

	spinlock int32
	free     bool
	size     int // total allocated size

	block  *block   // current block
	blocks []*block // all blocks
}

// NewArena returns a new arena with a global heap.
func NewArena() *Arena {
	return newArenaHeap(globalHeap)
}

// newArena returns a new arena with a new heap, for tests only.
func newArena() *Arena {
	return &Arena{heap: newHeap()}
}

// newArenaHeap returns a new arena with a given heap.
func newArenaHeap(heap *heap) *Arena {
	return &Arena{heap: heap}
}

// Size returns the total arena memory size in bytes.
func (a *Arena) Size() int64 {
	a.lock()
	defer a.unlock()

	return int64(a.size)
}

// Used returns the allocated arena memory size in bytes.
func (a *Arena) Used() int64 {
	a.lock()
	defer a.unlock()

	total := int64(0)
	for _, block := range a.blocks {
		total += int64(block.allocated())
	}
	return total
}

// Free frees the arena and releases its memory.
func (a *Arena) Free() {
	a.lock()
	defer a.unlock()

	if a.free {
		return
	}

	a.free = true
	a.size = 0

	blocks := a.blocks
	a.block = nil
	a.blocks = nil

	a.heap.freeBlocks(blocks...)
}

// internal

// alloc allocates data and returns a pointer to it.
func (a *Arena) alloc(size int) unsafe.Pointer {
	a.lock()
	defer a.unlock()

	if a.free {
		panic("allocation in a free arena")
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
	a.size += a.block.size()

	return a.block.alloc(size)
}

// private

func (a *Arena) lock() {
	for {
		if atomic.CompareAndSwapInt32(&a.spinlock, 0, 1) {
			return
		}
		runtime.Gosched()
	}
}

func (a *Arena) unlock() {
	atomic.StoreInt32(&a.spinlock, 0)
}

func (a *Arena) loadBlock() *block {
	ptr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&a.block)))
	if ptr == nil {
		return nil
	}

	return (*block)(ptr)
}
