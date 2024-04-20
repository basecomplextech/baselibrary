package refstream

import (
	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/ref"
)

var _ Stream[ref.Ref] = (*streamMap[ref.Ref, ref.Ref])(nil)

type streamMap[T, R ref.Ref] struct {
	src Stream[T]
	fn  func(T) R
}

func newStreamMap[T, R ref.Ref](src Stream[T], fn func(T) R) *streamMap[T, R] {
	return &streamMap[T, R]{
		src: src,
		fn:  fn,
	}
}

// Filter returns a new stream that only contains elements that satisfy the predicate.
func (s *streamMap[T, R]) Filter(fn func(R) bool) Stream[R] {
	// return newStreamFilter[R](s, fn)
	panic("unsupported")
}

// Listen adds a listener to the stream, and returns an unsubscribe function.
func (s *streamMap[T, R]) Listen(ln async.StreamListener[R]) (unsub func()) {
	map_ := newStreamMapListener(s.fn, ln)
	return s.src.Listen(map_)
}

// Subscribe returns a queue subscribed to the stream, free the queue to unsubscribe.
func (s *streamMap[T, R]) Subscribe() Queue[R] {
	queue := newQueue[R]()
	map_ := newStreamMapListener(s.fn, queue)
	queue.unsub = s.src.Listen(map_)
	return queue
}

// listener

var _ async.StreamListener[ref.Ref] = (*streamMapListener[ref.Ref, ref.Ref])(nil)

type streamMapListener[T, R ref.Ref] struct {
	fn   func(T) R
	next async.StreamListener[R]
}

func newStreamMapListener[T, R ref.Ref](fn func(T) R, next async.StreamListener[R]) *streamMapListener[T, R] {
	return &streamMapListener[T, R]{
		fn:   fn,
		next: next,
	}
}

// OnStreamMessage is called when a new message is available on the stream.
func (l *streamMapListener[T, R]) OnStreamMessage(msg T) {
	next := l.fn(msg)
	defer next.Release()

	l.next.OnStreamMessage(next)
}
