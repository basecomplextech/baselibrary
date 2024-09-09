// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package pools

import "sync"

// Pool is a generic pool interface, which wraps sync.Pool.
type Pool[T any] interface {
	// Get returns a value from the pool, or false.
	Get() (T, bool)

	// New returns a value from the pool, or creates a new one, panics if no new function.
	New() T

	// Put returns a value to the pool.
	Put(T)
}

// NewPool returns a new pool without a new function.
func NewPool[T any]() Pool[T] {
	return newPool[T](nil)
}

// NewPoolFunc returns a new pool with a new function.
func NewPoolFunc[T any](new func() T) Pool[T] {
	return newPool(new)
}

// internal

var _ Pool[any] = &pool[any]{}

type pool[T any] struct {
	sync.Pool
}

func newPool[T any](new func() T) Pool[T] {
	var new1 func() any
	if new != nil {
		new1 = func() any {
			return new()
		}
	}

	return &pool[T]{sync.Pool{New: new1}}
}

// Get returns a value from the pool, or false.
func (p *pool[T]) Get() (value T, ok bool) {
	v := p.Pool.Get()
	if v == nil {
		return
	}
	return v.(T), true
}

// New returns a value from the pool, or creates a new one, panics if no new function.
func (p *pool[T]) New() T {
	v := p.Pool.Get()
	if v == nil {
		panic("no pool new function")
	}
	return v.(T)
}

// Put returns a value to the pool.
func (p *pool[T]) Put(v T) {
	p.Pool.Put(v)
}
