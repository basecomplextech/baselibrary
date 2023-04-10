package arena

import (
	"runtime"
	"sync/atomic"
	"unsafe"

	"github.com/complex1tech/baselibrary/alloc/internal/heap"
	"github.com/complex1tech/baselibrary/collect/slices"
)

// Arena is an arena allocator, which internally allocates memory in blocks.
// It is goroutine-safe, but operations on a freed arena panic.
type Arena interface {
	// Cap returns the arena capacity.
	Cap() int64

	// Len calculates and returns the arena used size.
	Len() int64

	// Bytes allocates a byte slice.
	Bytes(cap int) []byte

	// Reset resets the arena.
	Reset()

	// Internal

	// Free frees the arena and releases its memory.
	Free()
}

// New returns an empty arena.
func New(h *heap.Heap) Arena {
	return newArena(h, heap.MinBlockSize)
}

// NewSize returns an arena with a preallocated memory storage.
func NewSize(heap *heap.Heap, size int) Arena {
	return newArena(heap, size)
}

// internal

type arena struct {
	heap    *heap.Heap
	initCap int // initial capacity

	spinlock int32
	free     bool
	cap      int64 // total allocated capacity

	blocks []*heap.Block
}

func newArena(heap *heap.Heap, size int) *arena {
	a := &arena{heap: heap}
	if size > 0 {
		a.initCap = a.allocBlock(size).Cap()
	}
	return a
}

// Cap returns the arena capacity.
func (a *arena) Cap() int64 {
	a.lock()
	defer a.unlock()

	return a.cap
}

// Len calculates and returns the arena used size.
func (a *arena) Len() int64 {
	a.lock()
	defer a.unlock()

	n := int64(0)
	for _, block := range a.blocks {
		n += int64(block.Len())
	}
	return n
}

// Bytes allocates a byte slice.
func (a *arena) Bytes(cap int) []byte {
	if cap == 0 {
		return nil
	}

	ptr := a.alloc(cap)
	return unsafe.Slice((*byte)(ptr), cap)
}

// Reset resets the arena.
func (a *arena) Reset() {
	a.lock()
	defer a.unlock()

	if len(a.blocks) == 0 {
		return
	}

	// Maybe just reset the first block
	n := 0
	if b := a.blocks[0]; b.Cap() == a.initCap {
		n = 1

		b.Reset()
		a.cap = int64(b.Cap())

		if len(a.blocks) == 1 {
			return
		}
	}

	// Free other blocks
	a.heap.FreeMany(a.blocks[n:]...)
	slices.Zero(a.blocks[n:]) // for gc

	a.cap = 0
	a.blocks = a.blocks[:n]
}

// Internal

// Free frees the arena and releases its memory.
func (a *arena) Free() {
	a.lock()
	defer a.unlock()

	if a.free {
		return
	}

	a.free = true
	a.cap = 0

	blocks := a.blocks
	a.blocks = nil
	a.heap.FreeMany(blocks...)
}

// alloc

// alloc allocates data and returns a pointer to it.
func (a *arena) alloc(n int) unsafe.Pointer {
	a.lock()
	defer a.unlock()

	if a.free {
		panic("operation on a free arena")
	}

	if len(a.blocks) > 0 {
		b := a.blocks[len(a.blocks)-1]
		ptr := b.Alloc(n)
		if ptr != nil {
			return ptr
		}
	}

	b := a.allocBlock(n)
	return b.Alloc(n)
}

func (a *arena) allocBlock(n int) *heap.Block {
	// Double last block capacity
	size := 0
	if len(a.blocks) > 0 {
		last := a.blocks[len(a.blocks)-1]
		size = last.Cap() * 2
	}
	if n > size {
		size = n
	}

	// Alloc next block
	b := a.heap.Alloc(size)
	a.blocks = append(a.blocks, b)
	a.cap += int64(b.Cap())
	return b
}

// lock

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
