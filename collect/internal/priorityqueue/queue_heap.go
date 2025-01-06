// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package priorityqueue

import (
	"container/heap"

	"github.com/basecomplextech/baselibrary/compare"
)

var _ heap.Interface = (*queueHeap[any, any])(nil)

type queueHeap[T any, P any] struct {
	compare compare.Compare[P]
	items   []Item[T, P]
}

// Len is the number of elements in the collection.
func (q *queueHeap[T, P]) Len() int {
	return len(q.items)
}

// Less reports whether the element with index i must sort before the element with index j.
func (q *queueHeap[T, P]) Less(i, j int) bool {
	a, b := q.items[i], q.items[j]
	return q.compare(a.Priority, b.Priority) < 0
}

// Swap swaps the elements with indexes i and j.
func (q *queueHeap[T, P]) Swap(i, j int) {
	q.items[i], q.items[j] = q.items[j], q.items[i]
}

// Push adds x as element Len()
func (q *queueHeap[T, P]) Push(x any) {
	item := x.(Item[T, P])
	q.items = append(q.items, item)
}

// Pop removes and return element Len() - 1.
func (q *queueHeap[T, P]) Pop() any {
	n := len(q.items)
	x := q.items[n-1]
	q.items = q.items[:n-1]
	return x
}
