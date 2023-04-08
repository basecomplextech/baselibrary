package alloc

import (
	"github.com/complex1tech/baselibrary/alloc/arena"
	"github.com/complex1tech/baselibrary/alloc/internal/buf"
	"github.com/complex1tech/baselibrary/alloc/internal/heap"
)

// Allocator allocates arenas and buffers.
type Allocator interface {
	// Arena allocates a new arena.
	Arena() arena.Arena

	// ArenaSize allocates a new arena with a preallocated capacity.
	ArenaSize(size int) arena.Arena

	// Buffer allocates a new buffer.
	Buffer() *Buffer

	// BufferSize allocates a new buffer with a preallocated capacity.
	BufferSize(size int) *Buffer
}

// Buffer is a buffer allocated by an allocator.
// The buffer must be freed after usage.
type Buffer = buf.Buffer

// New returns a new allocator.
func New() Allocator {
	h := heap.New()
	return newAllocator(h)
}

// internal

type allocator struct {
	heap *heap.Heap
}

func newAllocator(heap *heap.Heap) *allocator {
	return &allocator{heap: heap}
}

// Arena allocates a new arena.
func (a *allocator) Arena() arena.Arena {
	return arena.New(a.heap)
}

// ArenaSize allocates a new arena with a preallocated capacity.
func (a *allocator) ArenaSize(size int) arena.Arena {
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
