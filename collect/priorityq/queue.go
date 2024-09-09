// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package priorityq

import (
	"container/heap"
	"slices"
)

// Queue is a priority queue that is implemented using a heap.
type Queue[V any, P any] struct {
	queue *heapq[V, P]
}

// Item is an element of the queue.
type Item[V any, P any] struct {
	Value    V
	Priority P
}

// CompareFunc compares two priorities, and returns -1 if a < b, 0 if a == b, 1 if a > b.
type CompareFunc[P any] func(a, b P) int

// New returns a new priority queue with a priority compare function.
func New[V any, P any](compare CompareFunc[P], items ...Item[V, P]) *Queue[V, P] {
	q := &Queue[V, P]{
		queue: &heapq[V, P]{
			compare: compare,
			items:   slices.Clone(items),
		},
	}

	heap.Init(q.queue)
	return q
}

// Len returns the number of elements in the queue.
func (q *Queue[V, P]) Len() int {
	return q.queue.Len()
}

// Clear removes all elements from the queue.
func (q *Queue[V, P]) Clear() {
	q.queue.items = nil
}

// Pop removes and returns the minimum element (according to Less) from the queue.
func (q *Queue[V, P]) Pop() (value V, priority P, ok bool) {
	if q.Len() == 0 {
		return value, priority, false
	}

	item := heap.Pop(q.queue).(Item[V, P])
	return item.Value, item.Priority, true
}

// Push pushes an element onto the queue.
func (q *Queue[V, P]) Push(value V, priority P) {
	item := Item[V, P]{
		Value:    value,
		Priority: priority,
	}

	heap.Push(q.queue, item)
}

// internal

var _ heap.Interface = (*heapq[any, any])(nil)

type heapq[V any, P any] struct {
	compare CompareFunc[P]
	items   []Item[V, P]
}

// Len is the number of elements in the collection.
func (q *heapq[V, P]) Len() int {
	return len(q.items)
}

// Less reports whether the element with index i must sort before the element with index j.
func (q *heapq[V, P]) Less(i, j int) bool {
	a, b := q.items[i], q.items[j]
	return q.compare(a.Priority, b.Priority) < 0
}

// Swap swaps the elements with indexes i and j.
func (q *heapq[V, P]) Swap(i, j int) {
	q.items[i], q.items[j] = q.items[j], q.items[i]
}

// Push adds x as element Len()
func (q *heapq[V, P]) Push(x any) {
	item := x.(Item[V, P])
	q.items = append(q.items, item)
}

// Pop removes and return element Len() - 1.
func (q *heapq[V, P]) Pop() any {
	n := len(q.items)
	x := q.items[n-1]
	q.items = q.items[:n-1]
	return x
}
