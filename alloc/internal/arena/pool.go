// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package arena

import (
	"sync/atomic"
	"unsafe"
)

// Pool is a pool of objects allocated in the arena.
// It is thread-safe but only if backed by [MutexArena].
// The pool itself is allocated in the arena.
type Pool[T any] interface {
	// Get acquires an object and returns true, or allocates a new one and returns false.
	Get() (*T, bool)

	// Put puts an object back into the pool, does not zero it.
	// The object must be allocated in this pool.
	Put(obj *T)
}

// NewPool returns a new pool which allocates objects in the given arena.
func NewPool[T any](a Arena) Pool[T] {
	return newPool[T](a)
}

// internal

const (
	poolGetAttempts = 3
	poolPutAttempts = 3
)

type pool[T any] struct {
	arena  Arena
	head   atomic.Pointer[poolItem[T]] // head free item
	offset uintptr                     // object offset in poolItem
}

type poolItem[T any] struct {
	next atomic.Pointer[poolItem[T]] // next free item
	obj  T
}

func newPool[T any](arena Arena) *pool[T] {
	var m poolItem[T]

	p := Alloc[pool[T]](arena)
	p.arena = arena
	p.offset = unsafe.Offsetof(m.obj)
	return p
}

// Get acquires an object and returns true, or allocates a new one and returns false.
func (p *pool[T]) Get() (_ *T, ok bool) {
	for i := 0; i < poolGetAttempts; i++ {
		// Load current item
		item := p.head.Load()
		if item == nil {
			break
		}

		// Swap it with next
		next := item.next.Load()
		if !p.head.CompareAndSwap(item, next) {
			continue
		}

		// Return object
		item.next.Store(nil)
		return &item.obj, true
	}

	// Alloc new item
	m := Alloc[poolItem[T]](p.arena)
	return &m.obj, false
}

// Put puts an object back into the pool, does not zero it.
// The object must be allocated in this pool.
func (p *pool[T]) Put(obj *T) {
	// Compute item address
	ptr := unsafe.Pointer(obj)
	ptr1 := unsafe.Add(ptr, -p.offset)
	item := (*poolItem[T])(ptr1)

	for i := 0; i < poolPutAttempts; i++ {
		// Load current item
		head := p.head.Load()
		item.next.Store(head)

		// Swap it with next
		if p.head.CompareAndSwap(head, item) {
			return
		}
	}
}
