// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package priorityqueue

import (
	"sync"

	"github.com/basecomplextech/baselibrary/collect"
	"github.com/basecomplextech/baselibrary/collect/chans"
	"github.com/basecomplextech/baselibrary/compare"
	"golang.org/x/exp/constraints"
)

// Queue is a priority queue guarded by a mutex.
type Queue[T any, P any] interface {
	collect.PriorityQueue[T, P]

	// Wait returns a channel that is notified on new items.
	Wait() <-chan struct{}
}

// New returns a new priority queue with an ordered priority.
func New[T any, P constraints.Ordered]() Queue[T, P] {
	return newOrdered[T, P]()
}

// NewCompare returns a new priority queue with a priority compare function.
func NewCompare[T any, P any](compare compare.Compare[P]) Queue[T, P] {
	return newQueue[T, P](compare)
}

// internal

type queue[T any, P any] struct {
	mu      sync.RWMutex
	collect collect.PriorityQueue[T, P]
	wait    chan struct{}
}

func newOrdered[T any, P constraints.Ordered]() *queue[T, P] {
	return &queue[T, P]{
		collect: collect.NewPriorityQueue[T, P](),
		wait:    make(chan struct{}, 1),
	}
}

func newQueue[T any, P any](compare compare.Compare[P]) *queue[T, P] {
	return &queue[T, P]{
		collect: collect.NewPriorityQueueCompare[T, P](compare),
		wait:    make(chan struct{}, 1),
	}
}

// Len returns the number of elements in the queue.
func (q *queue[T, P]) Len() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.collect.Len()
}

// Clear removes all elements from the queue.
func (q *queue[T, P]) Clear() {
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

// Push pushes an element onto the queue.
func (q *queue[T, P]) Push(value T, priority P) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Push item
	q.collect.Push(value, priority)

	// Notify waiter
	select {
	case q.wait <- struct{}{}:
	default:
	}
}

// Poll removes and returns the minimum element (according to Less) from the queue.
func (q *queue[T, P]) Poll() (value T, priority P, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.collect.Poll()
}

// Wait returns a channel that is notified on new items.
func (q *queue[T, P]) Wait() <-chan struct{} {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.collect.Len() > 0 {
		return chans.Closed()
	}
	return q.wait
}
