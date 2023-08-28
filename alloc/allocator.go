package alloc

import (
	"github.com/basecomplextech/baselibrary/alloc/internal/arena"
	"github.com/basecomplextech/baselibrary/alloc/internal/buf"
	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
	"github.com/basecomplextech/baselibrary/alloc/internal/mq"
	"github.com/basecomplextech/baselibrary/alloc/internal/msgqueue"
)

var global = newAllocator(heap.Global())

type allocator struct {
	heap *heap.Heap
}

func newAllocator(heap *heap.Heap) *allocator {
	return &allocator{heap: heap}
}

// Arena allocates a new arena.
func (a *allocator) Arena() Arena {
	return arena.New(a.heap)
}

// ArenaSize allocates a new arena with a preallocated capacity.
func (a *allocator) ArenaSize(size int) Arena {
	return arena.NewSize(a.heap, size)
}

// Buffer allocates a new buffer.
func (a *allocator) Buffer() *Buffer {
	return buf.New(a.heap)
}

// BufferSize allocates a new buffer with a preallocated capacity.
func (a *allocator) BufferSize(size int) *Buffer {
	return buf.NewSize(a.heap, size)
}

// MessageQueue allocates an unbounded buffer queue.
func (a *allocator) MessageQueue() MessageQueue {
	return msgqueue.New(a.heap)
}

// MessageQueueCap allocates a new buffer queue with a max capacity.
func (a *allocator) MessageQueueCap(cap int) MessageQueue {
	return msgqueue.NewCap(a.heap, cap)
}

// NewMQueue allocates an unbounded message queue.
func (a *allocator) MQueue() MQueue {
	return mq.New(a.heap)
}

// NewMQueueCap allocates a message queue with a max capacity.
func (a *allocator) MQueueCap(cap int) MQueue {
	return mq.NewCap(a.heap, cap)
}
