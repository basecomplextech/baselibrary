// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"slices"

	"github.com/basecomplextech/baselibrary/collect/slices2"
	"github.com/basecomplextech/baselibrary/pools"
)

type shardedValue[T any] struct {
	refs   Atomic64
	ref    R[T]
	counts []shardedValueCount

	pool pools.Pool[*shardedValue[T]]
}

func newShardedValue[T any](pool pools.Pool[*shardedValue[T]], r R[T], n int) *shardedValue[T] {
	v := acquireShardedValue(pool)
	v.ref = r
	v.counts = slices.Grow(v.counts, n)
	return v
}

func (v *shardedValue[T]) add() *Atomic64 {
	v.refs.Acquire()
	v.counts = append(v.counts, shardedValueCount{})

	i := len(v.counts) - 1
	c := &v.counts[i]
	n := &c.count
	return n
}

func (v *shardedValue[T]) release() {
	if ok := v.refs.Release(); !ok {
		return
	}

	v.ref.Release()
	releaseShardedValue(v)
}

func (v *shardedValue[T]) unwrap() T {
	return v.ref.Unwrap()
}

// refcount

type shardedValueCount struct {
	count Atomic64
	_     [256 - 8]byte
}

// pool

var shardedValuePools = pools.NewPools()

func acquireShardedValue[T any](pool pools.Pool[*shardedValue[T]]) *shardedValue[T] {
	v, ok := pool.Get()
	if ok {
		return v
	}

	v = &shardedValue[T]{}
	v.pool = pool
	return v
}

func releaseShardedValue[T any](v *shardedValue[T]) {
	pool := v.pool
	counts := slices2.Truncate(v.counts)

	*v = shardedValue[T]{}
	v.pool = pool
	v.counts = counts

	pool.Put(v)
}

func acquireShardedValuePool[T any]() pools.Pool[*shardedValue[T]] {
	return pools.GetPool[*shardedValue[T]](shardedValuePools)
}
