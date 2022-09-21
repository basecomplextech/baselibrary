package alloc

import (
	"sync/atomic"
	"unsafe"
)

// ArenaFreeList keeps a linked list of free objects in the arena.
type ArenaFreeList[T any] interface {
	// Get returns a free object from the list, or allocates a new one in the arena.
	Get() *T

	// Put puts an object back into the free list.
	Put(obj *T)
}

// NewArenaFreeList returns a new free list which allocates objects in the arena.
func NewArenaFreeList[T any](a Arena) ArenaFreeList[T] {
	arena := a.(*arena)
	return newArenaFreeList[T](arena)
}

// internal

const (
	arenaFreeListGetAttempts = 5
	arenaFreeListPutAttempts = 5
)

// arenaFreeList keeps a linked list of free objects in the arena.
type arenaFreeList[T any] struct {
	arena *arena
	size  uintptr
	free  uintptr // last free item
}

type arenaFreeListItem struct {
	next uintptr // next free item
}

func newArenaFreeList[T any](arena *arena) *arenaFreeList[T] {
	var zero T
	size := unsafe.Sizeof(zero)

	// increase size to hold item
	itemSize := unsafe.Sizeof(arenaFreeListItem{})
	if size < itemSize {
		size = itemSize
	}

	return &arenaFreeList[T]{
		arena: arena,
		size:  size,
	}
}

// Get returns a free object from the list or allocates a new one in the arena.
func (l *arenaFreeList[T]) Get() *T {
	var zero T

	for i := 0; i < arenaFreeListGetAttempts; i++ {
		// load current item
		free := atomic.LoadUintptr(&l.free)
		if free == 0 {
			break
		}
		uptr := *(*unsafe.Pointer)(unsafe.Pointer(&free))

		// swap it with previous
		item := (*arenaFreeListItem)(uptr)
		if !atomic.CompareAndSwapUintptr(&l.free, free, item.next) {
			continue
		}

		// reset and return object
		result := (*T)(uptr)
		*result = zero
		return result
	}

	// allocate new object
	return ArenaAlloc[T](l.arena)
}

// Put puts an object back into the free list.
func (l *arenaFreeList[T]) Put(ptr *T) {
	// reset object
	var zero T
	*ptr = zero

	// cast it into item
	item := (*arenaFreeListItem)(unsafe.Pointer(ptr))

	for i := 0; i < arenaFreeListPutAttempts; i++ {
		// load current item
		free := atomic.LoadUintptr(&l.free)
		item.next = free

		// swap it with next
		next := (uintptr)(unsafe.Pointer(ptr))
		if atomic.CompareAndSwapUintptr(&l.free, free, next) {
			return
		}
	}
}
