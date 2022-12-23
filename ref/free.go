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
