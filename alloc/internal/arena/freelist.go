// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the Business Source License (BSL 1.1)
// that can be found in the LICENSE file.

package arena

import (
	"unsafe"
)

// FreeList is a linked list of free objects allocated in the arena.
// The list itself is allocated in the arena.
type FreeList[T any] interface {
	// Get returns a free object from the list, or allocates a new one.
	Get() *T

	// Put puts an object back into the free list.
	Put(obj *T)
}

// NewFreeList returns a new free list which allocates objects in the given arena.
func NewFreeList[T any](a Arena) FreeList[T] {
	return newFreeList[T](a)
}

// internal

type freeList[T any] struct {
	arena Arena
	size  uintptr // element size rounded up to freeItem size
	free  uintptr // last free item
}

type freeItem struct {
	next uintptr // next free item
}

func newFreeList[T any](arena Arena) *freeList[T] {
	var zero T
	size := unsafe.Sizeof(zero)

	// Increase size to hold item
	itemSize := unsafe.Sizeof(freeItem{})
	if size < itemSize {
		size = itemSize
	}

	// Make list
	l := Alloc[freeList[T]](arena)
	l.arena = arena
	l.size = size
	return l
}

// Get returns a free object from the list, or allocates a new one.
func (l *freeList[T]) Get() *T {
	var zero T

	// Allocate new element if empty
	free := l.free
	if free == 0 {
		return Alloc[T](l.arena)
	}

	// Get free item
	uptr := *(*unsafe.Pointer)(unsafe.Pointer(&free))

	// Swap it with previous
	item := (*freeItem)(uptr)
	l.free = item.next

	// Zero and return object
	result := (*T)(uptr)
	*result = zero
	return result
}

// Put puts an object back into the free list.
func (l *freeList[T]) Put(obj *T) {
	// Zero object
	var zero T
	*obj = zero

	// Cast it into item
	item := (*freeItem)(unsafe.Pointer(obj))
	item.next = l.free

	// Swap current item with next
	next := (uintptr)(unsafe.Pointer(obj))
	l.free = next
}
