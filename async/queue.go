package async

import (
	"sync"

	"github.com/basecomplextech/baselibrary/collect/slices2"
)

// Queue is an unbounded FIFO queue.
type Queue[T any] interface {
	// Clear clears the queue.
	Clear()

	// Len returns the number of elements in the queue.
	Len() int

	// Push adds an element to the queue, panics if the queue is freed.
	Push(v T)

	// Pop removes an element from the queue, returns false if the queue is empty.
	Pop() (v T, ok bool)

	// Wait returns a channel which is notified on new elements.
	Wait() <-chan struct{}

	// Free clears and closes the queue.
	Free()
}

// NewQueue returns an empty queue.
func NewQueue[T any](items ...T) Queue[T] {
	q := newQueue[T]()

	if len(items) > 0 {
		q.list = append(q.list, items...)
	}
	return q
}

// internal

type queue[T any] struct {
	mu    sync.RWMutex
	freed bool

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
	q.assertNotFreed()

	q.list = slices2.Truncate(q.list)
}

// Len returns the number of elements in the queue.
func (q *queue[T]) Len() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.list)
}

// Push adds an element to the queue, panics if the queue is freed.
func (q *queue[T]) Push(v T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.assertNotFreed()

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
	q.assertNotFreed()

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
	q.mu.RLock()
	defer q.mu.RUnlock()
	q.assertNotFreed()

	if len(q.list) > 0 {
		return closedChan
	}
	return q.wait
}

// Free clears and closes the queue.
func (q *queue[T]) Free() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.freed {
		return
	}

	q.list = nil
	q.freed = true
	close(q.wait)
}

// private

func (q *queue[T]) assertNotFreed() {
	if q.freed {
		panic("operation on freed queue")
	}
}

var closedChan = func() chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}()
