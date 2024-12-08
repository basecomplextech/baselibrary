// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"fmt"
)

var _ R[any] = (*shardedVarRef[any])(nil)

type shardedVarRef[T any] struct {
	refs *Atomic64
	set  *shardedVarRefset[T]
}

func (r *shardedVarRef[T]) init(refs *Atomic64, set *shardedVarRefset[T]) {
	r.set = set
	r.refs = refs
	r.refs.Init(1)
}

// Acquire increments refcount and returns the reference.
func (r *shardedVarRef[T]) Acquire() bool {
	if ok := r.refs.Acquire(); ok {
		return true
	}

	r.Release()
	return false
}

// Refcount returns the number of current references.
func (r *shardedVarRef[T]) Refcount() int64 {
	return r.refs.Refcount()
}

// Retain increments refcount, panics when count is <= 0.
func (r *shardedVarRef[T]) Retain() {
	if ok := r.refs.Acquire(); ok {
		return
	}
	r.refs.Release()

	var obj T
	panic(fmt.Sprintf("retain of already released %T reference", obj))
}

// Release decrements refcount and releases the object if the count is 0.
func (r *shardedVarRef[T]) Release() {
	if ok := r.refs.Release(); !ok {
		return
	}

	v := r.set
	r.set = nil
	v.release()
}

// Unwrap returns the object or panics if the refcount is 0.
func (r *shardedVarRef[T]) Unwrap() T {
	// Unwrapping a released reference is already a programming error.
	// So let's ignore a possible race condition here.

	count := r.Refcount()
	if count <= 0 {
		panic("unwrap of released reference")
	}

	return r.set.unwrap()
}
