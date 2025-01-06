// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asynccol

import "github.com/basecomplextech/baselibrary/async/asynccol/internal/queue"

type (
	// Queue is an unbounded FIFO queue.
	Queue[T any] = queue.Queue[T]
)

// NewQueue returns an empty queue.
func NewQueue[T any](items ...T) Queue[T] {
	return queue.New[T](items...)
}
