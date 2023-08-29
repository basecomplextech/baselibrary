package alloc

import (
	"github.com/basecomplextech/baselibrary/alloc/internal/arena"
	"github.com/basecomplextech/baselibrary/ref"
)

// Arena is an arena memory allocator.
// The arena must be freed after usage.
type Arena = arena.Arena

// NewArena allocates an arena.
func NewArena() Arena {
	return arena.New(globalHeap)
}

// NewArenaSize allocates an arena of a preallocated capacity.
func NewArenaSize(size int) Arena {
	return arena.NewSize(globalHeap, size)
}

// NewArenaRef allocates an arena and returns a reference to it.
func NewArenaRef() *ref.R[Arena] {
	return ref.New(NewArena())
}
