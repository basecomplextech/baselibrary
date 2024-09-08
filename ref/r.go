// Copyright 2021 Ivan Korobkov. All rights reserved.

package ref

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

// R is a generic atomic countable reference.
// It wraps an object and frees it when refcount reaches 0.
type R[T any] interface {
	// Refcount returns the number of current references.
	Refcount() int64

	// Retain increments refcount, panics when count is <= 0.
	Retain()

	// Release decrements refcount and releases the object if the count is 0.
	Release()

	// Unwrap returns the object or panics if the refcount is 0.
	Unwrap() T
}

// New returns a new reference with refcount 1.
func New[T Freer](obj T) R[T] {
	return &ref[T]{
		refs: 1,
		obj:  obj,
	}
}

// NewFree returns a new reference with a free function.
func NewFree[T any](obj T, free func()) R[T] {
	freer := freeFunc(free)
	return NewFreer(obj, freer)
}

// NewFreer returns a new reference with a custom freer.
func NewFreer[T any](obj T, freer Freer) R[T] {
	size := unsafe.Sizeof(obj)
	if size <= maxUnpooledSize {
		return &refFreer[T]{
			refs:  1,
			obj:   obj,
			freer: freer,
		}
	}

	s := acquireRefFreer[T]()
	s.obj = obj
	s.freer = freer

	return &refFreerPooled[T]{
		refs:          1,
		refFreerState: s,
	}
}

// NewNoop returns a new reference with no freer.
func NewNoop[T any](obj T) R[T] {
	return &refNoop[T]{
		refs: 1,
		obj:  obj,
	}
}

// Map returns a reference which reuses parent's reference counter.
// The method allows to cast or map an object and return another reference.
//
// Example:
//
//	func map(event ref.R[Event]) ref.R[*UsedAdded] {
//		v := event.Unwrap().(*UserAdded)
//		return ref.Map[*UserAdded, Event](v, event)
//	}
func Map[T, T1 any](obj T, parent R[T1]) R[T] {
	return &refMap[T, T1]{
		parent: parent,
		obj:    obj,
	}
}

// Next returns a child reference with a parent reference as a freer.
// The parent is not retained.
//
// Example:
//
//	func parse(buf ref.R[buffer.Buffer]) ref.R[*Event] {
//		event := parseEvent(buf.Unwrap())
//		return ref.Next(event, buf)
//	}
func Next[T, T1 any](obj T, parent R[T1]) R[T] {
	size := unsafe.Sizeof(obj)
	if size <= maxUnpooledSize {
		return &refNext[T, T1]{
			refs:   1,
			obj:    obj,
			parent: parent,
		}
	}

	s := acquireRefNext[T, T1]()
	s.obj = obj
	s.parent = parent

	return &refNextPooled[T, T1]{
		refs:         1,
		refNextState: s,
	}
}

// NextRetain returns a new reference with a parent reference as a freer, retains the parent.
//
// Example:
//
//	func parse(buf ref.R[buffer.Buffer]) ref.R[*Event] {
//		event := parseEvent(buf.Unwrap())
//		return ref.NextRetain(event, buf)
//	}
func NextRetain[T, T1 any](obj T, parent R[T1]) R[T] {
	ref := Next(obj, parent)
	parent.Retain()
	return ref
}

// internal

// maxUnpooledSize specifies the maximum object size before using a pooled reference.
const maxUnpooledSize = 24

// types

type ref[T Freer] struct {
	refs int64
	obj  T
}

type refFreer[T any] struct {
	refs  int64
	freer Freer
	obj   T
}

type refNoop[T any] struct {
	refs int64
	obj  T
}

type refNext[T, T1 any] struct {
	refs   int64
	parent R[T1]
	obj    T
}

type refMap[T, T1 any] struct {
	parent R[T1]
	obj    T
}

// Refcount

func (r *ref[T]) Refcount() int64 {
	return atomic.LoadInt64(&r.refs)
}

func (r *refFreer[T]) Refcount() int64 {
	return atomic.LoadInt64(&r.refs)
}

func (r *refNoop[T]) Refcount() int64 {
	return atomic.LoadInt64(&r.refs)
}

func (r *refNext[T, T1]) Refcount() int64 {
	return atomic.LoadInt64(&r.refs)
}

func (r *refMap[T, T1]) Refcount() int64 {
	return r.parent.Refcount()
}

// Retain

func (r *ref[T]) Retain() {
	v := atomic.AddInt64(&r.refs, 1)
	if v <= 1 {
		panic(fmt.Sprintf("retain: %T already released", r.obj))
	}
}

func (r *refFreer[T]) Retain() {
	v := atomic.AddInt64(&r.refs, 1)
	if v <= 1 {
		panic(fmt.Sprintf("retain: %T already released", r.obj))
	}
}

func (r *refNoop[T]) Retain() {
	v := atomic.AddInt64(&r.refs, 1)
	if v <= 1 {
		panic(fmt.Sprintf("retain: %T already released", r.obj))
	}
}

func (r *refNext[T, T1]) Retain() {
	v := atomic.AddInt64(&r.refs, 1)
	if v <= 1 {
		panic(fmt.Sprintf("retain: %T already released", r.obj))
	}
}

func (r *refMap[T, T1]) Retain() {
	r.parent.Retain()
}

// Release

func (r *ref[T]) Release() {
	v := atomic.AddInt64(&r.refs, -1)
	switch {
	case v > 0:
		return
	case v < 0:
		panic(fmt.Sprintf("release: %T already released", r.obj))
	}

	var zero T
	r.obj.Free()
	r.obj = zero
}

func (r *refFreer[T]) Release() {
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

func (r *refNoop[T]) Release() {
	v := atomic.AddInt64(&r.refs, -1)
	switch {
	case v > 0:
		return
	case v < 0:
		panic(fmt.Sprintf("release: %T already released", r.obj))
	}
}

func (r *refNext[T, T1]) Release() {
	v := atomic.AddInt64(&r.refs, -1)
	switch {
	case v > 0:
		return
	case v < 0:
		panic(fmt.Sprintf("release: %T already released", r.obj))
	}

	var zero T
	r.parent.Release()
	r.parent = nil
	r.obj = zero
}

func (r *refMap[T, T1]) Release() {
	r.parent.Release()
	v := r.parent.Refcount()
	if v > 0 {
		return
	}

	var zero T
	r.obj = zero
}

// Unwrap

func (r *ref[T]) Unwrap() T {
	v := atomic.LoadInt64(&r.refs)
	if v <= 0 {
		panic(fmt.Sprintf("unwrap: %T already released", r.obj))
	}
	return r.obj
}

func (r *refFreer[T]) Unwrap() T {
	v := atomic.LoadInt64(&r.refs)
	if v <= 0 {
		panic(fmt.Sprintf("unwrap: %T already released", r.obj))
	}
	return r.obj
}

func (r *refNoop[T]) Unwrap() T {
	v := atomic.LoadInt64(&r.refs)
	if v <= 0 {
		panic(fmt.Sprintf("unwrap: %T already released", r.obj))
	}
	return r.obj
}

func (r *refNext[T, T1]) Unwrap() T {
	v := atomic.LoadInt64(&r.refs)
	if v <= 0 {
		panic(fmt.Sprintf("unwrap: %T already released", r.obj))
	}
	return r.obj
}

func (r *refMap[T, T1]) Unwrap() T {
	v := r.parent.Refcount()
	if v <= 0 {
		panic(fmt.Sprintf("unwrap: %T already released", r.obj))
	}
	return r.obj
}
