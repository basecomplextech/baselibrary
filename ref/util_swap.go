// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

// SwapFree frees an old object and returns the new.
//
// Usage:
//
//	tbl := table.Clone()
//	...
//	s.table = SwapFree(s.table, tbl)
func SwapFree[T Freer](old T, new T) T {
	old.Free()
	return new
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
