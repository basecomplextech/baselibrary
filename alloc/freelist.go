// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package alloc

import "github.com/basecomplextech/baselibrary/alloc/internal/arena"

// FreeList is a linked list of free objects allocated in the arena.
// The list itself is allocated in the arena.
type FreeList[T any] interface {
	// Get returns a free object from the list, or allocates a new one.
	Get() *T

	// Put puts an object back into the free list.
	Put(obj *T)
}

// NewFreeList returns a new free list which allocates objects in the given arena.
// The list itself is allocated in the arena.
func NewFreeList[T any](a Arena) FreeList[T] {
	return arena.NewFreeList[T](a)
}
