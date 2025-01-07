// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package collect

import (
	"container/heap"
	"slices"

	"github.com/basecomplextech/baselibrary/compare"
	"golang.org/x/exp/constraints"
)

// PriorityQueue is a priority queue.
type PriorityQueue[T any, P any] interface {
	// Len returns the number of elements in the queue.
	Len() int

	// Clear removes all elements from the queue.
	Clear()

	// Poll removes and returns the minimum element (according to Less) from the queue.
	Poll() (value T, priority P, ok bool)

	// Push pushes an element onto the queue.
	Push(value T, priority P)
}

// PriorityQueueItem is an element of the queue.
type PriorityQueueItem[T any, P any] struct {
	Value    T
	Priority P
}

// NewPriorityQueue returns a new priority queue with an ordered priority.
func NewPriorityQueue[T any, P constraints.Ordered](
	items ...PriorityQueueItem[T, P]) PriorityQueue[T, P] {

	compare := compare.Ordered[P]()
	return newPriorityQueue[T, P](compare, items...)
}

// NewPriorityQueueCompare returns a new priority queue with a priority compare function.
func NewPriorityQueueCompare[T any, P any](compare compare.Compare[P],
	items ...PriorityQueueItem[T, P]) PriorityQueue[T, P] {

	return newPriorityQueue[T, P](compare, items...)
}

// internal

// priorityQueue is a priority queue that is implemented using a heap.
type priorityQueue[T any, P any] struct {
	heap *priorityHeap[T, P]
}

// newPriorityQueue returns a new priority queue with a priority compare function.
func newPriorityQueue[T any, P any](compare compare.Compare[P],
	items ...PriorityQueueItem[T, P]) *priorityQueue[T, P] {

	q := &priorityQueue[T, P]{
		heap: &priorityHeap[T, P]{
			compare: compare,
			items:   slices.Clone(items),
		},
	}

	heap.Init(q.heap)
	return q
}

// Len returns the number of elements in the queue.
func (q *priorityQueue[T, P]) Len() int {
	return q.heap.Len()
}

// Clear removes all elements from the queue.
func (q *priorityQueue[T, P]) Clear() {
	q.heap.items = nil
}

// Push pushes an element onto the queue.
func (q *priorityQueue[T, P]) Push(value T, priority P) {
	item := PriorityQueueItem[T, P]{
		Value:    value,
		Priority: priority,
	}

	heap.Push(q.heap, item)
}

// Poll removes and returns the minimum element (according to Less) from the queue.
func (q *priorityQueue[T, P]) Poll() (value T, priority P, ok bool) {
	if q.Len() == 0 {
		return value, priority, false
	}

	item := heap.Pop(q.heap).(PriorityQueueItem[T, P])
	return item.Value, item.Priority, true
}
