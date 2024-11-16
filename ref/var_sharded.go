// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"runtime"
	"sync"

	"github.com/basecomplextech/baselibrary/opt"
)

// ShardedVar is a sharded atomic non-blocking variable which holds a value reference.
//
// The sharded variable is optimized for high contention scenarios.
// Internally it uses multiple shards to reduce contention.
//
// Fastrand is used to select a shard on access. The scalability of this approach is limited,
// but it is still faster than using a single mutex/atomic.
type ShardedVar[T any] interface {
	Var[T]
}

// NewShardedVar returns a new concurrent variable.
func NewShardedVar[T any]() ShardedVar[T] {
	return newShardedVar[T]()
}

// internal

var _ ShardedVar[any] = (*shardedVar[any])(nil)

type shardedVar[T any] struct {
	shards []*varImpl[T] // no need to use cache-line padding, contention is on varRef
	wmu    sync.RWMutex
}

func newShardedVar[T any]() *shardedVar[T] {
	cpus := runtime.NumCPU()
	v := &shardedVar[T]{
		shards: make([]*varImpl[T], cpus),
	}

	for i := range v.shards {
		v.shards[i] = newVar[T]()
	}
	return v
}

// Acquire acquires, retains and returns a value reference, or false.
func (v *shardedVar[T]) Acquire() (R[T], bool) {
	i := int(fastrand()) % len(v.shards)
	return v.shards[i].Acquire()
}

// Set sets a value, releases the previous reference.
func (v *shardedVar[T]) Set(value T) {
	v.wmu.Lock()
	defer v.wmu.Unlock()

	ref := NewNoop(value)
	defer ref.Release()

	for i := range v.shards {
		v.shards[i].SetRetain(ref)
	}
}

// SetRetain sets a value reference, retains the new one and releases the old one.
func (v *shardedVar[T]) SetRetain(ref R[T]) {
	v.wmu.Lock()
	defer v.wmu.Unlock()

	for i := range v.shards {
		v.shards[i].SetRetain(ref)
	}
}

// Unset clears the value.
func (v *shardedVar[T]) Unset() {
	v.wmu.Lock()
	defer v.wmu.Unlock()

	for i := range v.shards {
		v.shards[i].Unset()
	}
}

// Internal

// Unwrap returns the current value.
// The method must be externally synchronized.
func (v *shardedVar[T]) Unwrap() opt.Opt[T] {
	return v.shards[0].Unwrap()
}

// UnwrapRef returns the current reference.
// The method must be externally synchronized.
func (v *shardedVar[T]) UnwrapRef() opt.Opt[R[T]] {
	return v.shards[0].UnwrapRef()
}
