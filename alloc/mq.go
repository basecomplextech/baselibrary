package alloc

import (
	"github.com/basecomplextech/baselibrary/alloc/internal/mq"
	"github.com/basecomplextech/baselibrary/alloc/internal/msgqueue"
)

// MessageQueue is a queue which transfers binary messages, internaly allocates memory in blocks.
type MessageQueue = msgqueue.MessageQueue

// NewMessageQueue allocates an unbounded buffer queue.
func NewMessageQueue() MessageQueue {
	return global.MessageQueue()
}

// NewMessageQueueCap allocates a buffer queue with a max capacity.
func NewMessageQueueCap(cap int) MessageQueue {
	return global.MessageQueueCap(cap)
}

// MQueue is a single reader multiple writers binary message queue.
// The queue can be unbounded, or can be configured with a soft max capacity.
// Writes mostly do not block readers.
type MQueue = mq.MQueue

// NewMQueue allocates an unbounded buffer queue.
func NewMQueue() MQueue {
	return global.MQueue()
}

// NewMQueueCap allocates a buffer queue with a max capacity.
func NewMQueueCap(cap int) MQueue {
	return global.MQueueCap(cap)
}
