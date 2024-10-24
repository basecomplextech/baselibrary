// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"fmt"
	"sync/atomic"
)

// NewNoop returns a new reference with no freer.
func NewNoop[T any](obj T) R[T] {
	return &refNoop[T]{
		refs: 1,
		obj:  obj,
	}
}

// internal

var _ R[any] = (*refNoop[any])(nil)

type refNoop[T any] struct {
	refs int64
	obj  T
}

func (r *refNoop[T]) Refcount() int64 {
	return atomic.LoadInt64(&r.refs)
}

func (r *refNoop[T]) Retain() {
	v := atomic.AddInt64(&r.refs, 1)
	if v <= 1 {
		panic(fmt.Sprintf("retain: %T already released", r.obj))
	}
}

func (r *refNoop[T]) Release() {
	v := atomic.AddInt64(&r.refs, -1)
	switch {
	case v > 0:
		return
	case v < 0:
		panic(fmt.Sprintf("release: %T already released", r.obj))
	}
}

func (r *refNoop[T]) Unwrap() T {
	v := atomic.LoadInt64(&r.refs)
	if v <= 0 {
		panic(fmt.Sprintf("unwrap: %T already released", r.obj))
	}
	return r.obj
}
