// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

var _ R[any] = (*dummyRef[any])(nil)

type dummyRef[T any] struct{}

func newDummyRef[T any]() R[T] {
	return dummyRef[T]{}
}

// Refcount returns the number of current references.
func (dummyRef[T]) Refcount() int64 { return 0 }

// Acquire tries to increment refcount and returns true, or false if already released.
func (dummyRef[T]) Acquire() bool { return false }

// Retain increments refcount, panics when count is <= 0.
func (dummyRef[T]) Retain() {}

// Release decrements refcount and releases the object if the count is 0.
func (dummyRef[T]) Release() {}

// Unwrap returns the object or panics if the refcount is 0.
func (dummyRef[T]) Unwrap() (value T) {
	return
}
