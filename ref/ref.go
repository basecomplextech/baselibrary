package ref

// Ref is a countable reference interface without generics, i.e. R[?].
type Ref interface {
	// Retain increments refcount, panics when count is <= 0.
	Retain()

	// Release decrements refcount and releases the object if the count is 0.
	Release()
}

// Retain retains and returns a reference.
//
// Usage:
//
//	tree.table = Retain(table)
func Retain[R Ref](r R) R {
	r.Retain()
	return r
}

// RetainAll retains all references.
func RetainAll[R Ref](refs ...R) []R {
	for _, r := range refs {
		r.Retain()
	}
	return refs
}

// ReleaseAll releases all references.
func ReleaseAll[R Ref](refs ...R) {
	for _, r := range refs {
		r.Release()
	}
}

// SwapRetain retains a new reference, releases an old one and returns the new.
//
// Usage:
//
//	tbl := table.Clone()
//	defer tbl.Release()
//	...
//	s.table = SwapRetain(s.table, tbl)
func SwapRetain[R Ref](old R, new R) R {
	new.Retain()
	old.Release()
	return new
}

// SwapNoRetain releases an old reference, and returns a new one.
//
// Usage:
//
//	tbl := newTable()
//	...
//	s.table = SwapNoRetain(s.table, tbl)
func SwapNoRetain[R Ref](old R, new R) R {
	old.Release()
	return new
}
