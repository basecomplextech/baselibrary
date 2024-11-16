// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"fmt"
)

var _ R[any] = (*varRef[any])(nil)

type varRef[T any] struct {
	refs Atomic64
	ref  R[T]
}

func newVarRef[T any](ref R[T]) *varRef[T] {
	r := &varRef[T]{}
	r.refs.Init(1)
	r.ref = ref
	return r
}

// Refcount returns the number of current references.
func (r *varRef[T]) Refcount() int64 {
	return r.refs.Refcount()
}

// Retain increments refcount, panics when count is <= 0.
func (r *varRef[T]) Retain() {
	// Increment refs, return if alive
	acquired := r.refs.Acquire()
	if acquired {
		return
	}

	// Otherwise, release immediately
	r.release()

	// Panic on retain of released object
	var zero T
	panic(fmt.Sprintf("retain: %T already released", zero))
}

// Release decrements refcount and releases the object if the count is 0.
func (r *varRef[T]) Release() {
	r.release()
}

// Unwrap returns the object or panics if the refcount is 0.
func (r *varRef[T]) Unwrap() (v T) {
	if r.refs.Refcount() > 0 {
		return r.ref.Unwrap()
	}

	var zero T
	panic(fmt.Sprintf("unwrap: %T already released", zero))
}

// private

// acquire increments refcount and returns true if the released bit is not set.
// otherwise, it releases the reference and returns false.
func (r *varRef[T]) acquire() bool {
	// Increment refs, return if alive
	acquired := r.refs.Acquire()
	if acquired {
		return true
	}

	// Otherwise, release immediately
	r.release()
	return false
}

// release decrements refcount, and frees the object if refcount is 0.
func (r *varRef[T]) release() {
	// Decrement refs, return if alive
	released := r.refs.Release()
	if !released {
		return
	}

	// Free value when released
	r.free()
}

// free frees the reference, called only once.
func (r *varRef[T]) free() {
	ref := r.ref
	r.ref = nil

	ref.Release()
}
