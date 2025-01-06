// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"github.com/basecomplextech/baselibrary/async/internal/stream"
	"github.com/basecomplextech/baselibrary/ref"
)

type (
	// Stream is an async stream of values.
	Stream[T any] = stream.Stream[T]

	// StreamNextFunc is a function that returns the next value from a stream.
	StreamNextFunc[T any] = stream.NextFunc[T]
)

// NewStream returns a new stream.
func NewStream[T any](next StreamNextFunc[T]) Stream[T] {
	return stream.New(next)
}

// NewStreamFree returns a new stream with a free function.
func NewStreamFree[T any](next StreamNextFunc[T], free func()) Stream[T] {
	return stream.NewFree(next, free)
}

// NewStreamFreer returns a new stream with a freer.
func NewStreamFreer[T any](next StreamNextFunc[T], freer ref.Freer) Stream[T] {
	return stream.NewFreer(next, freer)
}
