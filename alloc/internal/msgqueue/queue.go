package msgqueue

import (
	"sync"

	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
	"github.com/basecomplextech/baselibrary/status"
)

// MessageQueue is byte message queue, which internaly allocates memory in blocks.
type MessageQueue interface {
	// Cap returns the maximum capacity of the queue in bytes, 0 means unlimited.
	Cap() int

	// Len returns the number of unread bytes in the queue.
	Len() int

	// Close

	// Close closes the queue for writing, it is still possible to read the existing messages.
	Close()

	// CloseWithError closes the queue for writing, it is still possible to read the existing messages.
	CloseWithError(st status.Status)

	// Read/write

	// Read reads an message from the queue, the message is valid until the next iteration.
	// The method returns a close status when there are no more items and the queue is closed.
	Read() ([]byte, bool, status.Status)

	// Write writes an message to the queue, returns (false, status.Wait) if the queue is full.
	// The method returns a close status if the queue is closed.
	Write(msg []byte) (bool, status.Status)

	// Wait

	// Wait returns a channel which is notified when the queue is not empty.
	// The method returns a closed channel if the queue is closed.
	Wait() <-chan struct{}

	// WaitNotFull returns a channel which is notified when the queue is not full.
	// The method returns a closed channel if the queue is closed.
	WaitNotFull(size int) <-chan struct{}

	// Internal

	// Free releases the queue and its iternal resources.
	Free()
}

// New returns an unbounded queue.
func New(heap *heap.Heap) MessageQueue {
	return newQueue(heap, 0)
}

// NewCap returns a bounded queue with a maximum capacity in bytes, 0 means unlimited.
// The capacity specifies a soft limit on the maximum number of bytes in the queue.
// The queue is still able to allocate more memory if needed, for example,
// to send bigger messages.
func NewCap(heap *heap.Heap, cap int) MessageQueue {
	return newQueue(heap, cap)
}

// internal

const maxBlockSize = 1 << 17 // 128K

var _ MessageQueue = (*queue)(nil)

type queue struct {
	heap *heap.Heap
	cap  int // total capacity

	mu     sync.Mutex
	st     status.Status
	blocks []*block

	wait     chan struct{}
	waitFlag bool

	waitNotFull     chan struct{}
	waitNotFullFlag bool

	maxCap int
}

func newQueue(heap *heap.Heap, cap int) *queue {
	if cap < 0 {
		panic("negative capacity")
	}

	return &queue{
		heap: heap,
		cap:  cap,
		st:   status.OK,

		wait:        make(chan struct{}, 1),
		waitNotFull: make(chan struct{}, 1),
	}
}

// Cap returns the maximum capacity of the queue in bytes, 0 means unlimited.
func (q *queue) Cap() int {
	return q.cap
}

// Len returns the number of unread bytes in the queue.
func (q *queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	n := 0
	for _, b := range q.blocks {
		n += b.unread()
	}
	return n
}

// Close

// Close closes the queue for writing, it is still possible to read the existing messages.
func (q *queue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.st.OK() {
		return
	}

	q.st = status.End
	q.notifyNotEmpty()
	q.notifyNotFull()
}

// CloseWithError closes the queue for writing, it is still possible to read the existing messages.
func (q *queue) CloseWithError(st status.Status) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.st.OK() {
		return
	}

	q.st = st
	q.notifyNotEmpty()
	q.notifyNotFull()
}

// Read/write

// Read reads an message from the queue, the message is valid until the next iteration.
// The method returns a close status when there are no more items and the queue is closed.
func (q *queue) Read() ([]byte, bool, status.Status) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.waitFlag = false

	for {
		// Get first block
		block, ok := q.first()
		if !ok {
			return nil, false, q.st
		}

		// Move to the next message
		msg, ok := block.next()
		switch {
		case ok:
			// Return the next message
			q.notifyNotFull()
			return msg, true, status.OK

		case len(q.blocks) == 1:
			// Reset if the only block and no more messages
			block.reset()
			q.notifyNotFull()
			return nil, false, status.OK

		default:
			// Release the block if read and no more messages
			q.releaseBlock()
			q.notifyNotFull()
		}
	}
}

