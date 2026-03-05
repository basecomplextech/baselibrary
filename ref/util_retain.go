// Copyright 2021 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import "github.com/basecomplextech/baselibrary/opt"

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

// RetainOpt retains a optional reference and returns it.
func RetainOpt[R Ref](r opt.Opt[R]) opt.Opt[R] {
	r1, ok := r.Unwrap()
	if !ok {
		return r
	}

	r1.Retain()
	return r
}
