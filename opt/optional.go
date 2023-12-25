// Package opt provides a generic optional value.
package opt

// Value is an optional value, an empty value is undefined.
type Value[T any] struct {
	Valid bool
	Value T
}

// New returns a new valid value.
func New[T any](value T) Value[T] {
	return Value[T]{
		Valid: true,
		Value: value,
	}
}

// None returns a new undefined value.
func None[T any]() Value[T] {
	return Value[T]{}
}

// Unwrap returns the value and true if valid.
func (v Value[T]) Unwrap() (T, bool) {
	return v.Value, v.Valid
}

// MustUnwrap returns the value or panics if not valid.
func (v Value[T]) MustUnwrap() T {
	if !v.Valid {
		panic("undefined value")
	}
	return v.Value
}
