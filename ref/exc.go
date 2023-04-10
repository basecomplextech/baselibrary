package ref

// Ext is a wrapper which indicates that an object is external, not owned and must not be retained.
//
// For example, Ext[[]byte] may indicate that the byte slice is stored outside of the heap,
// or in an arena allocator, and cannot be retained.
//
// Use strings.Clone or bytes.Clone to create an owned copy of a string or byte slice if required.
type Ext[T any] struct {
	obj T
}

// WrapExt returns an unowned reference of an object.
func WrapExt[T any](obj T) Ext[T] {
	return Ext[T]{obj}
}

// Unwrap returns the object.
func (u Ext[T]) Unwrap() T {
	return u.obj
}
