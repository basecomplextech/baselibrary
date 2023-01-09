package alloc

import (
	"sync/atomic"
	"unsafe"
)

// internal

const (
	arenaListGetAttempts = 5
	arenaListPutAttempts = 5
)

// arenaList keeps a linked list of free objects in the arena.
type arenaList[T any] struct {
	arena *arena
	size  uintptr
	free  uintptr // last free item
}

type arenaListItem struct {
	next uintptr // next free item
}

func newArenaList[T any](arena *arena) *arenaList[T] {
	var zero T
	size := unsafe.Sizeof(zero)

	// increase size to hold item
	itemSize := unsafe.Sizeof(arenaListItem{})
	if size < itemSize {
		size = itemSize
	}

	return &arenaList[T]{
		arena: arena,
		size:  size,
	}
}

// Get returns a free object from the list or allocates a new one in the arena.
func (l *arenaList[T]) Get() *T {
	var zero T

	for i := 0; i < arenaListGetAttempts; i++ {
		// load current item
		free := atomic.LoadUintptr(&l.free)
		if free == 0 {
			break
		}
		uptr := *(*unsafe.Pointer)(unsafe.Pointer(&free))

		// swap it with previous
		item := (*arenaListItem)(uptr)
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
func (l *arenaList[T]) Put(ptr *T) {
	// reset object
	var zero T
	*ptr = zero

	// cast it into item
	item := (*arenaListItem)(unsafe.Pointer(ptr))

	for i := 0; i < arenaListPutAttempts; i++ {
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
