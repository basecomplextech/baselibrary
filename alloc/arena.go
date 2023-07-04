package alloc

import (
	"github.com/complex1tech/baselibrary/alloc/internal/arena"
	"github.com/complex1tech/baselibrary/ref"
)

// Arena is an arena memory allocator.
// The arena must be freed after usage.
type Arena = arena.Arena

// NewArena allocates an arena in the global allocator.
func NewArena() Arena {
	return global.Arena()
}

// NewArenaSize allocates an arena of a preallocated capacity in the global allocator.
func NewArenaSize(size int) Arena {
	return global.ArenaSize(size)
}

// NewArenaRef allocates an arena in the global allocator and returns a reference to it.
func NewArenaRef() *ref.R[Arena] {
	return ref.New(NewArena())
}
