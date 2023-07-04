package ref

import (
	"fmt"
	"sync/atomic"
)

// R is a generic atomic countable reference.
// It wraps an object and frees it when refcount reaches 0.
type R[T Freer] struct {
	// only one is set, src is used only in Map[S, D]
	obj T
	src refMap[T]

	refs int64
}

// New returns a new reference with refcount 1.
func New[T Freer](obj T) *R[T] {
	return &R[T]{
		obj:  obj,
		refs: 1,
	}
}

// Map returns a new reference which retains another reference and maps its object.
// The method exists mostly to support casting in references.
//
// Example:
//
//	// *ref.R[*revision]
//	src := machine.Head()
//	defer src.Release()
//
//	// *ref.R[blockchain.Revision]
//	dst := ref.Map(src, func(r Revision) blockchain.Revision {
//		return r.(blockchain.Revision)
//	})
func Map[S Freer, D Freer](src *R[S], cast func(S) D) *R[D] {
	m := newRefMap[S, D](src, cast)

	return &R[D]{
		src:  m,
		refs: 1,
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

	if r.src != nil {
		r.src.free()
		return
	}

	r.obj.Free()
}

// Unwrap returns the object or panics if the refcount is 0.
func (r *R[T]) Unwrap() T {
	refs := atomic.LoadInt64(&r.refs)
	if refs <= 0 {
		panic(fmt.Sprintf("unwrap: %T already released", r.obj))
	}

	if r.src != nil {
		return r.src.unwrap()
	}
	return r.obj
}

// private

type refMap[T Freer] interface {
	free()
	unwrap() T
}

type refMapImpl[S Freer, D Freer] struct {
	src  *R[S]
	cast func(S) D
}

func newRefMap[S Freer, D Freer](src *R[S], cast func(S) D) refMap[D] {
	return &refMapImpl[S, D]{
		src:  Retain(src),
		cast: cast,
	}
}

func (m *refMapImpl[S, D]) free() {
	m.src.Release()
	m.src = nil
}

func (m *refMapImpl[S, D]) unwrap() D {
	return m.cast(m.src.Unwrap())
}
