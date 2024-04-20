package refstream

import (
	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/ref"
	"github.com/basecomplextech/baselibrary/ref/refqueue"
)

// Stream is a wrapper around async.Stream with reference counting.
type Stream[T ref.Ref] interface {
	async.Stream[T]
}

// Source is a wrapper around async.StreamSource with reference counting.
type Source[T ref.Ref] interface {
	async.StreamSource[T]
}

// Stream is a wrapper around async.Queue with reference counting.
type Queue[T ref.Ref] interface {
	refqueue.Queue[T]
}

// NewSource returns a new stream source with reference counting.
func NewSource[T ref.Ref]() async.StreamSource[T] {
	src := async.NewStreamSource[T]()
	return newSource[T](src)
}

// Map returns a stream which maps messages from another stream.
func Map[T, R ref.Ref](s Stream[T], fn func(T) R) Stream[R] {
	s1 := async.MapStream[T, R](s, fn)
	return newStream(s1)
}

// source

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

// stream

var _ Stream[ref.Ref] = (*stream[ref.Ref])(nil)

type stream[T ref.Ref] struct {
	src async.Stream[T]
}

func newStream[T ref.Ref](src async.Stream[T]) *stream[T] {
	return &stream[T]{src}
}

// Filter returns a new stream that only contains elements that satisfy the predicate.
func (s *stream[T]) Filter(fn func(T) bool) async.Stream[T] {
	src1 := s.src.Filter(fn)
	return newStream(src1)
}

// Listen adds a listener to the stream, and returns an unsubscribe function.
func (s *stream[T]) Listen(ln async.StreamListener[T]) (unsub func()) {
	return s.src.Listen(ln)
}

// Subscribe returns a queue subscribed to the stream, free the queue to unsubscribe.
func (s *stream[T]) Subscribe() async.Queue[T] {
	q := newQueue[T]()
	q.unsub = s.src.Listen(q)
	return q
}

// queue

var _ async.StreamListener[ref.Ref] = (*queue[ref.Ref])(nil)

type queue[T ref.Ref] struct {
	refqueue.Queue[T]
	unsub func()
}

func newQueue[T ref.Ref]() *queue[T] {
	return &queue[T]{
		Queue: refqueue.New[T](),
	}
}

// OnStreamMessage is called when a new message is available on the stream.
func (q *queue[T]) OnStreamMessage(msg T) {
	q.Queue.Push(msg)
}
