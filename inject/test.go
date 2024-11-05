// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package inject

import (
	"github.com/basecomplextech/baselibrary/tests"
)

// Test returns a test context.
func Test(t tests.T, vv ...any) Context {
	tfn := func() tests.T { return t } // returns tests.T interface

	x := New(t, tfn)
	x.Add(vv...)
	return x
}
