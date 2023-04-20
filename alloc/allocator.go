package alloc

import (
	"github.com/complex1tech/baselibrary/alloc/internal/arena"
	"github.com/complex1tech/baselibrary/alloc/internal/buf"
	"github.com/complex1tech/baselibrary/alloc/internal/heap"
)

var global = newAllocator(heap.Global())

type allocator struct {
	heap *heap.Heap
}

func newAllocator(heap *heap.Heap) *allocator {
	return &allocator{heap: heap}
}

// Arena allocates a new arena.
func (a *allocator) Arena() Arena {
	return arena.New(a.heap)
}

// ArenaSize allocates a new arena with a preallocated capacity.
func (a *allocator) ArenaSize(size int) Arena {
	return arena.NewSize(a.heap, size)
}

// Buffer allocates a new buffer.
func (a *allocator) Buffer() *Buffer {
	return buf.New(a.heap)
}

// BufferSize allocates a new buffer with a preallocated capacity.
func (a *allocator) BufferSize(size int) *Buffer {
	return buf.NewSize(a.heap, size)
}
