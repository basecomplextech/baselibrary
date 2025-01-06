// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package queue

import (
	"github.com/basecomplextech/baselibrary/collect/slices2"
)

// Queue is an unbounded FIFO queue.
type Queue[T any] interface {
	// Len returns the number of elements in the queue.
	Len() int

	// Clear clears the queue.
	Clear()

	// Push adds an element to the queue.
	Push(v T)

	// Poll removes an element from the queue, returns false if the queue is empty.
	Poll() (T, bool)
}

// New returns a new queue.
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
	return len(q.list)
}

// Clear clears the queue.
func (q *queue[T]) Clear() {
	q.list = slices2.Truncate(q.list)

	// Notify all
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
	q.list = append(q.list, v)
}

// Poll removes an element from the queue, returns false if the queue is empty.
func (q *queue[T]) Poll() (v T, ok bool) {
	if len(q.list) == 0 {
		return
	}

	// Get value
	v = q.list[0]

	// Shift remaining left
	q.list = slices2.ShiftLeft(q.list, 1)
	return v, true
}
