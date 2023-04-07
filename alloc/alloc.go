package alloc

import (
	"github.com/complex1tech/baselibrary/alloc/arena"
	"github.com/complex1tech/baselibrary/alloc/internal/buf"
	"github.com/complex1tech/baselibrary/alloc/internal/heap"
	"github.com/complex1tech/baselibrary/ref"
)

// Allocator allocates arenas and buffers.
type Allocator interface {
	// Arena allocates a new arena.
	Arena() arena.Arena

	// Buffer allocates a new buffer.
	Buffer(size int) *Buffer
}

// Buffer is a buffer allocated by an allocator.
// The buffer must be freed after usage.
type Buffer = buf.Buffer

// New returns a new allocator.
func New() Allocator {
	h := heap.New()
	return newAllocator(h)
}

// Global returns the global allocator.
func Global() Allocator {
	return global
}

// Alloc

// NewArena allocates an arena in the global allocator.
func NewArena() arena.Arena {
	return global.Arena()
}

// NewArenaRef allocates an arena in the global allocator and returns a reference to it.
func NewArenaRef() *ref.R[arena.Arena] {
	return ref.Wrap(NewArena())
}

// NewBuffer allocates a buffer in the global allocator.
func NewBuffer() *Buffer {
	return global.Buffer(0)
}

// NewBuffer allocates a buffer of a preallocated capacity in the global allocator.
func NewBufferSize(size int) *Buffer {
	return global.Buffer(size)
}

// internal

var global = newAllocator(heap.Global())

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

// Buffer allocates a new buffer with a preallocated memory storage.
func (a *allocator) Buffer(size int) *Buffer {
	return buf.NewSize(a.heap, size)
}
