// Copyright 2021 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"fmt"
)

// R is a generic atomic countable reference.
// It wraps an object and frees it when refcount reaches 0.
type R[T any] interface {
	// Refcount returns the number of current references.
	Refcount() int64

	// Retain increments refcount, panics when count is <= 0.
	Retain()

	// Release decrements refcount and releases the object if the count is 0.
	Release()

	// Unwrap returns the object or panics if the refcount is 0.
	Unwrap() T
}

// Ref is a countable reference interface without generics, i.e. R[?].
type Ref interface {
	// Retain increments refcount, panics when count is <= 0.
	Retain()

	// Release decrements refcount and releases the object if the count is 0.
	Release()
}

// New

// New returns a new reference with refcount 1.
//
// Constructors:
//   - [New] returns a new reference.
//   - [NewFree] returns a new reference with a free function.
//   - [NewFreer] returns a new reference with a custom freer.
//   - [NewNoop] returns a new reference without any freer.
//   - [Next] returns a child reference with a parent reference as a freer.
//   - [NextRetain] returns a child reference and retains the parent.
func New[T Freer](obj T) R[T] {
	r := &ref[T]{obj: obj}
	r.refs.Init(1)
	return r
}

// internal

// maxUnpooledSize specifies the maximum object size before using a pooled reference.
const maxUnpooledSize = 24

var _ R[Freer] = (*ref[Freer])(nil)

type ref[T Freer] struct {
	refs Atomic64
	obj  T
}

func (r *ref[T]) Refcount() int64 {
	return r.refs.Refcount()
}

func (r *ref[T]) Acquire() bool {
	if ok := r.refs.Acquire(); ok {
		return true
	}

	r.Release()
	return false
}

func (r *ref[T]) Retain() {
	if ok := r.refs.Acquire(); ok {
		return
	}

	r.Release()
	panic(fmt.Sprintf("retain: %T already released", r.obj))
}

func (r *ref[T]) Release() {
	released := r.refs.Release()
	if !released {
		return
	}

	var zero T
	r.obj.Free()
	r.obj = zero
}

func (r *ref[T]) Unwrap() T {
	v := r.refs.Refcount()
	if v <= 0 {
		panic(fmt.Sprintf("unwrap: %T already released", r.obj))
	}
	return r.obj
}
