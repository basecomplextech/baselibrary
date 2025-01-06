// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package priorityqueue

import (
	"container/heap"
	"slices"

	"github.com/basecomplextech/baselibrary/compare"
	"golang.org/x/exp/constraints"
)

// Queue is a priority queue.
type Queue[T any, P any] interface {
	// Len returns the number of elements in the queue.
	Len() int

	// Clear removes all elements from the queue.
	Clear()

	// Poll removes and returns the minimum element (according to Less) from the queue.
	Poll() (value T, priority P, ok bool)

	// Push pushes an element onto the queue.
	Push(value T, priority P)
}

// Item is an element of the queue.
type Item[T any, P any] struct {
	Value    T
	Priority P
}

// New returns a new priority queue with an ordered priority.
func New[T any, P constraints.Ordered](items ...Item[T, P]) Queue[T, P] {
	compare := compare.Ordered[P]()
	return newQueue[T, P](compare, items...)
}

// NewCompare returns a new priority queue with a priority compare function.
func NewCompare[T any, P any](compare compare.Compare[P], items ...Item[T, P]) Queue[T, P] {
	return newQueue(compare, items...)
}

// internal

// queue is a priority queue that is implemented using a heap.
type queue[T any, P any] struct {
	heap *queueHeap[T, P]
}

// newQueue returns a new priority queue with a priority compare function.
func newQueue[T any, P any](compare compare.Compare[P], items ...Item[T, P]) *queue[T, P] {
	q := &queue[T, P]{
		heap: &queueHeap[T, P]{
			compare: compare,
			items:   slices.Clone(items),
		},
	}

	heap.Init(q.heap)
	return q
}

// Len returns the number of elements in the queue.
func (q *queue[T, P]) Len() int {
	return q.heap.Len()
}

// Clear removes all elements from the queue.
func (q *queue[T, P]) Clear() {
	q.heap.items = nil
}

// Push pushes an element onto the queue.
func (q *queue[T, P]) Push(value T, priority P) {
	item := Item[T, P]{
		Value:    value,
		Priority: priority,
	}

	heap.Push(q.heap, item)
}

// Poll removes and returns the minimum element (according to Less) from the queue.
func (q *queue[T, P]) Poll() (value T, priority P, ok bool) {
	if q.Len() == 0 {
		return value, priority, false
	}

	item := heap.Pop(q.heap).(Item[T, P])
	return item.Value, item.Priority, true
}
