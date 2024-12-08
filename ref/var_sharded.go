// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"runtime"
	"sync"

	_ "unsafe"

	"github.com/basecomplextech/baselibrary/opt"
	"github.com/basecomplextech/baselibrary/pools"
)

// ShardedVar is a sharded rwmutex-based variable which holds a value reference.
//
// The sharded variable is optimized for high contention scenarios.
// Internally it uses multiple shards to reduce contention.
//
// Fastrand is used to select a shard on access. The scalability of this approach is not linear,
// but it is still faster than using a single mutex/atomic.
//
// Benchmarks
//
//	cpu: Apple M1 Pro
//	BenchmarkShardedVar-10                        	61101717	        19.77 ns/op	        50.58 mops	       0 B/op	       0 allocs/op
//	BenchmarkShardedVar_Parallel-10               	23703859	        51.03 ns/op	        19.60 mops	       0 B/op	       0 allocs/op
//	BenchmarkShardedVar_AcquireSet_Parallel-10    	19628034	        59.33 ns/op	        16.86 mops	       0 B/op	       0 allocs/op
//	BenchmarkShardedVar_Acquire-10                	74216665	        16.21 ns/op	        61.71 mops	       0 B/op	       0 allocs/op
//	BenchmarkShardedVar_Acquire_Parallel-10       	29167982	        42.29 ns/op	        23.64 mops	       0 B/op	       0 allocs/op
//	BenchmarkShardedVar_SetRetain-10              	 2595574	       463.50 ns/op	         2.15 mops	     160 B/op	       1 allocs/op
type ShardedVar[T any] interface {
	Var[T]
}

// NewShardedVar returns a new concurrent variable.
func NewShardedVar[T any]() ShardedVar[T] {
	return newShardedVar[T]()
}

// internal

var _ Var[any] = (*shardedVar[any])(nil)

type shardedVar[T any] struct {
	wmu    sync.Mutex
	pool   pools.Pool[*shardedVarRefset[T]]
	shards []shardedVarShard[T]
}

func newShardedVar[T any]() *shardedVar[T] {
	cpus := runtime.NumCPU()
	pool := acquireVarRefsetPool[T]()
	shards := make([]shardedVarShard[T], cpus)

	return &shardedVar[T]{
		pool:   pool,
		shards: shards,
	}
}

// Acquire acquires, retains and returns a value reference, or false.
func (v *shardedVar[T]) Acquire() (R[T], bool) {
	i := int(fastrand()) % len(v.shards)
	return v.shards[i].acquire()
}

// Set sets a value, releases the previous reference.
func (v *shardedVar[T]) Set(value T) {
	ref := NewNoop(value)
	defer ref.Release()

	v.SetRetain(ref)
}

// SetRetain sets a value reference, retains the new one and releases the old one.
func (v *shardedVar[T]) SetRetain(ref R[T]) {
	v.wmu.Lock()
	defer v.wmu.Unlock()

	// Retain user reference
	ref.Retain()
	done := false
	defer func() {
		if !done {
			ref.Release()
		}
	}()

	// Make new refset
	set := newShardedVarRefset(v.pool, ref, len(v.shards))

	// Make new references
	refs := make([]shardedVarRef[T], len(v.shards))
	for i := range refs {
		count := set.add()
		refs[i].init(count, set)
	}

	// Set shard variables
	for i := range refs {
		ref := &refs[i]
		v.shards[i].set(ref)
	}

	// Done
	done = true
}

// Unset clears the value.
func (v *shardedVar[T]) Unset() {
	v.wmu.Lock()
	defer v.wmu.Unlock()

	// Clear shard variables
	for i := range v.shards {
		v.shards[i].unset()
	}
}

// Internal

// Unwrap returns the current value.
// The method must be externally synchronized.
func (v *shardedVar[T]) Unwrap() opt.Opt[T] {
	return v.shards[0].unwrap()
}

// UnwrapRef returns the current reference.
// The method must be externally synchronized.
func (v *shardedVar[T]) UnwrapRef() opt.Opt[R[T]] {
	return v.shards[0].unwrapRef()
}
