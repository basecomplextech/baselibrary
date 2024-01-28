// Package opt provides a generic optional value.
package opt

// Opt is an optional value, an empty value is unset.
type Opt[T any] struct {
	Set   bool
	Value T
}

// New returns a new set value.
func New[T any](value T) Opt[T] {
	return Opt[T]{
		Set:   true,
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
		Set:   true,
		Value: value,
	}
}

// None returns a new unset value.
func None[T any]() Opt[T] {
	return Opt[T]{}
}

// Unwrap returns the value and true if set.
func (o Opt[T]) Unwrap() (T, bool) {
	return o.Value, o.Set
}

// MustUnwrap returns the value or panics if not set.
func (o Opt[T]) MustUnwrap() T {
	if !o.Set {
		panic("unset value")
	}
	return o.Value
}
