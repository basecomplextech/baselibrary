// Copyright 2024 Ivan Korobkov. All rights reserved.

package alloc

import (
	"github.com/basecomplextech/baselibrary/alloc/internal/byteq"
	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
)

// ByteQueue is a single reader multiple writers binary message queue.
//
// The queue can be unbounded, or can be configured with a soft max capacity.
// Writes mostly do not block readers.
type ByteQueue = byteq.Queue

// NewByteQueue allocates an unbounded buffer queue.
func NewByteQueue() ByteQueue {
	return byteq.New(heap.Global)
}

// NewByteQueueCap allocates a buffer queue with a max capacity.
func NewByteQueueCap(cap int) ByteQueue {
	return byteq.NewCap(heap.Global, cap)
}
