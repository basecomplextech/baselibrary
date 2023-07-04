package ref

// Box wraps any value and its freer.
type Box[T any] struct {
	obj   T
	freer Freer
}

// NewBox returns a new boxed value with a custom freer.
func NewBox[T any](obj T, freer Freer) *Box[T] {
	return &Box[T]{
		obj:   obj,
		freer: freer,
	}
}

// NewBoxFunc returns a new boxed value with a free function.
func NewBoxFunc[T any](obj T, free func()) *Box[T] {
	return &Box[T]{
		obj:   obj,
		freer: freeFunc(free),
	}
}

// NewBoxed returns a new reference with a boxed value with a custom freer.
func NewBoxed[T any](obj T, freer Freer) *R[*Box[T]] {
	return New[*Box[T]](NewBox[T](obj, freer))
}

// NewBoxedFunc returns a new reference with a boxed value and a free function.
func NewBoxedFunc[T any](obj T, free func()) *R[*Box[T]] {
	return New[*Box[T]](NewBoxFunc[T](obj, free))
}

// Free frees the wrapped object.
func (b *Box[T]) Free() {
	var zero T
	b.obj = zero
	b.freer.Free()
}

// Unwrap returns the wrapped object.
func (b *Box[T]) Unwrap() T {
	return b.obj
}

// private

// freeFunc is an adapter which allows to use a function as a Freer.
type freeFunc func()

// Free frees the object.
func (f freeFunc) Free() {
	f()
}
