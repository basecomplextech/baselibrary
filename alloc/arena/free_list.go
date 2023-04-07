package arena

import (
	"sync/atomic"
	"unsafe"
)

const (
	maxListGet = 5
	maxListPut = 5
)

// FreeList keeps a linked list of free objects.
type FreeList[T any] interface {
	// Get returns a free object from the list, or allocates a new one.
	Get() *T

	// Put puts an object back into the free list.
	Put(obj *T)
}

// NewFreeList returns a new free list which allocates objects in the given arena.
func NewFreeList[T any](a Arena) FreeList[T] {
	arena := a.(*arena)
	return newFreeList[T](arena)
}

// internal

var _ FreeList[struct{}] = (*freeList[struct{}])(nil)

// freeList keeps a linked list of free objects in the arena.
type freeList[T any] struct {
	arena *arena
	size  uintptr
	free  uintptr // last free item
}

type freeListItem struct {
	next uintptr // next free item
}

func newFreeList[T any](arena *arena) *freeList[T] {
	var zero T
	size := unsafe.Sizeof(zero)

	// Increase size to hold item
	itemSize := unsafe.Sizeof(freeListItem{})
	if size < itemSize {
		size = itemSize
	}

	return &freeList[T]{
		arena: arena,
		size:  size,
	}
}

// Get returns a free object from the list or allocates a new one in the arena.
func (l *freeList[T]) Get() *T {
	var zero T

	for i := 0; i < maxListGet; i++ {
		// Load current item
		free := atomic.LoadUintptr(&l.free)
		if free == 0 {
			break
		}
		uptr := *(*unsafe.Pointer)(unsafe.Pointer(&free))

		// Swap it with previous
		item := (*freeListItem)(uptr)
		if !atomic.CompareAndSwapUintptr(&l.free, free, item.next) {
			continue
		}

		// Reset and return object
		result := (*T)(uptr)
		*result = zero
		return result
	}

	// Allocate new object
	return Alloc[T](l.arena)
}

// Put puts an object back into the free list.
func (l *freeList[T]) Put(ptr *T) {
	// Reset object
	var zero T
	*ptr = zero

	// Cast it into item
	item := (*freeListItem)(unsafe.Pointer(ptr))

	for i := 0; i < maxListPut; i++ {
		// Load current item
		free := atomic.LoadUintptr(&l.free)
		item.next = free

		// Swap it with next
		next := (uintptr)(unsafe.Pointer(ptr))
		if atomic.CompareAndSwapUintptr(&l.free, free, next) {
			return
		}
	}
}
