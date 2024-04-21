package async

var _ Stream[any] = (*streamMap[any, any])(nil)

type streamMap[T, R any] struct {
	src Stream[T]
	fn  func(T) R
}

func newStreamMap[T, R any](src Stream[T], fn func(T) R) *streamMap[T, R] {
	return &streamMap[T, R]{
		src: src,
		fn:  fn,
	}
}

// Filter returns a new stream that only contains elements that satisfy the predicate.
func (s *streamMap[T, R]) Filter(fn func(R) bool) Stream[R] {
	return newStreamFilter(s, fn)
}

// Listen adds a listener to the stream, and returns an unsubscribe function.
func (s *streamMap[T, R]) Listen(ln StreamListener[R]) (unsub func()) {
	mp := newStreamMapListener(s.fn, ln)
	return s.src.Listen(mp)
}

// Subscribe returns a queue subscribed to the stream, free the queue to unsubscribe.
func (s *streamMap[T, R]) Subscribe() Queue[R] {
	q := newStreamQueue[R]()
	mp := newStreamMapListener(s.fn, q)
	q.unsub = s.src.Listen(mp)
	return q
}

// listener

var _ StreamListener[any] = (*streamMapListener[any, any])(nil)

type streamMapListener[T, R any] struct {
	fn   func(T) R
	next StreamListener[R]
}

func newStreamMapListener[T, R any](fn func(T) R, next StreamListener[R]) *streamMapListener[T, R] {
	return &streamMapListener[T, R]{
		fn:   fn,
		next: next,
	}
}

// OnStreamMessage is called when a new message is available on the stream.
func (l *streamMapListener[T, R]) OnStreamMessage(msg T) {
	next := l.fn(msg)
	l.next.OnStreamMessage(next)
}
