package freelist

import (
	"sync/atomic"
	"unsafe"

	"github.com/basecomplextech/baselibrary/alloc/internal/arena"
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

// New returns a free list, which allocs objects in the given arena.
func New[T any](arena arena.Arena) FreeList[T] {
	return newList[T](arena)
}

// internal

var _ FreeList[struct{}] = (*list[struct{}])(nil)

// list keeps a linked list of free objects in the arena.
type list[T any] struct {
	arena arena.Arena
	size  uintptr
	free  uintptr // last free item
}

type item struct {
	next uintptr // next free item
}

func newList[T any](a arena.Arena) *list[T] {
	var zero T
	size := unsafe.Sizeof(zero)

	// Increase size to hold item
	itemSize := unsafe.Sizeof(item{})
	if size < itemSize {
		size = itemSize
	}

	return &list[T]{
		arena: a,
		size:  size,
	}
}

// Get returns a free object from the list or allocates a new one in the arena.
func (l *list[T]) Get() *T {
	var zero T

	for i := 0; i < maxListGet; i++ {
		// Load current item
		free := atomic.LoadUintptr(&l.free)
		if free == 0 {
			break
		}
		uptr := *(*unsafe.Pointer)(unsafe.Pointer(&free))

		// Swap it with previous
		item := (*item)(uptr)
		if !atomic.CompareAndSwapUintptr(&l.free, free, item.next) {
			continue
		}

		// Reset and return object
		result := (*T)(uptr)
		*result = zero
		return result
	}

	// Allocate new object
	// TODO: Arena is not thread-safe any more, do we need to do anything?
	return arena.Alloc[T](l.arena)
}

// Put puts an object back into the free list.
func (l *list[T]) Put(ptr *T) {
	// Reset object
	var zero T
	*ptr = zero

	// Cast it into item
	item := (*item)(unsafe.Pointer(ptr))

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
