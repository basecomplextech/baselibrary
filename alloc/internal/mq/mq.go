package mq

import "github.com/basecomplextech/baselibrary/status"

// MessageQueue is a single reader multiple writers binary message queue.
// Writes mostly do not block readers.
type MessageQueue interface {
	// Closed returns true if the queue is closed.
	Closed() bool

	// Methods

	// Clear releases all unread messages.
	Clear()

	// Close closes the queue for writing, it is still possible to read the existing messages.
	Close()

	// Read reads an message from the queue, the message is valid until the next iteration.
	// The method returns a close status when there are no more items and the queue is closed.
	Read() ([]byte, bool, status.Status)

	// ReadWait returns a channel which is notified when more messages are available.
	// The method returns a closed channel if the queue is closed.
	ReadWait() <-chan struct{}

	// Write writes an message to the queue, returns false if the queue is full.
	// The method returns a close status if the queue is closed.
	Write(msg []byte) (bool, status.Status)

	// WriteWait returns a channel which is notified when a message can be written.
	// The method returns a closed channel if the queue is closed.
	WriteWait(size int) <-chan struct{}

	// Reset resets the queue, releases all unread messages, the queue can be used again.
	Reset()

	// Internal

	// Free releases the queue and its iternal resources.
	Free()
}

// Closed returns true if the queue is closed.
func (q *queue) Closed() bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.closed
}

// Methods

// Clear releases all unread messages.
func (q *queue) Clear() {}

// Close closes the queue for writing, it is still possible to read the existing messages.
func (q *queue) Close() {
	q.close()
}

// Read reads an message from the queue, the message is valid until the next iteration.
// The method returns a close status when there are no more items and the queue is closed.
func (q *queue) Read() ([]byte, bool, status.Status) {
	return q.read()
}

// ReadWait returns a channel which is notified when more messages are available.
// The method returns a closed channel if the queue is closed.
func (q *queue) ReadWait() <-chan struct{} {
	return q.readWait()
}

// Write writes an message to the queue, returns false if the queue is full.
// The method returns a close status if the queue is closed.
func (q *queue) Write(msg []byte) (bool, status.Status) {
	return q.write(msg)
}

// WriteWait returns a channel which is notified when a message can be written.
// The method returns a closed channel if the queue is closed.
func (q *queue) WriteWait(size int) <-chan struct{} {
	return q.writeWait(size)
}

// Reset resets the queue, releases all unread messages, the queue can be used again.
func (q *queue) Reset() {}

// Free releases the queue and its iternal resources.
func (q *queue) Free() {}
