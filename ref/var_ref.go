// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"fmt"
	"sync/atomic"
)

const varDetachedBit = int64(1 << 62)

var _ R[any] = (*varRef[any])(nil)

type varRef[T any] struct {
	refs atomic.Int64
	ref  atomic.Pointer[R[T]]
	_ref R[T]
}

func newVarRef[T any](ref R[T]) *varRef[T] {
	r := &varRef[T]{}
	r._ref = ref
	r.refs.Store(1)
	r.ref.Store(&r._ref)
	return r
}

// Refcount returns the number of current references.
func (r *varRef[T]) Refcount() int64 {
	n := r.refs.Load()
	n &= ^varDetachedBit // clear detached bit
	return int64(n)
}

// Retain increments refcount, panics when count is <= 0.
func (r *varRef[T]) Retain() {
	n := r.refs.Add(1)
	n &= ^varDetachedBit // clear detached bit

	if n <= 1 {
		var zero T
		panic(fmt.Sprintf("retain: %T already released", zero))
	}
}

// Release decrements refcount and releases the object if the count is 0.
func (r *varRef[T]) Release() {
	r.release()
}

// Unwrap returns the object or panics if the refcount is 0.
func (r *varRef[T]) Unwrap() (v T) {
	ref := r.ref.Load()
	if ref == nil {
		return v
	}
	return (*ref).Unwrap()
}

// private

// acquire increments refcount and returns true if the detached bit is not set.
// otherwise, it releases the reference and returns false.
func (r *varRef[T]) acquire() bool {
	// Increment refcount
	n := r.refs.Add(1)
	detached := n&varDetachedBit != 0 // get detached bit
	n &= ^varDetachedBit              // clear detached bit

	// Return if not detached
	if !detached {
		return true
	}

	// Release immediately otherwise
	r.release()
	return false
}

// release decrements refcount, and frees the object if refcount is 0.
func (r *varRef[T]) release() {
	n := r.refs.Add(-1)
	n &= ^varDetachedBit // clear detached bit

	// Panic if negative refcount
	switch {
	case n > 0:
		return
	case n < 0:
		var zero T
		panic(fmt.Sprintf("release: %T already released", zero))
	}

	// Free value when refcount is 0
	r.free()
}

// detach sets the detached bit and releases the reference.
func (r *varRef[T]) detach() {
	// It is still possible that some readers increment the refcount when it is 0.
	// Yet, they will see the detached bit and release the reference immediatelly.

	r.refs.Or(varDetachedBit)
	r.release()
}

// free frees the reference.
func (r *varRef[T]) free() {
	// The method may be called multiple times because refcount may reach 0 multiple times.
	// For example, when readers fail to acquire a detached reference.

	ptr := r.ref.Swap(nil)
	if ptr == nil {
		return
	}

	(*ptr).Release()
}
