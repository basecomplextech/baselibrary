// Copyright 2021 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"fmt"
	"sync/atomic"
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
	return &ref[T]{
		refs: 1,
		obj:  obj,
	}
}

// internal

// maxUnpooledSize specifies the maximum object size before using a pooled reference.
const maxUnpooledSize = 24

var _ R[Freer] = (*ref[Freer])(nil)

type ref[T Freer] struct {
	refs int64
	obj  T
}

func (r *ref[T]) Refcount() int64 {
	return atomic.LoadInt64(&r.refs)
}

func (r *ref[T]) Retain() {
	v := atomic.AddInt64(&r.refs, 1)
	if v <= 1 {
		panic(fmt.Sprintf("retain: %T already released", r.obj))
	}
}

func (r *ref[T]) Release() {
	v := atomic.AddInt64(&r.refs, -1)
	switch {
	case v > 0:
		return
	case v < 0:
		panic(fmt.Sprintf("release: %T already released", r.obj))
	}

	var zero T
	r.obj.Free()
	r.obj = zero
}

func (r *ref[T]) Unwrap() T {
	v := atomic.LoadInt64(&r.refs)
	if v <= 0 {
		panic(fmt.Sprintf("unwrap: %T already released", r.obj))
	}
	return r.obj
}
