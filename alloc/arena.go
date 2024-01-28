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

// Pin

// Pinned is a wrapper for an object pinned to an arena.
type Pinned[T any] struct {
	Obj T
	Set bool
}

// Pin pins an object to an arena.
func Pin[T any](arena Arena, obj T) Pinned[T] {
	arena.Pin(obj)
	return Pinned[T]{Obj: obj}
}

// Reset clears the pinned object.
func (p *Pinned[T]) Reset() {
	var zero T
	p.Obj = zero
	p.Set = false
}

// Unwrap returns the pinned object and panics if the object is not pinned.
func (p Pinned[T]) Unwrap() T {
	if !p.Set {
		panic("object not pinned")
	}
	return p.Obj
}
