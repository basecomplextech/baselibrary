// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"fmt"
	"unsafe"

	"github.com/basecomplextech/baselibrary/pools"
)

// NewFree returns a new reference with a free function.
func NewFree[T any](obj T, free func()) R[T] {
	freer := freeFunc(free)
	return NewFreer(obj, freer)
}

// NewFreer returns a new reference with a custom freer.
func NewFreer[T any](obj T, freer Freer) R[T] {
	size := unsafe.Sizeof(obj)
	if size <= maxUnpooledSize {
		r := &refFreer[T]{
			obj:   obj,
			freer: freer,
		}
		r.refs.Init(1)
		return r
	}

	s := acquireRefFreer[T]()
	s.obj = obj
	s.freer = freer

	r := &refFreerPooled[T]{refFreerState: s}
	r.refs.Init(1)
	return r
}

// internal

var _ R[any] = (*refFreer[any])(nil)

type refFreer[T any] struct {
	refs  Atomic64
	freer Freer
	obj   T
}

func (r *refFreer[T]) Refcount() int64 {
	return r.refs.Refcount()
}

// Acquire tries to increment refcount and returns true, or false if already released.
func (r *refFreer[T]) Acquire() bool {
	if ok := r.refs.Acquire(); ok {
		return true
	}

	r.Release()
	return false
}

func (r *refFreer[T]) Retain() {
	if ok := r.refs.Acquire(); ok {
		return
	}

	r.Release()
	panic(fmt.Sprintf("retain: %T already released", r.obj))
}

func (r *refFreer[T]) Release() {
	released := r.refs.Release()
	if !released {
		return
	}

	var zero T
	r.freer.Free()
	r.freer = nil
	r.obj = zero
}

func (r *refFreer[T]) Unwrap() T {
	v := r.refs.Refcount()
	if v <= 0 {
		panic(fmt.Sprintf("unwrap: %T already released", r.obj))
	}
	return r.obj
}

// pooled

var _ R[any] = (*refFreerPooled[any])(nil)

type refFreerPooled[T any] struct {
	refs Atomic64
	*refFreerState[T]
}

type refFreerState[T any] struct {
	pool  pools.Pool[*refFreerState[T]]
	freer Freer
	obj   T
}

func (r *refFreerPooled[T]) Refcount() int64 {
	return r.refs.Refcount()
}

func (r *refFreerPooled[T]) Acquire() bool {
	if ok := r.refs.Acquire(); ok {
		return true
	}

	r.Release()
	return false
}

func (r *refFreerPooled[T]) Retain() {
	if ok := r.refs.Acquire(); ok {
		return
	}

	r.Release()
	panic(fmt.Sprintf("retain: %T already released", r))
}

func (r *refFreerPooled[T]) Release() {
	released := r.refs.Release()
	if !released {
		return
	}

	r.freer.Free()

	s := r.refFreerState
	r.refFreerState = nil
	releaseRefFreer(s)
}

func (r *refFreerPooled[T]) Unwrap() T {
	v := r.refs.Refcount()
	if v <= 0 {
		panic(fmt.Sprintf("unwrap: %T already released", r))
	}
	return r.obj
}

// pool

var refFreerPools = pools.NewPools()

func acquireRefFreer[T any]() *refFreerState[T] {
	s, ok, pool := pools.Acquire1[*refFreerState[T]](refFreerPools)
	if ok {
		return s
	}

	s = &refFreerState[T]{}
	s.pool = pool
	return s
}

func releaseRefFreer[T any](s *refFreerState[T]) {
	p := s.pool
	*s = refFreerState[T]{}
	s.pool = p

	p.Put(s)
}
