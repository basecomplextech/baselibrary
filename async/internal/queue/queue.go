// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package queue

import (
	"sync"

	"github.com/basecomplextech/baselibrary/collect/chans"
	"github.com/basecomplextech/baselibrary/collect/slices2"
)

// Queue is an unbounded FIFO queue.
type Queue[T any] interface {
	// Len returns the number of elements in the queue.
	Len() int

	// Clear clears the queue.
	Clear()

	// Push adds an element to the queue, panics if the queue is closed.
	Push(v T)

	// Poll removes an element from the queue, returns false if the queue is empty.
	Poll() (T, bool)

	// Wait returns a channel which is notified on new elements.
	Wait() <-chan struct{}
}

// New returns an empty queue.
func New[T any](items ...T) Queue[T] {
	q := newQueue[T]()

	if len(items) > 0 {
		q.list = append(q.list, items...)
	}
	return q
}

// internal

var _ Queue[int] = (*queue[int])(nil)

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

// Read

// Len returns the number of elements in the queue.
func (q *queue[T]) Len() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.list)
}

// Poll removes an element from the queue, returns false if the queue is empty.
func (q *queue[T]) Poll() (v T, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.assertOpen()

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
	q.assertOpen()

	if len(q.list) > 0 {
		return chans.Closed()
	}
	return q.wait
}

// Write

// Clear clears the queue.
func (q *queue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.assertOpen()

	q.list = slices2.Truncate(q.list)
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
	q.assertOpen()

	q.list = append(q.list, v)
	select {
	case q.wait <- struct{}{}:
	default:
	}
}

// private

func (q *queue[T]) assertOpen() {
	if q.closed {
		panic("queue is closed")
	}
}
