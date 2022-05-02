package ref

// Count specifies an interface which can be implemented by any countable reference.
type Count interface {
	// Retain increments refcount, panics when count is 0.
	Retain()

	// Release decrements refcount and releases the object if the count is 0.
	Release()
}

// Retain retains and returns a reference.
//
// Usage:
//	tree.table = Retain(table)
//
func Retain[C Count](count C) C {
	count.Retain()
	return count
}

// RetainAll retains all references.
func RetainAll[C Count](counts ...C) {
	for _, count := range counts {
		count.Retain()
	}
}

// ReleaseAll releases all references.
func ReleaseAll[C Count](counts ...C) {
	for _, count := range counts {
		count.Release()
	}
}

// Swap retains a new reference, releases an old one and returns the new.
//
// Usage:
//
//	tbl := table.Clone()
//  defer tbl.Release()
//  ...
//	s.table = Swap(s.table, tbl)
//
func Swap[C Count](old C, new C) C {
	new.Retain()
	old.Release()
	return new
}

// SwapNoRetain releases an old reference, and returns a new one.
//
// Usage:
//
//	tbl := newTable()
//  ...
//	s.table = SwapNoRetain(s.table, tbl)
//
func SwapNoRetain[C Count](old C, new C) C {
	old.Release()
	return new
}
