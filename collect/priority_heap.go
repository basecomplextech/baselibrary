// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package collect

import (
	"container/heap"

	"github.com/basecomplextech/baselibrary/compare"
)

var _ heap.Interface = (*priorityHeap[any, any])(nil)

type priorityHeap[T any, P any] struct {
	compare compare.Compare[P]
	items   []PriorityQueueItem[T, P]
}

// Len is the number of elements in the collection.
func (q *priorityHeap[T, P]) Len() int {
	return len(q.items)
}

// Less reports whether the element with index i must sort before the element with index j.
func (q *priorityHeap[T, P]) Less(i, j int) bool {
	a, b := q.items[i], q.items[j]
	return q.compare(a.Priority, b.Priority) < 0
}

// Swap swaps the elements with indexes i and j.
func (q *priorityHeap[T, P]) Swap(i, j int) {
	q.items[i], q.items[j] = q.items[j], q.items[i]
}

// Push adds x as element Len()
func (q *priorityHeap[T, P]) Push(x any) {
	item := x.(PriorityQueueItem[T, P])
	q.items = append(q.items, item)
}

// Pop removes and return element Len() - 1.
func (q *priorityHeap[T, P]) Pop() any {
	n := len(q.items)
	x := q.items[n-1]
	q.items = q.items[:n-1]
	return x
}
