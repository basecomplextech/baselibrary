package ref

import (
	"fmt"
	"sync/atomic"

	"github.com/basecomplextech/baselibrary/pools"
)

var (
	_ R[any] = (*refFreerPooled[any])(nil)
	_ R[any] = (*refNextPooled[any, any])(nil)
)

// freer

type refFreerPooled[T any] struct {
	refs int64
	*refFreerState[T]
}

type refFreerState[T any] struct {
	pool  pools.Pool[*refFreerState[T]]
	freer Freer
	obj   T
}

// next

type refNextPooled[T, T1 any] struct {
	refs int64
	*refNextState[T, T1]
}

type refNextState[T, T1 any] struct {
	pool   pools.Pool[*refNextState[T, T1]]
	parent R[T1]
	obj    T
}

// Refcount

func (r *refFreerPooled[T]) Refcount() int64 {
	return atomic.LoadInt64(&r.refs)
}

func (r *refNextPooled[T, T1]) Refcount() int64 {
	return atomic.LoadInt64(&r.refs)
}

// Retain

func (r *refFreerPooled[T]) Retain() {
	v := atomic.AddInt64(&r.refs, 1)
	if v <= 1 {
		panic(fmt.Sprintf("retain: %T already released", r))
	}
}

func (r *refNextPooled[T, T1]) Retain() {
	v := atomic.AddInt64(&r.refs, 1)
	if v <= 1 {
		panic(fmt.Sprintf("retain: %T already released", r))
	}
}

// Release

func (r *refFreerPooled[T]) Release() {
	v := atomic.AddInt64(&r.refs, -1)
	switch {
	case v > 0:
		return
	case v < 0:
		panic(fmt.Sprintf("release: %T already released", r))
	}

	r.freer.Free()

	s := r.refFreerState
	r.refFreerState = nil
	releaseRefFreer(s)
}

func (r *refNextPooled[T, T1]) Release() {
	v := atomic.AddInt64(&r.refs, -1)
	switch {
	case v > 0:
		return
	case v < 0:
		panic(fmt.Sprintf("release: %T already released", r))
	}

	r.parent.Release()

	s := r.refNextState
	r.refNextState = nil
	releaseRefNext(s)
}

// Unwrap

func (r *refFreerPooled[T]) Unwrap() T {
	v := atomic.LoadInt64(&r.refs)
	if v <= 0 {
		panic(fmt.Sprintf("unwrap: %T already released", r))
	}
	return r.obj
}

func (r *refNextPooled[T, T1]) Unwrap() T {
	v := atomic.LoadInt64(&r.refs)
	if v <= 0 {
		panic(fmt.Sprintf("unwrap: %T already released", r))
	}
	return r.obj
}

// pools

var (
	refFreerPools = pools.New()
	refNextPools  = pools.New()
)

func acquireRefFreer[T any]() *refFreerState[T] {
	s, ok, pool := pools.Acquire1[T, *refFreerState[T]](refFreerPools)
	if ok {
		return s
	}

	s = &refFreerState[T]{}
	s.pool = pool
	return s
}

func acquireRefNext[T, T1 any]() *refNextState[T, T1] {
	s, ok, pool := pools.Acquire1[T, *refNextState[T, T1]](refNextPools)
	if ok {
		return s
	}

	s = &refNextState[T, T1]{}
	s.pool = pool
	return s
}

func releaseRefFreer[T any](s *refFreerState[T]) {
	p := s.pool
	*s = refFreerState[T]{}
	s.pool = p

	p.Put(s)
}

func releaseRefNext[T, T1 any](s *refNextState[T, T1]) {
	p := s.pool
	*s = refNextState[T, T1]{}
	s.pool = p

	p.Put(s)
}
