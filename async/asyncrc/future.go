// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncrc

import "github.com/basecomplextech/baselibrary/async"

// Future is a reference counted future.
type Future[T any] interface {
	async.Future[T]

	// Retain increments the reference count, panics if the future is already released.
	Retain()

	// Release decrements the reference count, frees the future if it reaches zero.
	Release()
}
