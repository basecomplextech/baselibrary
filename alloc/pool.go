// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the Business Source License (BSL 1.1)
// that can be found in the LICENSE file.

package alloc

import "github.com/basecomplextech/baselibrary/alloc/internal/arena"

// Pool is a pool of objects allocated in the arena.
// It is thread-safe but only if backed by [MutexArena].
// The pool itself is allocated in the arena.
type Pool[T any] interface {
	// Get acquires an object and returns true, or allocates a new one and returns false.
	Get() (*T, bool)

	// Put puts an object back into the pool.
	// The object must be allocated in this pool.
	Put(obj *T)
}

// NewPool returns a new pool which allocates objects in the given arena.
func NewPool[T any](a Arena) Pool[T] {
	return arena.NewPool[T](a)
}
