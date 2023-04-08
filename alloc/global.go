package alloc

import (
	"github.com/complex1tech/baselibrary/alloc/arena"
	"github.com/complex1tech/baselibrary/alloc/internal/heap"
	"github.com/complex1tech/baselibrary/ref"
)

var global = newAllocator(heap.Global())

// Global returns the global allocator.
func Global() Allocator {
	return global
}

// NewArena allocates an arena in the global allocator.
func NewArena() arena.Arena {
	return global.Arena()
}

// NewArenaSize allocates an arena of a preallocated capacity in the global allocator.
func NewArenaSize(size int) arena.Arena {
	return global.ArenaSize(size)
}

// NewArenaRef allocates an arena in the global allocator and returns a reference to it.
func NewArenaRef() *ref.R[arena.Arena] {
	return ref.Wrap(NewArena())
}

// NewBuffer allocates a buffer in the global allocator.
func NewBuffer() *Buffer {
	return global.Buffer()
}

// NewBuffer allocates a buffer of a preallocated capacity in the global allocator.
func NewBufferSize(size int) *Buffer {
	return global.BufferSize(size)
}
