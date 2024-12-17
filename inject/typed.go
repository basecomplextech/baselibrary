// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package inject

// Typed returns a function that returns the object passed to it.
// This allows to pass an object as an interface.
//
// Usage:
//
//	inject.New(inject.Typed(obj))
func Typed[T any](obj T) func() T {
	return func() T { return obj }
}
