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

// FreeFunc returns an adapter which allows to use a function as a Freer.
func FreeFunc(f func()) Freer {
	return freeFunc(f)
}

// private

// freeFunc is an adapter which allows to use a function as a Freer.
type freeFunc func()

// Free frees the object.
func (f freeFunc) Free() {
	f()
}
