// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"slices"
	"sync"

	"github.com/basecomplextech/baselibrary/collect/chans"
	"github.com/basecomplextech/baselibrary/collect/slices2"
)

// Queue is an unbounded FIFO queue guarded by a mutex.
type Queue[T any] interface {
	// Len returns the number of elements in the queue.
	Len() int

	// Clear clears the queue.
	Clear()

	// Push adds an element to the queue.
	Push(v T)

	// Poll removes an element from the queue, returns false if the queue is empty.
	Poll() (T, bool)

	// Wait returns a channel which is notified on new elements.
	Wait() <-chan struct{}
}

// NewQueue returns a new queue.
func NewQueue[T any](items ...T) Queue[T] {
	return newQueue[T](items...)
}

// DrainQueue drains all events from the topic queue.
func DrainQueue[T any](q Queue[T]) {
	for {
		_, ok := q.Poll()
		if !ok {
			return
		}
	}
}

// internal

var _ Queue[int] = (*queue[int])(nil)

type queue[T any] struct {
	mu   sync.RWMutex
	list []T
	wait chan struct{}
}

func newQueue[T any](items ...T) *queue[T] {
	return &queue[T]{
		list: slices.Clone(items),
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

// Clear clears the queue.
func (q *queue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Clear queue
	q.list = slices2.Truncate(q.list)

	// Notify waiters
	for len(q.wait) > 0 {
		select {
		case <-q.wait:
		default:
		}
	}
}

// Push adds an element to the queue.
func (q *queue[T]) Push(v T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Push item
	q.list = append(q.list, v)

	// Notify waiter
	select {
	case q.wait <- struct{}{}:
	default:
	}
}

// Poll removes an element from the queue, returns false if the queue is empty.
func (q *queue[T]) Poll() (v T, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.list) == 0 {
		return
	}

	// Get value
	v = q.list[0]

	// Shift remaining left
	q.list = slices2.ShiftLeft(q.list, 1)
	return v, true
}

// Wait returns a channel which is notified on new elements.
func (q *queue[T]) Wait() <-chan struct{} {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if len(q.list) > 0 {
		return chans.Closed()
	}
	return q.wait
}
