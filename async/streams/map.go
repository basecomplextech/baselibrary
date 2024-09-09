// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package streams

import "github.com/basecomplextech/baselibrary/async"

// Map returns a stream which maps messages from another stream.
func Map[T, R any](s Stream[T], fn func(T) R) Stream[R] {
	return newMapStream(s, fn)
}

// internal

var _ Stream[any] = (*mapStream[any, any])(nil)

type mapStream[T, R any] struct {
	src Stream[T]
	fn  func(T) R
}

func newMapStream[T, R any](src Stream[T], fn func(T) R) *mapStream[T, R] {
	return &mapStream[T, R]{
		src: src,
		fn:  fn,
	}
}

// Filter returns a new stream that only contains elements that satisfy the predicate.
func (s *mapStream[T, R]) Filter(fn func(R) bool) Stream[R] {
	return newStreamFilter(s, fn)
}

// Listen adds a listener to the stream, and returns an unsubscribe function.
func (s *mapStream[T, R]) Listen(ln Listener[R]) (unsub func()) {
	mp := newMapListener(s.fn, ln)
	return s.src.Listen(mp)
}

// Subscribe returns a queue subscribed to the stream, free the queue to unsubscribe.
func (s *mapStream[T, R]) Subscribe() async.Queue[R] {
	q := newQueue[R]()
	mp := newMapListener(s.fn, q)
	q.unsub = s.src.Listen(mp)
	return q
}

// listener

var _ Listener[any] = (*mapListener[any, any])(nil)

type mapListener[T, R any] struct {
	fn   func(T) R
	next Listener[R]
}

func newMapListener[T, R any](fn func(T) R, next Listener[R]) *mapListener[T, R] {
	return &mapListener[T, R]{
		fn:   fn,
		next: next,
	}
}

// OnStreamMessage is called when a new message is available on the stream.
func (l *mapListener[T, R]) OnStreamMessage(msg T) {
	next := l.fn(msg)
	l.next.OnStreamMessage(next)
}
