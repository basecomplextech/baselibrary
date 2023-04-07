package ref

// Unowned is a wrapper which indicates that the object is not owned and must not be retained.
//
// For example, Unowned[[]byte] may indicate that the byte slice is stored outside the heap,
// or in an arena allocator, and cannot be retained.
//
// Use strings.Clone or bytes.Clone to create an owned copy of a string or byte slice if required.
type Unowned[T any] struct {
	obj T
}

// Unown returns an unowned reference of an object.
func Unown[T any](obj T) Unowned[T] {
	return Unowned[T]{obj}
}

// Unwrap returns the object.
func (u Unowned[T]) Unwrap() T {
	return u.obj
}
