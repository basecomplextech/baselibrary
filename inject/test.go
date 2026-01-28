// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package inject

import (
	"github.com/basecomplextech/baselibrary/tests"
)

// Test returns a test context which eagerly initializes all dependencies.
func Test(t tests.T, providers ...any) Context {
	x := TestLazy(t, providers...).(*context)

	for typ := range x.providers {
		x.get(typ)
	}
	return x
}

// TestLazy returns a lazy test context.
func TestLazy(t tests.T, providers ...any) Context {
	// testFn returns tests.T interface
	testFn := func() tests.T { return t }

	x := New(t, testFn)
	x.Add(providers...)
	return x
}
