package pqueue

import (
	"container/heap"

	"github.com/complex1tech/baselibrary/compare"
	"github.com/complex1tech/baselibrary/constraints"
	"github.com/complex1tech/baselibrary/slices"
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

// New returns a new priority queue with a priority compare function.
func New[V any, P any](compare compare.Func[P], items ...Item[V, P]) *Queue[V, P] {
	q := &Queue[V, P]{
		queue: &heapq[V, P]{
			compare: compare,
			items:   slices.Clone(items),
		},
	}

	heap.Init(q.queue)
	return q
}

// Ordered returns a new priority queue with a priority natural ordering.
func Ordered[V any, P constraints.Ordered](items ...Item[V, P]) *Queue[V, P] {
	compare := compare.Ordered[P]()
	return New(compare, items...)
}

// Len returns the number of elements in the queue.
func (q *Queue[V, P]) Len() int {
	return q.queue.Len()
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// Pop is equivalent to Remove(h, 0).
func (q *Queue[V, P]) Pop() (value V, ok bool) {
	if q.Len() == 0 {
		return value, false
	}

	item := heap.Pop(q.queue).(Item[V, P])
	return item.Value, true
}

// Push pushes the element x onto the heap.
func (q *Queue[V, P]) Push(value V, priority P) {
	item := Item[V, P]{
		Value:    value,
		Priority: priority,
	}

	heap.Push(q.queue, item)
}

// // Remove removes and returns the element at index i from the heap.
// func (q *Queue[V, P]) Remove(i int) (value V, ok bool) {
// 	if i >= q.Len() {
// 		return value, false
// 	}
//
// 	item := heap.Remove(q.queue, i).(Item[V, P])
// 	return item.Value, true
// }
