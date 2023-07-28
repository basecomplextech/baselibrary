package async

import (
	"sync"

	"github.com/basecomplextech/baselibrary/collect/slices"
)

// Queue is an unbounded FIFO queue.
type Queue[T any] interface {
	// Clear clears the queue.
	Clear()

	// Close clears and closes the queue.
	Close()

	// Push adds an element to the queue, panics if the queue is closed.
	Push(v T)

	// Pop removes an element from the queue, returns false if the queue is empty.
	Pop() (v T, ok bool)

	// Wait returns a channel which is notified on new elements.
	Wait() <-chan struct{}
}

// NewQueue returns an empty queue.
func NewQueue[T any]() Queue[T] {
	return newQueue[T]()
}

// internal

type queue[T any] struct {
	mu     sync.RWMutex
	closed bool

	list []T
	wait chan struct{}
}

func newQueue[T any]() *queue[T] {
	return &queue[T]{
		wait: make(chan struct{}, 1),
	}
}

// Clear clears the queue.
func (q *queue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.list = slices.Clear(q.list)
}

// Close clears and closes the queue.
func (q *queue[T]) Close() {
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
func (q *queue[T]) Push(v T) {
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
func (q *queue[T]) Pop() (v T, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.list) == 0 {
		return
	}

	// Get value
	v = q.list[0]

	// Move remaining to front
	copy(q.list, q.list[1:])
	q.list = q.list[:len(q.list)-1]
	return v, true
}

// Wait returns a channel which is notified on new elements.
func (q *queue[T]) Wait() <-chan struct{} {
	return q.wait
}
