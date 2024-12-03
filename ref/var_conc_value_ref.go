// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"fmt"
)

var _ R[any] = (*concValueRef[any])(nil)

type concValueRef[T any] struct {
	val  *concValue[T]
	refs *Atomic64
}

func (r *concValueRef[T]) init(val *concValue[T], refs *Atomic64) {
	r.val = val
	r.refs = refs
	r.refs.Init(1)
}

// Acquire increments refcount and returns the reference.
func (r *concValueRef[T]) Acquire() bool {
	if ok := r.refs.Acquire(); ok {
		return true
	}

	r.Release()
	return false
}

// Refcount returns the number of current references.
func (r *concValueRef[T]) Refcount() int64 {
	return r.refs.Refcount()
}

// Retain increments refcount, panics when count is <= 0.
func (r *concValueRef[T]) Retain() {
	if ok := r.refs.Acquire(); ok {
		return
	}
	r.refs.Release()

	var obj T
	panic(fmt.Sprintf("retain of already released %T reference", obj))
}

// Release decrements refcount and releases the object if the count is 0.
func (r *concValueRef[T]) Release() {
	if ok := r.refs.Release(); !ok {
		return
	}

	v := r.val
	r.val = nil
	v.release()
}

// Unwrap returns the object or panics if the refcount is 0.
func (r *concValueRef[T]) Unwrap() T {
	// Unwrapping a released reference is already a programming error.
	// So let's ignore the race condition here.

	count := r.Refcount()
	if count <= 0 {
		panic("unwrap of released reference")
	}

	return r.val.unwrap()
}
