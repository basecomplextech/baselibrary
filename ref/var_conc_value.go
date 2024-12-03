// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"slices"

	"github.com/basecomplextech/baselibrary/collect/slices2"
	"github.com/basecomplextech/baselibrary/pools"
)

type concValue[T any] struct {
	refs   Atomic64
	ref    R[T]
	counts []concValueCount

	pool pools.Pool[*concValue[T]]
}

func newConcValue[T any](pool pools.Pool[*concValue[T]], r R[T], n int) *concValue[T] {
	v := acquireConcValue(pool)
	v.ref = r
	v.counts = slices.Grow(v.counts, n)
	return v
}

func (v *concValue[T]) add() *Atomic64 {
	v.refs.Acquire()
	v.counts = append(v.counts, concValueCount{})

	i := len(v.counts) - 1
	c := &v.counts[i]
	n := &c.count
	return n
}

func (v *concValue[T]) release() {
	if ok := v.refs.Release(); !ok {
		return
	}

	v.ref.Release()
	releaseConcValue(v)
}

func (v *concValue[T]) unwrap() T {
	return v.ref.Unwrap()
}

// refcount

type concValueCount struct {
	count Atomic64
	_     [256 - 8]byte
}

// pool

var concValuePools = pools.NewPools()

func acquireConcValue[T any](pool pools.Pool[*concValue[T]]) *concValue[T] {
	v, ok := pool.Get()
	if ok {
		return v
	}

	v = &concValue[T]{}
	v.pool = pool
	return v
}

func releaseConcValue[T any](v *concValue[T]) {
	pool := v.pool
	counts := slices2.Truncate(v.counts)

	*v = concValue[T]{}
	v.pool = pool
	v.counts = counts

	pool.Put(v)
}

func acquireConcValuePool[T any]() pools.Pool[*concValue[T]] {
	return pools.GetPool[*concValue[T]](concValuePools)
}
