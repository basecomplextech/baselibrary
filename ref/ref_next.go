// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"fmt"
	"unsafe"
)

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
		r := &refNext[T, T1]{
			obj:    obj,
			parent: parent,
		}
		r.refs.Init(1)
		return r
	}

	s := acquireRefNext[T, T1]()
	s.obj = obj
	s.parent = parent

	r := &refNextPooled[T, T1]{refNextState: s}
	r.refs.Init(1)
	return r
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

var _ R[any] = (*refNext[any, any])(nil)

type refNext[T, T1 any] struct {
	refs   Atomic64
	parent R[T1]
	obj    T
}

func (r *refNext[T, T1]) Refcount() int64 {
	return r.refs.Refcount()
}

func (r *refNext[T, T1]) Acquire() bool {
	if ok := r.refs.Acquire(); ok {
		return true
	}

	r.Release()
	return false
}

func (r *refNext[T, T1]) Retain() {
	if ok := r.refs.Acquire(); ok {
		return
	}

	r.Release()
	panic(fmt.Sprintf("retain: %T already released", r.obj))
}

func (r *refNext[T, T1]) Release() {
	released := r.refs.Release()
	if !released {
		return
	}

	var zero T
	r.parent.Release()
	r.parent = nil
	r.obj = zero
}

func (r *refNext[T, T1]) Unwrap() T {
	v := r.refs.Refcount()
	if v <= 0 {
		panic(fmt.Sprintf("unwrap: %T already released", r.obj))
	}
	return r.obj
}
