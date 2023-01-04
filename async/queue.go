package async

import "sync"

// Queue is a FIFO queue.
type Queue[T any] struct {
	mu     sync.RWMutex
	closed bool

	list []T
	wait chan struct{}
}

// NewQueue returns an empty queue.
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		wait: make(chan struct{}, 1),
	}
}

// Clear clears the queue.
func (q *Queue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.list = nil
}

// Close clears and closes the queue.
func (q *Queue[T]) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return
	}

	q.list = nil
	q.closed = true
	close(q.wait)
}

// Push adds an element to the queue, panics if the queue is closed.
func (q *Queue[T]) Push(v T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		panic("operation on closed queue")
	}

	q.list = append(q.list, v)
	select {
	case q.wait <- struct{}{}:
	default:
	}
}

// Pop removes an element from the queue, returns false if the queue is empty.
func (q *Queue[T]) Pop() (v T, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.list) == 0 {
		return
	}

	// get value
	v = q.list[0]

	// move remaining to front
	copy(q.list, q.list[1:])
	q.list = q.list[:len(q.list)-1]
	return v, true
}

// Wait returns a channel which is notified on new elements.
func (q *Queue[T]) Wait() <-chan struct{} {
	return q.wait
}
