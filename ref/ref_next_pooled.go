// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the Business Source License (BSL 1.1)
// that can be found in the LICENSE file.

package ref

import (
	"fmt"

	"github.com/basecomplextech/baselibrary/pools"
)

var _ R[any] = (*refNextPooled[any, any])(nil)

type refNextPooled[T, T1 any] struct {
	refs Atomic64
	*refNextState[T, T1]
}

type refNextState[T, T1 any] struct {
	pool   pools.Pool[*refNextState[T, T1]]
	parent R[T1]
	obj    T
}

func (r *refNextPooled[T, T1]) Refcount() int64 {
	return r.refs.Refcount()
}

func (r *refNextPooled[T, T1]) Acquire() bool {
	if ok := r.refs.Acquire(); ok {
		return true
	}

	r.Release()
	return false
}

func (r *refNextPooled[T, T1]) Retain() {
	if ok := r.refs.Acquire(); ok {
		return
	}

	r.Release()
	panic(fmt.Sprintf("retain: %T already released", r))
}

func (r *refNextPooled[T, T1]) Release() {
	released := r.refs.Release()
	if !released {
		return
	}

	r.parent.Release()

	s := r.refNextState
	r.refNextState = nil
	releaseRefNext(s)
}

func (r *refNextPooled[T, T1]) Unwrap() T {
	v := r.refs.Refcount()
	if v <= 0 {
		panic(fmt.Sprintf("unwrap: %T already released", r))
	}
	return r.obj
}

// pool

var refNextPools = pools.NewPools()

func acquireRefNext[T, T1 any]() *refNextState[T, T1] {
	s, ok, pool := pools.Acquire1[*refNextState[T, T1]](refNextPools)
	if ok {
		return s
	}

	s = &refNextState[T, T1]{}
	s.pool = pool
	return s
}

func releaseRefNext[T, T1 any](s *refNextState[T, T1]) {
	p := s.pool
	*s = refNextState[T, T1]{}
	s.pool = p

	p.Put(s)
}
