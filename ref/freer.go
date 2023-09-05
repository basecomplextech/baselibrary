package ref

// Freer specifies an object which can be freed.
type Freer interface {
	Free()
}

// FreeAll frees all objects.
func FreeAll[T Freer](objs ...T) {
	for _, obj := range objs {
		obj.Free()
	}
}

// private

// freeFunc is an adapter which allows to use a function as a Freer.
type freeFunc func()

// Free frees the object.
func (f freeFunc) Free() {
	f()
}

// freeRef is an adapter which allows to use a reference as a Freer.
type freeRef[T any] R[T]

func (f *freeRef[T]) Free() {
	r := (*R[T])(f)
	r.Release()
}
