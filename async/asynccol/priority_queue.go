// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asynccol

import (
	"github.com/basecomplextech/baselibrary/async/asynccol/internal/priorityqueue"
	"github.com/basecomplextech/baselibrary/compare"
	"golang.org/x/exp/constraints"
)

// PriorityQueue is a priority queue guarded by a mutex.
type PriorityQueue[T any, P any] = priorityqueue.Queue[T, P]

// NewPriorityQueue returns a new priority queue with an ordered priority.
func NewPriorityQueue[T any, P constraints.Ordered]() PriorityQueue[T, P] {
	return priorityqueue.New[T, P]()
}

// NewPriorityQueueCompare returns a new priority queue with a priority compare function.
func NewPriorityQueueCompare[T any, P any](compare compare.Compare[P]) PriorityQueue[T, P] {
	return priorityqueue.NewCompare[T, P](compare)
}
