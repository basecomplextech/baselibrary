package async

var _ Stream[any] = (*streamFilter[any])(nil)

type streamFilter[T any] struct {
	src Stream[T]
	fn  func(T) bool
}

func newStreamFilter[T any](src Stream[T], fn func(T) bool) *streamFilter[T] {
	return &streamFilter[T]{
		src: src,
		fn:  fn,
	}
}

// Filter returns a new stream that only contains elements that satisfy the predicate.
func (s *streamFilter[T]) Filter(fn func(T) bool) Stream[T] {
	return &streamFilter[T]{
		src: s,
		fn:  fn,
	}
}

// Listen adds a listener to the stream, and returns an unsubscribe function.
func (s *streamFilter[T]) Listen(ln StreamListener[T]) (unsub func()) {
	filter := newStreamFilterListener(s.fn, ln)
	return s.src.Listen(filter)
}

// Subscribe returns a queue subscribed to the stream, free the queue to unsubscribe.
func (s *streamFilter[T]) Subscribe() Queue[T] {
	queue := newStreamQueue[T]()
	filter := newStreamFilterListener(s.fn, queue)
	queue.unsub = s.src.Listen(filter)
	return queue
}

// listener

var _ StreamListener[any] = (*streamFilterListener[any])(nil)

type streamFilterListener[T any] struct {
	fn   func(T) bool
	next StreamListener[T]
}

func newStreamFilterListener[T any](fn func(T) bool, next StreamListener[T]) *streamFilterListener[T] {
	return &streamFilterListener[T]{
		fn:   fn,
		next: next,
	}
}

// OnStreamMessage is called when a new message is available on the stream.
func (l *streamFilterListener[T]) OnStreamMessage(msg T) {
	if l.fn(msg) {
		l.next.OnStreamMessage(msg)
	}
}
