package streams

import (
	"sync"

	"github.com/basecomplextech/baselibrary/async"
)

// Source is a stream which can be written to.
type Source[T any] interface {
	Stream[T]

	// Send sends a message to the stream.
	Send(T)
}

// NewSource returns a new stream source.
func NewSource[T any]() Source[T] {
	return newSource[T]()
}

// internal

var _ Source[any] = (*streamSource[any])(nil)

type streamSource[T any] struct {
	mu sync.Mutex
	i  int
	ln map[int]Listener[T]
}

func newSource[T any]() *streamSource[T] {
	return &streamSource[T]{
		ln: make(map[int]Listener[T]),
	}
}

// Filter returns a new stream that only contains elements that satisfy the predicate.
func (s *streamSource[T]) Filter(fn func(T) bool) Stream[T] {
	return newStreamFilter(s, fn)
}

// Listen adds a listener to the stream, and returns an unsubscribe function.
func (s *streamSource[T]) Listen(ln Listener[T]) (unsub func()) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.i++
	seq := s.i
	s.ln[seq] = ln

	unsub = func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		delete(s.ln, seq)
	}
	return unsub
}

// Subscribe returns a queue subscribed to the stream, free the queue to unsubscribe.
func (s *streamSource[T]) Subscribe() async.Queue[T] {
	queue := newQueue[T]()
	queue.unsub = s.Listen(queue)
	return queue
}

// Send sends a message to the stream.
func (s *streamSource[T]) Send(msg T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, ln := range s.ln {
		ln.OnStreamMessage(msg)
	}
}
