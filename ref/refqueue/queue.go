// Copyright 2024 Ivan Korobkov. All rights reserved.

package refqueue

import (
	"sync"

	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/ref"
)

// Queue is a wrapper around async.Queue with reference counting.
type Queue[T any] interface {
	async.Queue[ref.R[T]]
}

// New returns a new queue.
func New[T any]() Queue[T] {
	return newQueue[T]()
}

// internal

var _ Queue[any] = (*queue[any])(nil)

type queue[T any] struct {
	mu sync.Mutex
	q  async.Queue[ref.R[T]]
}

func newQueue[T any]() *queue[T] {
	return &queue[T]{
		q: async.NewQueue[ref.R[T]](),
	}
}

// Clear clears the queue.
func (q *queue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.clear()
}

// Len returns the number of elements in the queue.
func (q *queue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.q.Len()
}

// Push adds an element to the queue, panics if the queue is closed.
func (q *queue[T]) Push(v ref.R[T]) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.q.Push(v)
	v.Retain()
}

// Pop removes an element from the queue, returns false if the queue is empty.
func (q *queue[T]) Pop() (v ref.R[T], ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.q.Pop()
}

// Wait returns a channel which is notified on new elements.
func (q *queue[T]) Wait() <-chan struct{} {
	return q.q.Wait()
}

// Free closes and frees the queue.
func (q *queue[T]) Free() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.clear()
	q.q.Free()
}

// internal

func (q *queue[T]) clear() {
	for {
		v, ok := q.q.Pop()
		if !ok {
			break
		}
		v.Release()
	}
}
