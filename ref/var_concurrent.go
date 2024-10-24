// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"runtime"
	"sync"

	"github.com/basecomplextech/baselibrary/opt"
)

// ConcurrentVar is a sharded atomic non-blocking variable which holds a value reference.
//
// The concurrent variable is optimized for high contention scenarios.
// Internally it uses multiple shards to reduce contention.
//
// Fastrand is used to select a shard on access. The scalability of this approach is limited,
// but it is still faster than using a single mutex/atomic.
type ConcurrentVar[T any] interface {
	Var[T]
}

// NewConcurrentVar returns a new concurrent variable.
func NewConcurrentVar[T any]() ConcurrentVar[T] {
	return newCVar[T]()
}

// internal

var _ ConcurrentVar[any] = (*cvar[any])(nil)

type cvar[T any] struct {
	shards []cvarShard[T]
	wmu    sync.RWMutex
}

type cvarShard[T any] struct {
	*varImpl[T]
	_ [256 - 8]byte
}

func newCVar[T any]() *cvar[T] {
	cpus := runtime.NumCPU()
	v := &cvar[T]{
		shards: make([]cvarShard[T], cpus),
	}

	for i := range v.shards {
		v.shards[i].varImpl = newVar[T]()
	}
	return v
}

// Acquire acquires, retains and returns a value reference, or false.
func (v *cvar[T]) Acquire() (R[T], bool) {
	i := int(fastrand()) % len(v.shards)
	return v.shards[i].Acquire()
}

// Set sets a value, releases the previous reference.
func (v *cvar[T]) Set(value T) {
	v.wmu.Lock()
	defer v.wmu.Unlock()

	ref := NewNoop(value)
	defer ref.Release()

	for i := range v.shards {
		v.shards[i].SetRetain(ref)
	}
}

// SetRetain sets a value reference, retains the new one and releases the old one.
func (v *cvar[T]) SetRetain(ref R[T]) {
	v.wmu.Lock()
	defer v.wmu.Unlock()

	for i := range v.shards {
		v.shards[i].SetRetain(ref)
	}
}

// Unset clears the value.
func (v *cvar[T]) Unset() {
	v.wmu.Lock()
	defer v.wmu.Unlock()

	for i := range v.shards {
		v.shards[i].Unset()
	}
}

// Internal

// Unwrap returns the current value.
// The method must be externally synchronized.
func (v *cvar[T]) Unwrap() opt.Opt[T] {
	return v.shards[0].Unwrap()
}

// UnwrapRef returns the current reference.
// The method must be externally synchronized.
func (v *cvar[T]) UnwrapRef() opt.Opt[R[T]] {
	return v.shards[0].UnwrapRef()
}
