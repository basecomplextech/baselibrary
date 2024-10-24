// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"sync"
	"sync/atomic"

	"github.com/basecomplextech/baselibrary/opt"
)

// Var is an atomic non-blocking variable which holds a value reference.
type Var[T any] interface {
	// Acquire acquires, retains and returns a value reference, or false.
	Acquire() (R[T], bool)

	// Set sets a value, releases the previous reference.
	Set(T)

	// SetRetain sets a value reference, retains the new one and releases the old one.
	SetRetain(R[T])

	// Unset clears the value.
	Unset()

	// Internal

	// Unwrap returns the current value.
	// The method must be externally synchronized.
	Unwrap() opt.Opt[T]

	// UnwrapRef returns the current reference.
	// The method must be externally synchronized.
	UnwrapRef() opt.Opt[R[T]]
}

// NewVar returns a new empty atomic variable.
func NewVar[T any]() Var[T] {
	return newVar[T]()
}

// internal

var _ Var[any] = (*varImpl[any])(nil)

type varImpl[T any] struct {
	cur atomic.Pointer[varRef[T]]
	wmu sync.Mutex // write mutex
}

func newVar[T any]() *varImpl[T] {
	return &varImpl[T]{}
}

// Acquire acquires, retains and returns a value reference, or false.
func (v *varImpl[T]) Acquire() (R[T], bool) {
	for {
		ref := v.cur.Load()
		if ref == nil {
			return nil, false
		}

		if ref.acquire() {
			return ref, true
		}
	}
}

// Set sets a value, releases the previous reference.
func (v *varImpl[T]) Set(value T) {
	ref := NewNoop(value)
	defer ref.Release()

	v.SetRetain(ref)
}

// SetRetain sets a value reference, retains the new one and releases the old one.
func (v *varImpl[T]) SetRetain(ref R[T]) {
	v.wmu.Lock()
	defer v.wmu.Unlock()

	// Swap the current reference
	next := newVarRef(ref)
	prev := v.cur.Swap(next)
	ref.Retain()

	// Detach and release the previous one
	if prev != nil {
		prev.detach()
	}
}

// Unset clears the value.
func (v *varImpl[T]) Unset() {
	v.wmu.Lock()
	defer v.wmu.Unlock()

	// Clear the current reference
	prev := v.cur.Swap(nil)

	// Detach and release the previous one
	if prev != nil {
		prev.detach()
	}
}

// Internal

// Unwrap returns the current value.
// The method must be externally synchronized.
func (v *varImpl[T]) Unwrap() opt.Opt[T] {
	ref := v.cur.Load()
	if ref == nil {
		return opt.Opt[T]{}
	}

	val := ref.Unwrap()
	return opt.New(val)
}

// UnwrapRef returns the current reference.
// The method must be externally synchronized.
func (v *varImpl[T]) UnwrapRef() opt.Opt[R[T]] {
	ref := v.cur.Load()
	if ref == nil {
		return opt.Opt[R[T]]{}
	}

	return opt.New[R[T]](ref)
}
