package pqueue

import (
	"container/heap"

	"github.com/complex1tech/baselibrary/compare"
)

var _ heap.Interface = (*heapq[any, any])(nil)

type heapq[V any, P any] struct {
	compare compare.Func[P]
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
