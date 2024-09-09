// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package streams

import "github.com/basecomplextech/baselibrary/async"

var _ Stream[any] = (*filterStream[any])(nil)

type filterStream[T any] struct {
	src Stream[T]
	fn  func(T) bool
}

func newStreamFilter[T any](src Stream[T], fn func(T) bool) *filterStream[T] {
	return &filterStream[T]{
		src: src,
		fn:  fn,
	}
}

// Filter returns a new stream that only contains elements that satisfy the predicate.
func (s *filterStream[T]) Filter(fn func(T) bool) Stream[T] {
	return &filterStream[T]{
		src: s,
		fn:  fn,
	}
}

// Listen adds a listener to the stream, and returns an unsubscribe function.
func (s *filterStream[T]) Listen(ln Listener[T]) (unsub func()) {
	filter := newFilterListener(s.fn, ln)
	return s.src.Listen(filter)
}

// Subscribe returns a queue subscribed to the stream, free the queue to unsubscribe.
func (s *filterStream[T]) Subscribe() async.Queue[T] {
	queue := newQueue[T]()
	filter := newFilterListener(s.fn, queue)
	queue.unsub = s.src.Listen(filter)
	return queue
}

// listener

var _ Listener[any] = (*filterListener[any])(nil)

type filterListener[T any] struct {
	fn   func(T) bool
	next Listener[T]
}

func newFilterListener[T any](fn func(T) bool, next Listener[T]) *filterListener[T] {
	return &filterListener[T]{
		fn:   fn,
		next: next,
	}
}

// OnStreamMessage is called when a new message is available on the stream.
func (l *filterListener[T]) OnStreamMessage(msg T) {
	if l.fn(msg) {
		l.next.OnStreamMessage(msg)
	}
}
