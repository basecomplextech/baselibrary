package alloc

import "github.com/basecomplextech/baselibrary/alloc/internal/bufqueue"

// BufferQueue is a queue which transfers binary messages, internaly allocates memory in blocks.
type BufferQueue = bufqueue.BufferQueue

// NewBufferQueue allocates an unbounded buffer queue.
func NewBufferQueue() BufferQueue {
	return global.BufferQueue()
}

// NewBufferQueueCap allocates a buffer queue with a max capacity.
func NewBufferQueueCap(cap int) BufferQueue {
	return global.BufferQueueCap(cap)
}
