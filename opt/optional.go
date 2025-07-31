// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

// Package opt provides a generic optional value.
package opt

// Opt is an optional value, an empty value is unset.
type Opt[T any] struct {
	Valid bool
	Value T
}

// New returns a new set value.
func New[T any](value T) Opt[T] {
	return Opt[T]{
		Valid: true,
		Value: value,
	}
}

// Maybe returns a new set value if not zero, otherwise an unset value.
func Maybe[T comparable](value T) Opt[T] {
	var zero T
	if value == zero {
		return Opt[T]{}
	}

	return Opt[T]{
		Valid: true,
		Value: value,
	}
}

// None returns a new unset value.
func None[T any]() Opt[T] {
	return Opt[T]{}
}

// Mutate

// Clear clears the value and returns the previous one if any.
func (o *Opt[T]) Clear() (v T, ok bool) {
	v, ok = o.Value, o.Valid
	*o = Opt[T]{}
	return
}

// Set sets the value.
func (o *Opt[T]) Set(value T) {
	*o = Opt[T]{
		Valid: true,
		Value: value,
	}
}

// Unwrap

// Or returns the value if set, otherwise returns the default value.
func (o Opt[T]) Or(default_ T) T {
	if o.Valid {
		return o.Value
	}
	return default_
}

// Unwrap returns the value and true if set.
func (o Opt[T]) Unwrap() (T, bool) {
	return o.Value, o.Valid
}

// MustUnwrap returns the value or panics if not set.
func (o Opt[T]) MustUnwrap() T {
	if !o.Valid {
		panic("unset value")
	}
	return o.Value
}