// Write writes an message to the queue, returns (false, status.Wait) if the queue is full.
// The method returns a close status if the queue is closed.
func (q *queue) Write(msg []byte) (bool, status.Status) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.st.OK() {
		return false, q.st
	}

	n := len(msg) + 4

	// Try to get the last block
	block, ok := q.last()
	switch {
	case !ok:
		// Alloc a new block if no blocks
		block = q.allocBlock(n)

	case block.rem() < n:
		// Try to alloc a new block, return false when full
		if !q.canAllocBlock(n) {
			return false, status.OK
		}
		block = q.allocBlock(n)
	}

	// Write message
	block.write(msg)
	q.notifyNotEmpty()
	return true, status.OK
}

// Wait

// Wait returns a channel which is notified when the queue is not empty.
// The method returns a closed channel if the queue is closed.
func (q *queue) Wait() <-chan struct{} {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.st.OK() {
		return closedChan
	}

	// Return if not empty
	if !q.empty() {
		return closedChan
	}

	// Wait for more messages
	select {
	case <-q.wait:
	default:
	}

	q.waitFlag = true
	return q.wait
}

// WaitNotFull returns a channel which is notified when the queue is not full.
// The method returns a closed channel if the queue is closed.
func (q *queue) WaitNotFull(size int) <-chan struct{} {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.st.OK() {
		return closedChan
	}

	n := 4 + size

	// Return if can write n
	last, ok := q.last()
	switch {
	case ok && last.rem() >= n:
		return closedChan
	case q.canAllocBlock(n):
		return closedChan
	}

	// Wait for more space
	select {
	case <-q.waitNotFull:
	default:
	}

	q.waitNotFullFlag = true
	return q.waitNotFull
}

// Internal

// Free releases the queue and its iternal resources.
func (q *queue) Free() {

}

// private

func (q *queue) first() (*block, bool) {
	if len(q.blocks) == 0 {
		return nil, false
	}

	return q.blocks[0], true
}

func (q *queue) last() (*block, bool) {
	if len(q.blocks) == 0 {
		return nil, false
	}

	return q.blocks[len(q.blocks)-1], true
}

func (q *queue) empty() bool {
	if len(q.blocks) == 0 {
		return true
	}

	unread := q.blocks[0].unread()
	return unread == 0
}

// allocBlock allocates a new block which can hold at least n bytes.
func (q *queue) allocBlock(n int) *block {
	size := 0

	// Double last block capacity if possible,
	// but no more than 1/4 of the queue capacity
	if len(q.blocks) > 0 {
		last := q.blocks[len(q.blocks)-1]
		size = last.Cap() * 2

		if q.cap > 0 {
			max := q.cap / 4
			if size > max {
				size = max
			}
		}
		if size > maxBlockSize {
			size = maxBlockSize
		}
	}

	// Use the requested size if bigger
	if n > size {
		size = n
	}

	b := q.heap.Alloc(size)
	block := &block{Block: b}
	q.blocks = append(q.blocks, block)

	cap := q.blocksCapacity()
	if cap > q.maxCap {
		q.maxCap = cap
	}
	return block
}

// canAllocBlock returns true if a new block can be allocated which can hold at least n bytes.
func (q *queue) canAllocBlock(n int) bool {
	switch {
	case q.cap == 0:
		return true
	case len(q.blocks) == 0:
		return true
	}

	cap := q.blocksCapacity()
	return cap+n <= q.cap
}

// releaseBlock releases the first block if present, panics on unread messages.
func (q *queue) releaseBlock() {
	if len(q.blocks) == 0 {
		return
	}

	first := q.blocks[0]
	if first.read != first.Len() {
		panic("cannot release block with unread data")
	}

	copy(q.blocks, q.blocks[1:])
	q.blocks[len(q.blocks)-1] = nil
	q.blocks = q.blocks[:len(q.blocks)-1]

	q.heap.Free(first.Block)
}

// blocksCapacity returns the total capacity of all blocks.
func (q *queue) blocksCapacity() int {
	cap := 0
	for _, block := range q.blocks {
		cap += block.Cap()
	}
	return cap
}

// notify

func (q *queue) notifyNotEmpty() {
	if !q.waitFlag {
		return
	}

	select {
	case q.wait <- struct{}{}:
	default:
	}
}

func (q *queue) notifyNotFull() {
	if !q.waitNotFullFlag {
		return
	}

	select {
	case q.waitNotFull <- struct{}{}:
	default:
	}
}

// closed channel

var closedChan = func() chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}()
