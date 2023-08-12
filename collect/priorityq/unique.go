package priorityq

import (
	"container/heap"
)

// UniqueQueue is a priority queue that contains unique elements.
type UniqueQueue[V comparable, P any] struct {
	queue *uniqueHeap[V, P]
}

// NewUnique returns a unique priority queue with a priority compare function.
func NewUnique[V comparable, P any](compare CompareFunc[P], items ...Item[V, P]) *UniqueQueue[V, P] {
	q := &UniqueQueue[V, P]{
		queue: &uniqueHeap[V, P]{
			compare: compare,
			items:   make([]uniqueItem[V, P], 0, len(items)),
			indexes: make(map[V]int),
		},
	}

	for i, item := range items {
		item := uniqueItem[V, P]{
			value:    item.Value,
			priority: item.Priority,
			index:    i,
		}

		q.queue.items = append(q.queue.items, item)
		q.queue.indexes[item.value] = i
	}

	heap.Init(q.queue)
	return q
}

// Len returns the number of elements in the queue.
func (q *UniqueQueue[V, P]) Len() int {
	return q.queue.Len()
}

// Clear removes all elements from the queue.
func (q *UniqueQueue[V, P]) Clear() {
	q.queue.items = nil
	q.queue.indexes = make(map[V]int)
}

// Contains returns true if the queue contains an element.
func (q *UniqueQueue[V, P]) Contains(value V) bool {
	i := q.queue.indexOf(value)
	return i >= 0
}

// Get returns the priority of an element.
func (q *UniqueQueue[V, P]) Get(value V) (p P, ok bool) {
	i := q.queue.indexOf(value)
	if i == -1 {
		return p, false
	}

	return q.queue.items[i].priority, true
}

// Pop removes and returns the minimum element (according to Less) from the queue.
func (q *UniqueQueue[V, P]) Pop() (value V, priority P, ok bool) {
	if q.Len() == 0 {
		return value, priority, false
	}

	item := heap.Pop(q.queue).(uniqueItem[V, P])
	return item.value, item.priority, true
}

// Push pushes an element onto the queue.
func (q *UniqueQueue[V, P]) Push(value V, priority P) {
	i := q.queue.indexOf(value)
	if i != -1 {
		q.queue.update(i, priority)
		return
	}

	item := uniqueItem[V, P]{
		value:    value,
		priority: priority,
	}
	heap.Push(q.queue, item)
}

// Remove removes an element from the queue, and returns its priority.
func (q *UniqueQueue[V, P]) Remove(value V) (p P, ok bool) {
	i := q.queue.indexOf(value)
	if i == -1 {
		return p, false
	}

	item := heap.Remove(q.queue, i).(uniqueItem[V, P])
	return item.priority, true
}

// internal

var _ heap.Interface = (*uniqueHeap[int, any])(nil)

type uniqueHeap[V comparable, P any] struct {
	compare CompareFunc[P]
	items   []uniqueItem[V, P]
	indexes map[V]int
}

type uniqueItem[V any, P any] struct {
	value    V
	priority P
	index    int
}

func (q *uniqueHeap[V, P]) indexOf(value V) int {
	i, ok := q.indexes[value]
	if !ok {
		return -1
	}
	return i
}

func (q *uniqueHeap[V, P]) update(index int, priority P) {
	q.items[index].priority = priority
	heap.Fix(q, index)
}

// heap.Interface

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
