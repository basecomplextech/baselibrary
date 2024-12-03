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

type ConcurrentVar[T any] interface {
	Var[T]
}

// NewConcurrentVar returns a new concurrent variable.
func NewConcurrentVar[T any]() ConcurrentVar[T] {
	return newConcVar[T]()
}

// internal

var _ ConcurrentVar[any] = (*concVar[any])(nil)

type concVar[T any] struct {
	wmu    sync.RWMutex
	pool   pools.Pool[*concValue[T]]
	shards []concVarShard[T]
}

func newConcVar[T any]() *concVar[T] {
	cpus := runtime.NumCPU()
	pool := acquireConcValuePool[T]()
	shards := make([]concVarShard[T], cpus)

	return &concVar[T]{
		pool:   pool,
		shards: shards,
	}
}

// Acquire acquires, retains and returns a value reference, or false.
func (v *concVar[T]) Acquire() (R[T], bool) {
	i := int(fastrand()) % len(v.shards)
	return v.shards[i].acquire()
}

// Set sets a value, releases the previous reference.
func (v *concVar[T]) Set(value T) {
	ref := NewNoop(value)
	defer ref.Release()

	v.SetRetain(ref)
}

// SetRetain sets a value reference, retains the new one and releases the old one.
func (v *concVar[T]) SetRetain(r R[T]) {
	v.wmu.Lock()
	defer v.wmu.Unlock()

	// Make new value
	val := newConcValue(v.pool, r, len(v.shards))

	// Make new refs
	refs := make([]concValueRef[T], len(v.shards))
	for i := range refs {
		count := val.add()
		refs[i].init(val, count)
	}

	// Set shard variables
	for i := range refs {
		ref := &refs[i]
		v.shards[i].set(ref)
	}

	// Retain user reference
	r.Retain()
}

// Unset clears the value.
func (v *concVar[T]) Unset() {
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
func (v *concVar[T]) Unwrap() opt.Opt[T] {
	return v.shards[0].unwrap()
}

// UnwrapRef returns the current reference.
// The method must be externally synchronized.
func (v *concVar[T]) UnwrapRef() opt.Opt[R[T]] {
	return v.shards[0].unwrapRef()
}
