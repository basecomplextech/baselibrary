package alloc

import (
	"github.com/basecomplextech/baselibrary/alloc/internal/mq"
)

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
