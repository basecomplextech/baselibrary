package refstream

import (
	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/ref"
)

var _ Source[ref.Ref] = (*source[ref.Ref])(nil)

type source[T ref.Ref] struct {
	src async.StreamSource[T]
}

func newSource[T ref.Ref](src async.StreamSource[T]) *source[T] {
	return &source[T]{src}
}

// Filter returns a new stream that only contains elements that satisfy the predicate.
func (s *source[T]) Filter(fn func(T) bool) async.Stream[T] {
	src1 := s.src.Filter(fn)
	return newStream(src1)
}

// Listen adds a listener to the stream, and returns an unsubscribe function.
func (s *source[T]) Listen(ln async.StreamListener[T]) (unsub func()) {
	return s.src.Listen(ln)
}

// Subscribe returns a queue subscribed to the stream, free the queue to unsubscribe.
func (s *source[T]) Subscribe() async.Queue[T] {
	q := newQueue[T]()
	q.unsub = s.src.Listen(q)
	return q
}

// Send sends a message to the stream.
func (s *source[T]) Send(msg T) {
	s.src.Send(msg)
}
