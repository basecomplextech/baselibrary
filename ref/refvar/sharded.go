// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package refvar

import (
	"runtime"
	"sync"

	_ "unsafe"

	"github.com/basecomplextech/baselibrary/opt"
	"github.com/basecomplextech/baselibrary/pools"
	"github.com/basecomplextech/baselibrary/ref"
)

// ShardedVar is a sharded rwmutex-based variable which holds a value reference.
//
// The sharded variable is optimized for high contention scenarios.
// Internally it uses multiple shards to reduce contention.
//
// Unsafe procPin/procUnpin functions are used to select a shard to access.
// This scales well with the number of CPUs.
//
// Benchmarks
//
//	cpu: Apple M1 Pro
//	BenchmarkShardedVar-10    						71493033	        16.61 ns/op	        60.21 mops	       0 B/op	       0 allocs/op
//	BenchmarkShardedVar_Parallel-10    				587905798	         2.06 ns/op	       485.20 mops	       0 B/op	       0 allocs/op
//	BenchmarkShardedVar_AcquireSet_Parallel-10    	100000000	        13.73 ns/op	        72.83 mops	     12787 sets	       0 B/op	       0 allocs/op
//	BenchmarkShardedVar_Acquire-10    				82076718	        14.38 ns/op	        69.56 mops	       0 B/op	       0 allocs/op
//	BenchmarkShardedVar_Acquire_Parallel-10    		541823913	         2.31 ns/op	       431.70 mops	       0 B/op	       0 allocs/op
//	BenchmarkShardedVar_SetRetain-10    	 		2402602	           481.70 ns/op	         2.07 mops	     160 B/op	       1 allocs/op
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
func (v *shardedVar[T]) Acquire() (ref.R[T], bool) {
	// Fastrand is left in case procPin/procUnpin functions are not available later.
	// i := int(fastrand()) % len(v.shards)

	i := procPin() % len(v.shards)
	shard := &v.shards[i]
	procUnpin()
	return shard.acquire()
}

// Set sets a value, releases the previous reference.
func (v *shardedVar[T]) Set(value T) {
	ref := ref.NewNoop(value)
	defer ref.Release()

	v.SetRetain(ref)
}

// SetRetain sets a value reference, retains the new one and releases the old one.
func (v *shardedVar[T]) SetRetain(ref ref.R[T]) {
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
func (v *shardedVar[T]) UnwrapRef() opt.Opt[ref.R[T]] {
	return v.shards[0].unwrapRef()
}

// private

//go:linkname procPin runtime.procPin
func procPin() int

//go:linkname procUnpin runtime.procUnpin
func procUnpin()
