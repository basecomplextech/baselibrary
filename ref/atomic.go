package ref

import (
	"fmt"
	"sync/atomic"
)

// R is a generic atomic countable reference.
// It wraps an object and frees it when refcount reaches 0.
type R[T Freer] struct {
	obj   T
	refs  *int64 // refcount pointer, can be shared with another R[T], see Cast
	_refs int64  // actual refcount value
}

// New returns a new reference with refcount 1.
func New[T Freer](obj T) *R[T] {
	r := &R[T]{
		obj:   obj,
		_refs: 1,
	}
	r.refs = &r._refs
	return r
}

// Cast casts one interface to another and returns another reference which shares the same refcount.
func Cast[T, T1 Freer](ref *R[T]) *R[T1] {
	obj0 := (any)(ref.obj)
	obj1 := obj0.(T1)

	r1 := &R[T1]{
		obj:  obj1,
		refs: ref.refs,
	}
	return r1
}

// Refcount returns the number of current references.
func (r *R[T]) Refcount() int64 {
	return atomic.LoadInt64(r.refs)
}

// Retain increments refcount, panics when count is 0.
func (r *R[T]) Retain() {
	v := atomic.AddInt64(r.refs, 1)
	if v <= 1 {
		panic(fmt.Sprintf("retain: %T already released", r.obj))
	}
	return
}

// Release decrements refcount and releases the object if the count is 0.
func (r *R[T]) Release() {
	v := atomic.AddInt64(r.refs, -1)
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
	refs := atomic.LoadInt64(r.refs)
	if refs <= 0 {
		panic(fmt.Sprintf("unwrap: %T already released", r.obj))
	}

	return r.obj
}
