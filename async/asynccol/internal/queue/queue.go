// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package queue

import (
	"sync"

	"github.com/basecomplextech/baselibrary/collect"
	"github.com/basecomplextech/baselibrary/collect/chans"
)

// Queue is an unbounded FIFO queue guarded by a mutex.
type Queue[T any] interface {
	collect.Queue[T]

	// Wait returns a channel which is notified on new elements.
	Wait() <-chan struct{}
}

// New returns an empty queue.
func New[T any](items ...T) Queue[T] {
	return newQueue[T](items...)
}

// internal

var _ Queue[int] = (*queue[int])(nil)

type queue[T any] struct {
	mu      sync.RWMutex
	collect collect.Queue[T]
	wait    chan struct{}
}

func newQueue[T any](items ...T) *queue[T] {
	return &queue[T]{
		collect: collect.NewQueue[T](items...),
		wait:    make(chan struct{}, 1),
	}
}

// Read

// Len returns the number of elements in the queue.
func (q *queue[T]) Len() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.collect.Len()
}

// Clear clears the queue.
func (q *queue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Clear queue
	q.collect.Clear()

	// Notify waiters
	for len(q.wait) > 0 {
		select {
		case <-q.wait:
		default:
			break
		}
	}
}

// Push adds an element to the queue.
func (q *queue[T]) Push(v T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Push item
	q.collect.Push(v)

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

	return q.collect.Poll()
}

// Wait returns a channel which is notified on new elements.
func (q *queue[T]) Wait() <-chan struct{} {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.collect.Len() > 0 {
		return chans.Closed()
	}
	return q.wait
}
