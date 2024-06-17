package refstream

import (
	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/async/streams"
	"github.com/basecomplextech/baselibrary/ref"
	"github.com/basecomplextech/baselibrary/ref/refqueue"
)

// Stream is a wrapper around async.Stream with reference counting.
type Stream[T any] interface {
	streams.Stream[ref.R[T]]
}

// Source is a wrapper around async.Source with reference counting.
type Source[T any] interface {
	streams.Source[ref.R[T]]
}

// Stream is a wrapper around async.Queue with reference counting.
type Queue[T any] interface {
	async.Queue[ref.R[T]]
}

// NewSource returns a new stream source with reference counting.
func NewSource[T any]() Source[T] {
	src := streams.NewSource[ref.R[T]]()
	return newSource(src)
}

// Map returns a stream which maps messages from another stream.
func Map[T, R any](s Stream[T], fn func(ref.R[T]) ref.R[R]) Stream[R] {
	s1 := streams.Map(s, fn)
	return newStream(s1)
}

// source

var _ Source[any] = (*source[any])(nil)

type source[T any] struct {
	src streams.Source[ref.R[T]]
}

func newSource[T any](src streams.Source[ref.R[T]]) *source[T] {
	return &source[T]{src}
}

// Filter returns a new stream that only contains elements that satisfy the predicate.
func (s *source[T]) Filter(fn func(ref.R[T]) bool) streams.Stream[ref.R[T]] {
	src1 := s.src.Filter(fn)
	return newStream(src1)
}

// Listen adds a listener to the stream, and returns an unsubscribe function.
func (s *source[T]) Listen(ln streams.Listener[ref.R[T]]) (unsub func()) {
	return s.src.Listen(ln)
}

// Subscribe returns a queue subscribed to the stream, free the queue to unsubscribe.
func (s *source[T]) Subscribe() async.Queue[ref.R[T]] {
	q := newQueue[T]()
	q.unsub = s.src.Listen(q)
	return q
}

// Send sends a message to the stream.
func (s *source[T]) Send(msg ref.R[T]) {
	s.src.Send(msg)
}

// stream

var _ Stream[any] = (*stream[any])(nil)

type stream[T any] struct {
	src streams.Stream[ref.R[T]]
}

func newStream[T any](src streams.Stream[ref.R[T]]) *stream[T] {
	return &stream[T]{src}
}

// Filter returns a new stream that only contains elements that satisfy the predicate.
func (s *stream[T]) Filter(fn func(ref.R[T]) bool) streams.Stream[ref.R[T]] {
	src1 := s.src.Filter(fn)
	return newStream(src1)
}

// Listen adds a listener to the stream, and returns an unsubscribe function.
func (s *stream[T]) Listen(ln streams.Listener[ref.R[T]]) (unsub func()) {
	return s.src.Listen(ln)
}

// Subscribe returns a queue subscribed to the stream, free the queue to unsubscribe.
func (s *stream[T]) Subscribe() async.Queue[ref.R[T]] {
	q := newQueue[T]()
	q.unsub = s.src.Listen(q)
	return q
}

// queue

var _ streams.Listener[ref.R[any]] = (*queue[any])(nil)

type queue[T any] struct {
	refqueue.Queue[T]
	unsub func()
}

func newQueue[T any]() *queue[T] {
	return &queue[T]{
		Queue: refqueue.New[T](),
	}
}

// OnStreamMessage is called when a new message is available on the stream.
func (q *queue[T]) OnStreamMessage(msg ref.R[T]) {
	q.Queue.Push(msg)
}
