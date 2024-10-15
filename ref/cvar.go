// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"fmt"
	"sync"
	"sync/atomic"
	_ "unsafe"

	"github.com/basecomplextech/baselibrary/opt"
)

// ConcurrentVar is a variable which holds a reference to a value.
//
// The variable is optimized for high contention scenarios.
// Internally it uses multiple slots and chained references to reduce contention.
type ConcurrentVar[T any] interface {
	// Acquire acquires and returns a reference, or false if empty.
	Acquire() (R[T], bool)

	// Clear clears the value, releases the previous reference.
	Clear()

	// Swas swaps the value, releases the previous reference.
	Swap(v T)

	// SwapRetain swaps the value reference, retains the new one and releases the old one.
	SwapRetain(r R[T])

	// Unwrap returns the current value.
	// The method must be externally synchronized.
	Unwrap() opt.Opt[T]

	// UnwrapRef returns the current reference.
	// The method must be externally synchronized.
	UnwrapRef() opt.Opt[R[T]]
}

// NewConcurrentVar returns an empty concurrent variable.
func NewConcurrentVar[T any]() ConcurrentVar[T] {
	return &concurrentVar[T]{}
}

// internal

const concurrentNum = 16

var _ ConcurrentVar[int] = (*concurrentVar[int])(nil)

type concurrentVar[T any] struct {
	wmu sync.Mutex

	cur   opt.Opt[R[T]]
	slots [concurrentNum]concurrentSlot[T]
}

// Acquire acquires and returns a reference, or false if empty.
func (v *concurrentVar[T]) Acquire() (R[T], bool) {
	i := fastrand() % concurrentNum
	return v.slots[i].acquire().Unwrap()
}

// Clear clears the value, releases the previous reference.
func (v *concurrentVar[T]) Clear() {
	v.wmu.Lock()
	defer v.wmu.Unlock()

	// Clear all slots
	var prevs [concurrentNum]opt.Opt[R[T]]
	for i := range v.slots {
		prev := v.slots[i].clear()
		prevs[i] = prev
	}

	// Clear current
	v.cur.Unset()

	// Release previous
	for _, prev := range prevs {
		r, ok := prev.Unwrap()
		if ok {
			r.Release()
		}
	}
}

// Swas swaps the value, releases the previous reference.
func (v *concurrentVar[T]) Swap(val T) {
	r := NewNoop(val)
	defer r.Release()

	v.SwapRetain(r)
}

// SwapRetain swaps the value reference, retains the new one and releases the old one.
func (v *concurrentVar[T]) SwapRetain(r R[T]) {
	v.wmu.Lock()
	defer v.wmu.Unlock()

	// Alloc refs as a single slice
	// Even with false sharing, it's faster than reusing the same ref.
	nexts := make([]concurrentRef[T], concurrentNum)

	// Swap slots
	var prevs [concurrentNum]opt.Opt[R[T]]
	for i := range v.slots {
		next := &nexts[i]
		next.init(r)

		prev := v.slots[i].swap(next)
		prevs[i] = prev
	}

	// Swap current
	v.cur.Set(r)

	// Retain next
	for i := 0; i < concurrentNum; i++ {
		r.Retain()
	}

	// Release previous
	for _, prev := range prevs {
		r, ok := prev.Unwrap()
		if ok {
			r.Release()
		}
	}
}

// Unwrap returns the current value.
// The method must be externally synchronized.
func (v *concurrentVar[T]) Unwrap() opt.Opt[T] {
	v.wmu.Lock()
	defer v.wmu.Unlock()

	r, ok := v.cur.Unwrap()
	if !ok {
		return opt.Opt[T]{}
	}

	val := r.Unwrap()
	return opt.New(val)
}

// UnwrapRef returns the current reference.
// The method must be externally synchronized.
func (v *concurrentVar[T]) UnwrapRef() opt.Opt[R[T]] {
	v.wmu.Lock()
	defer v.wmu.Unlock()

	return v.cur
}

// slot

type concurrentSlot[T any] struct {
	mu   sync.Mutex
	ref  opt.Opt[R[T]]
	_pad [224]byte
}

func (s *concurrentSlot[T]) acquire() opt.Opt[R[T]] {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.ref.Unwrap()
	if !ok {
		return opt.Opt[R[T]]{}
	}

	r.Retain()
	return opt.New(r)
}

func (s *concurrentSlot[T]) clear() opt.Opt[R[T]] {
	s.mu.Lock()
	defer s.mu.Unlock()

	prev := s.ref
	s.ref = opt.Opt[R[T]]{}
	return prev
}

func (s *concurrentSlot[T]) swap(r R[T]) opt.Opt[R[T]] {
	s.mu.Lock()
	defer s.mu.Unlock()

	prev := s.ref
	s.ref = opt.New(r)
	return prev
}

// ref

var _ R[int] = (*concurrentRef[int])(nil)

type concurrentRef[T any] struct {
	refs int64
	ref  R[T]
}

func (r *concurrentRef[T]) init(ref R[T]) {
	r.refs = 1
	r.ref = ref
}

// Refcount returns the number of current references.
func (r *concurrentRef[T]) Refcount() int64 {
	return atomic.LoadInt64(&r.refs)
}

// Retain increments refcount, panics when count is <= 0.
func (r *concurrentRef[T]) Retain() {
	v := atomic.AddInt64(&r.refs, 1)
	if v <= 1 {
		var zero T
		panic(fmt.Sprintf("retain: %T already released", zero))
	}
}

// Release decrements refcount and releases the object if the count is 0.
func (r *concurrentRef[T]) Release() {
	v := atomic.AddInt64(&r.refs, -1)
	switch {
	case v > 0:
		return
	case v < 0:
		var zero T
		panic(fmt.Sprintf("release: %T already released", zero))
	}

	r.ref.Release()
	r.ref = nil
}

// Unwrap returns the object or panics if the refcount is 0.
func (r *concurrentRef[T]) Unwrap() T {
	return r.ref.Unwrap()
}

// util

//go:linkname fastrand runtime.fastrand
func fastrand() uint32
