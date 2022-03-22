package alloc

import (
	"reflect"
	"runtime"
	"sync/atomic"
	"unsafe"
)

const alignment = 8

// Arena allocates objects in internal byte blocks.
// It is goroutine-safe, but operations on a free arena panic.
type Arena struct {
	heap *heap

	spinlock int32
	free     bool
	size     int64 // total allocated size

	blocks []*block
	lists  map[reflect.Type]unsafe.Pointer
}

// NewArena returns a new arena with a global heap.
func NewArena() *Arena {
	return newArenaHeap(globalHeap)
}

// newArena returns a new arena with a new heap, for tests only.
func newArena() *Arena {
	return newArenaHeap(newHeap())
}

// newArenaHeap returns a new arena with a given heap.
func newArenaHeap(heap *heap) *Arena {
	return &Arena{
		heap:  heap,
		lists: make(map[reflect.Type]unsafe.Pointer),
	}
}

// Size returns the total arena memory size in bytes.
func (a *Arena) Size() int64 {
	a.lock()
	defer a.unlock()

	return a.size
}

// Used returns the allocated arena memory size in bytes.
func (a *Arena) Used() int64 {
	a.lock()
	defer a.unlock()

	total := int64(0)
	for _, block := range a.blocks {
		total += int64(block.len())
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
	a.blocks = nil
	a.lists = make(map[reflect.Type]unsafe.Pointer)

	a.heap.freeBlocks(blocks...)
}

// internal

// alloc allocates data and returns a pointer to it.
func (a *Arena) alloc(n int) unsafe.Pointer {
	a.lock()
	defer a.unlock()

	if a.free {
		panic("operation on a free arena")
	}

	// maybe allocate block
	last := a.last()
	if last == nil || last.free() < n {
		last = a.allocBlock(n)
	}

	// grow buffer, add end padding for alignement
	start := len(last.buf)
	end := start + n
	if end > last.cap() {
		panic("allocation out of block range") // unreachable
	}

	end += (alignment - (end % alignment)) % alignment
	if end > last.cap() {
		end = last.cap()
	}

	last.buf = last.buf[:end]

	// slice buffer
	p := last.buf[start:end:end] // start:end:max, cap=max-start
	ptr := unsafe.Pointer(&p[0])
	return ptr
}

// allocFreeList allocates a new free list or returns an existing one.
func allocFreeList[T any](a *Arena) *FreeList[T] {
	var zero T
	typ := reflect.TypeOf(zero)

	a.lock()
	defer a.unlock()

	if a.free {
		panic("operation on a free arena")
	}

	// get existing list
	uptr, ok := a.lists[typ]
	if ok {
		return (*FreeList[T])(uptr)
	}

	// init new list
	list := newFreeList[T](a)
	a.lists[typ] = (unsafe.Pointer)(list)
	return list
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

func (a *Arena) last() *block {
	if len(a.blocks) == 0 {
		return nil
	}
	return a.blocks[len(a.blocks)-1]
}

func (a *Arena) allocBlock(n int) *block {
	// double last block size
	// limit it to maxBlockSize
	size := 0
	last := a.last()

	if last != nil {
		size = last.cap() * 2
	}
	if size > maxBlockSize {
		size = maxBlockSize
	}
	if n > size {
		size = n
	}

	block, _ := a.heap.allocBlock(size)
	a.blocks = append(a.blocks, block)
	a.size += int64(block.cap())
	return block
}
