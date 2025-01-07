// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package collect

import (
	"container/heap"

	"github.com/basecomplextech/baselibrary/collect/slices2"
	"github.com/basecomplextech/baselibrary/compare"
	"golang.org/x/exp/constraints"
)

// UniquePriorityQueue is a priority queue that contains unique elements.
type UniquePriorityQueue[T comparable, P any] interface {
	// Len returns the number of elements in the queue.
	Len() int

	// Clear removes all elements from the queue.
	Clear()

	// Contains returns true if the queue contains an element.
	Contains(value T) bool

	// Get returns the priority of an element.
	Get(value T) (p P, ok bool)

	// Push pushes an element onto the queue, or updates its priority.
	Push(value T, priority P)

	// Poll removes and returns the minimum element (according to Less) from the queue.
	Poll() (value T, priority P, ok bool)

	// Remove removes an element from the queue, and returns its priority.
	Remove(value T) (p P, ok bool)
}

// NewUniquePriorityQueue returns a unique priority queue.
func NewUniquePriorityQueue[T comparable, P constraints.Ordered](
	items ...PriorityQueueItem[T, P]) *uniquePriorityQueue[T, P] {

	compare := compare.Ordered[P]()
	return newUniquePriorityQueue[T, P](compare, items...)
}

// internal

var _ UniquePriorityQueue[int, int] = (*uniquePriorityQueue[int, int])(nil)

// uniquePriorityQueue is a priority queue that contains unique elements.
type uniquePriorityQueue[T comparable, P any] struct {
	heap *uniqueHeap[T, P]
}

func newUniquePriorityQueue[T comparable, P any](compare compare.Compare[P],
	items ...PriorityQueueItem[T, P]) *uniquePriorityQueue[T, P] {

	q := &uniquePriorityQueue[T, P]{
		heap: &uniqueHeap[T, P]{
			compare: compare,
			items:   make([]uniqueItem[T, P], 0, len(items)),
			indexes: make(map[T]int),
		},
	}

	for i, item := range items {
		item := uniqueItem[T, P]{
			value:    item.Value,
			priority: item.Priority,
			index:    i,
		}

		q.heap.items = append(q.heap.items, item)
		q.heap.indexes[item.value] = i
	}

	heap.Init(q.heap)
	return q
}

// Len returns the number of elements in the queue.
func (q *uniquePriorityQueue[T, P]) Len() int {
	return q.heap.Len()
}

// Clear removes all elements from the queue.
func (q *uniquePriorityQueue[T, P]) Clear() {
	q.heap.items = slices2.Truncate(q.heap.items)
	clear(q.heap.indexes)
}

// Contains returns true if the queue contains an element.
func (q *uniquePriorityQueue[T, P]) Contains(value T) bool {
	_, ok := q.heap.indexOf(value)
	return ok
}

// Get returns the priority of an element.
func (q *uniquePriorityQueue[T, P]) Get(value T) (p P, ok bool) {
	i, ok := q.heap.indexOf(value)
	if !ok {
		return p, false
	}

	return q.heap.items[i].priority, true
}

// Push pushes an element onto the queue, or updates its priority.
func (q *uniquePriorityQueue[T, P]) Push(value T, priority P) {
	i, ok := q.heap.indexOf(value)
	if ok {
		q.heap.update(i, priority)
		return
	}

	item := uniqueItem[T, P]{
		value:    value,
		priority: priority,
	}
	heap.Push(q.heap, item)
}

// Poll removes and returns the minimum element (according to Less) from the queue.
func (q *uniquePriorityQueue[T, P]) Poll() (value T, priority P, ok bool) {
	if len(q.heap.items) == 0 {
		return value, priority, false
	}

	item := heap.Pop(q.heap).(uniqueItem[T, P])
	return item.value, item.priority, true
}

// Remove removes an element from the queue, and returns its priority.
func (q *uniquePriorityQueue[T, P]) Remove(value T) (p P, ok bool) {
	i, ok := q.heap.indexOf(value)
	if !ok {
		return p, false
	}

	item := heap.Remove(q.heap, i).(uniqueItem[T, P])
	return item.priority, true
}
