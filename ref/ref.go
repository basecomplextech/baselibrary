package ref

import (
	"fmt"
	"sync/atomic"
)

// R is a generic atomically countable reference which wraps an object.
// The object is automatically freed when refcount reaches 0.
type R[T Freer] struct {
	obj  T
	refs int32
}

// Freer frees the object.
type Freer interface {
	Free()
}

// Wrap wraps an object into a reference.
func Wrap[T Freer](obj T) *R[T] {
	return &R[T]{
		obj:  obj,
		refs: 1,
	}
}

// Refcount returns the number of current references.
func (r *R[T]) Refcount() int32 {
	return r.refs
}

// Retain increments refcount, panics when count is 0.
func (r *R[T]) Retain() {
	v := atomic.AddInt32(&r.refs, 1)
	if v <= 1 {
		panic(fmt.Sprintf("retain: %T already released", r.obj))
	}
	return
}

// Release decrements refcount and releases the object if the count is 0.
func (r *R[T]) Release() {
	v := atomic.AddInt32(&r.refs, -1)
	switch {
	case v > 0:
		return
	case v < 0:
		panic(fmt.Sprintf("release: %T already released", r.obj))
	}

	r.obj.Free()
	return
}

// Unwrap returns the object or panics if the refcount is 0.
func (r *R[T]) Unwrap() T {
	if r.refs <= 0 {
		panic(fmt.Sprintf("unwrap: %T already released", r.obj))
	}

	return r.obj
}
