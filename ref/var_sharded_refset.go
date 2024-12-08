// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"slices"

	"github.com/basecomplextech/baselibrary/collect/slices2"
	"github.com/basecomplextech/baselibrary/pools"
)

type shardedVarRefset[T any] struct {
	refs   Atomic64
	ref    R[T]
	counts []varRefCount

	pool pools.Pool[*shardedVarRefset[T]]
}

func newShardedVarRefset[T any](pool pools.Pool[*shardedVarRefset[T]], r R[T], n int) *shardedVarRefset[T] {
	s := acquireShardedVarRefset(pool)
	s.ref = r
	s.counts = slices.Grow(s.counts, n)
	s.pool = pool
	return s
}

func (s *shardedVarRefset[T]) add() *Atomic64 {
	s.refs.Acquire()
	s.counts = append(s.counts, varRefCount{})

	i := len(s.counts) - 1
	c := &s.counts[i]
	n := &c.count
	return n
}

func (s *shardedVarRefset[T]) release() {
	if ok := s.refs.Release(); !ok {
		return
	}

	s.ref.Release()
	releaseShardedVarRefset(s)
}

func (s *shardedVarRefset[T]) unwrap() T {
	return s.ref.Unwrap()
}

// refcount

type varRefCount struct {
	count Atomic64
	_     [256 - 8]byte
}

// pool

var shardedVarRefsetPools = pools.NewPools()

func acquireShardedVarRefset[T any](pool pools.Pool[*shardedVarRefset[T]]) *shardedVarRefset[T] {
	s, ok := pool.Get()
	if ok {
		return s
	}

	s = &shardedVarRefset[T]{}
	s.pool = pool
	return s
}

func releaseShardedVarRefset[T any](s *shardedVarRefset[T]) {
	pool := s.pool
	s.reset()

	pool.Put(s)
}

func acquireVarRefsetPool[T any]() pools.Pool[*shardedVarRefset[T]] {
	return pools.GetPool[*shardedVarRefset[T]](shardedVarRefsetPools)
}

func (s *shardedVarRefset[T]) reset() {
	counts := slices2.Truncate(s.counts)

	*s = shardedVarRefset[T]{}
	s.counts = counts
}
