package alloc

import "github.com/basecomplextech/baselibrary/alloc/internal/freelist"

// FreeList keeps a linked list of free objects.
type FreeList[T any] interface {
	// Get returns a free object from the list, or allocates a new one.
	Get() *T

	// Put puts an object back into the free list.
	Put(obj *T)
}

// NewFreeList returns a new free list which allocates objects in the given arena.
func NewFreeList[T any](a Arena) FreeList[T] {
	return freelist.New[T](a)
}
