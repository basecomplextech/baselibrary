// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"fmt"
)

// NewNoop returns a new reference with no freer.
func NewNoop[T any](obj T) R[T] {
	r := &refNoop[T]{obj: obj}
	r.refs.Init(1)
	return r
}

// internal

var _ R[any] = (*refNoop[any])(nil)

type refNoop[T any] struct {
	refs Atomic64
	obj  T
}

func (r *refNoop[T]) Refcount() int64 {
	return r.refs.Refcount()
}

func (r *refNoop[T]) Acquire() bool {
	if ok := r.refs.Acquire(); ok {
		return true
	}

	r.Release()
	return false
}

func (r *refNoop[T]) Retain() {
	if ok := r.refs.Acquire(); ok {
		return
	}

	panic(fmt.Sprintf("retain: %T already released", r.obj))
}

func (r *refNoop[T]) Release() {
	r.refs.Release()
}

func (r *refNoop[T]) Unwrap() T {
	v := r.refs.Refcount()
	if v <= 0 {
		panic(fmt.Sprintf("unwrap: %T already released", r.obj))
	}
	return r.obj
}
