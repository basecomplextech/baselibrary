// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import "github.com/basecomplextech/baselibrary/async/internal/variable"

// [Experimental] Variable is an asynchronous variable which can be set, cleared, or failed.
type Variable[T any] = variable.Variable[T]

// New returns a new pending async variable.
func New[T any]() variable.Variable[T] {
	return variable.New[T]()
}
