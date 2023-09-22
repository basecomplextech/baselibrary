package ref

import (
	"fmt"
	"sync/atomic"
)

// R is a generic atomic countable reference.
// It wraps an object and frees it when refcount reaches 0.
type R[T any] struct {
	obj   T
	freer Freer

	refs int64
}

// New returns a new reference with refcount 1.
func New[T Freer](obj T) *R[T] {
	return &R[T]{
		obj:   obj,
		freer: obj,
		refs:  1,
	}
}

// NewFree returns a new reference with a free function.
func NewFree[T any](obj T, free func()) *R[T] {
	return &R[T]{
		obj:   obj,
		freer: freeFunc(free),
		refs:  1,
	}
}

// NewFreer returns a new reference with a custom freer.
func NewFreer[T any](obj T, freer Freer) *R[T] {
	return &R[T]{
		obj:   obj,
		freer: freer,
		refs:  1,
	}
}

// NewNoFreer returns a new reference with no freer.
func NewNoFreer[T any](obj T) *R[T] {
	return &R[T]{
		obj:   obj,
		freer: NoopFreer,
		refs:  1,
	}
}

// NewParent returns a new reference with a parent reference as a freer.
// The parent is not retained.
func NewParent[T any, T1 any](obj T, parent *R[T1]) *R[T] {
	r := (*refFreer[T1])(parent)

	return &R[T]{
		obj:   obj,
		freer: r,
		refs:  1,
	}
}

// NewParentRetain returns a new reference with a parent reference as a freer, retains the parent.
func NewParentRetain[T any, T1 any](obj T, parent *R[T1]) *R[T] {
	parent.Retain()
	r := (*refFreer[T1])(parent)

	return &R[T]{
		obj:   obj,
		freer: r,
		refs:  1,
	}
}

// Refcount returns the number of current references.
func (r *R[T]) Refcount() int64 {
	return atomic.LoadInt64(&r.refs)
}

// Retain increments refcount, panics when count is 0.
func (r *R[T]) Retain() {
	v := atomic.AddInt64(&r.refs, 1)
	if v <= 1 {
		panic(fmt.Sprintf("retain: %T already released", r.obj))
	}
	return
}

// Release decrements refcount and releases the object if the count is 0.
func (r *R[T]) Release() {
	v := atomic.AddInt64(&r.refs, -1)
	switch {
	case v > 0:
		return
	case v < 0:
		panic(fmt.Sprintf("release: %T already released", r.obj))
	}

	var zero T
	r.freer.Free()
	r.freer = nil
	r.obj = zero
}

// Unwrap returns the object or panics if the refcount is 0.
func (r *R[T]) Unwrap() T {
	refs := atomic.LoadInt64(&r.refs)
	if refs <= 0 {
		panic(fmt.Sprintf("unwrap: %T already released", r.obj))
	}

	return r.obj
}
