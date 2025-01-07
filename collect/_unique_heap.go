// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package priorityqueue

// Len is the number of elements in the collection.
func (q *uniqueHeap[V, P]) Len() int {
	return len(q.items)
}

// Less reports whether the element with index i must sort before the element with index j.
func (q *uniqueHeap[V, P]) Less(i, j int) bool {
	a, b := q.items[i], q.items[j]
	return q.compare(a.priority, b.priority) < 0
}

// Swap swaps the elements with indexes i and j.
func (q *uniqueHeap[V, P]) Swap(i, j int) {
	a, b := q.items[i], q.items[j]

	a.index = j
	b.index = i

	q.items[i] = b
	q.items[j] = a

	q.indexes[a.value] = j
	q.indexes[b.value] = i
}

// Push adds x as element Len()
func (q *uniqueHeap[V, P]) Push(x any) {
	item := x.(uniqueItem[V, P])
	item.index = len(q.items)

	q.items = append(q.items, item)
	q.indexes[item.value] = item.index
}

// Pop removes and return element Len() - 1.
func (q *uniqueHeap[V, P]) Pop() any {
	n := len(q.items)
	x := q.items[n-1]

	q.items = q.items[:n-1]
	delete(q.indexes, x.value)
	return x
}
