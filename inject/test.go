// Copyright 2024 Ivan Korobkov. All rights reserved.

package inject

import (
	"github.com/basecomplextech/baselibrary/tests"
)

// Test returns a test context.
func Test(t tests.T) Context {
	return New().
		Add(t).
		Add(func() tests.T { return t }) // Add interface
}
