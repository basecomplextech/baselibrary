// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package collect

import (
	"github.com/basecomplextech/baselibrary/collect/internal/priorityqueue"
	"github.com/basecomplextech/baselibrary/compare"
	"golang.org/x/exp/constraints"
)

type (
	// PriorityQueue is a priority queue.
	PriorityQueue[T any, P any] = priorityqueue.Queue[T, P]

	// PriorityQueueItem is an element of the queue.
	PriorityQueueItem[T any, P any] = priorityqueue.Item[T, P]
)

// NewPriorityQueue returns a new priority queue with an ordered priority.
func NewPriorityQueue[T any, P constraints.Ordered](
	items ...PriorityQueueItem[T, P]) PriorityQueue[T, P] {

	return priorityqueue.New(items...)
}

// NewPriorityQueueCompare returns a new priority queue with a priority compare function.
func NewPriorityQueueCompare[T any, P any](compare compare.Compare[P],
	items ...PriorityQueueItem[T, P]) PriorityQueue[T, P] {

	return priorityqueue.NewCompare(compare, items...)
}
