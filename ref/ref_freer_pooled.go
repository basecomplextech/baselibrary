// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the Business Source License (BSL 1.1)
// that can be found in the LICENSE file.

package ref

import (
	"fmt"

	"github.com/basecomplextech/baselibrary/pools"
)

var _ R[any] = (*refFreerPooled[any])(nil)

type refFreerPooled[T any] struct {
	refs Atomic64
	*refFreerState[T]
}

type refFreerState[T any] struct {
	pool  pools.Pool[*refFreerState[T]]
	freer Freer
	obj   T
}

func (r *refFreerPooled[T]) Refcount() int64 {
	return r.refs.Refcount()
}

func (r *refFreerPooled[T]) Acquire() bool {
	if ok := r.refs.Acquire(); ok {
		return true
	}

	r.Release()
	return false
}

func (r *refFreerPooled[T]) Retain() {
	if ok := r.refs.Acquire(); ok {
		return
	}

	r.Release()
	panic(fmt.Sprintf("retain: %T already released", r))
}

func (r *refFreerPooled[T]) Release() {
	released := r.refs.Release()
	if !released {
		return
	}

	r.freer.Free()

	s := r.refFreerState
	r.refFreerState = nil
	releaseRefFreer(s)
}

func (r *refFreerPooled[T]) Unwrap() T {
	v := r.refs.Refcount()
	if v <= 0 {
		panic(fmt.Sprintf("unwrap: %T already released", r))
	}
	return r.obj
}

// pool

var refFreerPools = pools.NewPools()

func acquireRefFreer[T any]() *refFreerState[T] {
	s, ok, pool := pools.Acquire1[*refFreerState[T]](refFreerPools)
	if ok {
		return s
	}

	s = &refFreerState[T]{}
	s.pool = pool
	return s
}

func releaseRefFreer[T any](s *refFreerState[T]) {
	p := s.pool
	*s = refFreerState[T]{}
	s.pool = p

	p.Put(s)
}
