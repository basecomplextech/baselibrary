// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

// Owner is an object wrapper which indicates an optional ownership of the object.
type Owner[T Freer] struct {
	Value T
	Valid bool
	Owned bool
}

// Owned returns an owner with the given value and ownership set to true.
func Owned[T Freer](value T) Owner[T] {
	return Owner[T]{
		Value: value,
		Valid: true,
		Owned: true,
	}
}

// Unowned returns an owner with the given value and ownership set to false.
func Unowned[T Freer](value T) Owner[T] {
	return Owner[T]{
		Value: value,
		Valid: true,
		Owned: false,
	}
}

// Clear clears the owner.
func (o *Owner[T]) Clear() {
	*o = Owner[T]{}
}

// SetOwned sets the ownership to true.
func (o *Owner[T]) SetOwned() {
	o.Owned = true
}

// Unowned returns a new owner with the same value and ownership set to false.
func (o Owner[T]) Unowned() Owner[T] {
	return Owner[T]{
		Value: o.Value,
		Valid: o.Valid,
		Owned: false,
	}
}

// Unwrap

// Unwrap returns the underlying value and a boolean indicating whether the value is valid.
func (o Owner[T]) Unwrap() (T, bool) {
	return o.Value, o.Valid
}

// MustUnwrap returns the underlying value if it is valid, otherwise it panics.
func (o Owner[T]) MustUnwrap() T {
	if !o.Valid {
		panic("value is not set")
	}
	return o.Value
}

// Free

// Free releases the owned object if it is owned.
func (o *Owner[T]) Free() {
	if o.Valid && o.Owned {
		o.Value.Free()
		*o = Owner[T]{}
	}
}
