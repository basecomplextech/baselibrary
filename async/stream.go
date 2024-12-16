// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"github.com/basecomplextech/baselibrary/opt"
	"github.com/basecomplextech/baselibrary/ref"
	"github.com/basecomplextech/baselibrary/status"
)

// Stream is an async stream of values.
type Stream[T any] interface {
	// Next returns the next value from the stream, or false if the stream has ended.
	Next(ctx Context) (T, bool, status.Status)

	// Internal

	// Free frees the stream.
	Free()
}

// StreamNextFunc is a function that returns the next value from a stream.
type StreamNextFunc[T any] func(ctx Context) (T, bool, status.Status)

// NewStream returns a new stream.
func NewStream[T any](next StreamNextFunc[T]) Stream[T] {
	return newStream(next, opt.None[ref.Freer]())
}

// NewStreamFree returns a new stream with a free function.
func NewStreamFree[T any](next StreamNextFunc[T], free func()) Stream[T] {
	freer := ref.FreeFunc(free)
	return NewStreamFreer(next, freer)
}

// NewStreamFreer returns a new stream with a freer.
func NewStreamFreer[T any](next StreamNextFunc[T], freer ref.Freer) Stream[T] {
	return &stream[T]{
		next:  next,
		freer: opt.New(freer),
	}
}

// internal

var _ Stream[any] = (*stream[any])(nil)

type stream[T any] struct {
	next  func(ctx Context) (T, bool, status.Status)
	freer opt.Opt[ref.Freer]
}

func newStream[T any](next StreamNextFunc[T], freer opt.Opt[ref.Freer]) Stream[T] {
	return &stream[T]{
		next:  next,
		freer: freer,
	}
}

// Next returns the next value from the stream, or false if the stream has ended.
func (s *stream[T]) Next(ctx Context) (T, bool, status.Status) {
	return s.next(ctx)
}

// Internal

// Free frees the stream.
func (s *stream[T]) Free() {
	if freer, ok := s.freer.Clear(); ok {
		freer.Free()
	}
}
