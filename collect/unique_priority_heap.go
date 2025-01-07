// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package collect

import (
	"container/heap"

	"github.com/basecomplextech/baselibrary/compare"
)

var _ heap.Interface = (*uniqueHeap[int, any])(nil)

type uniqueHeap[T comparable, P any] struct {
	compare compare.Compare[P]
	items   []uniqueItem[T, P]
	indexes map[T]int
}

type uniqueItem[T any, P any] struct {
	value    T
	priority P
	index    int
}

// Len is the number of elements in the collection.
func (q *uniqueHeap[T, P]) Len() int {
	return len(q.items)
}

// Less reports whether the element with index i must sort before the element with index j.
func (q *uniqueHeap[T, P]) Less(i, j int) bool {
	a, b := q.items[i], q.items[j]
	return q.compare(a.priority, b.priority) < 0
}

// Swap swaps the elements with indexes i and j.
func (q *uniqueHeap[T, P]) Swap(i, j int) {
	a, b := q.items[i], q.items[j]

	a.index = j
	b.index = i

	q.items[i] = b
	q.items[j] = a

	q.indexes[a.value] = j
	q.indexes[b.value] = i
}

// Push appends an element.
func (q *uniqueHeap[T, P]) Push(x any) {
	item := x.(uniqueItem[T, P])
	item.index = len(q.items)

	q.items = append(q.items, item)
	q.indexes[item.value] = item.index
}

// Pop removes and returns the last element.
func (q *uniqueHeap[T, P]) Pop() any {
	n := len(q.items)
	x := q.items[n-1]

	q.items = q.items[:n-1]
	delete(q.indexes, x.value)
	return x
}

// internal

func (q *uniqueHeap[T, P]) indexOf(value T) (int, bool) {
	i, ok := q.indexes[value]
	if !ok {
		return -1, false
	}
	return i, true
}

func (q *uniqueHeap[T, P]) update(index int, priority P) {
	q.items[index].priority = priority
	heap.Fix(q, index)
}
