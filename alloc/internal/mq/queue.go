package mq

import (
	"math"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
	"github.com/basecomplextech/baselibrary/status"
)

var _ MessageQueue = (*queue)(nil)

// maxBlockSize tries to keep blocks from growing too large.
// larger blocks can still be allocated to fit large messages.
const maxBlockSize = 1 << 17 // 128K

type queue struct {
	cap  int // maximum queue capacity, it is a soft limit, 0 means unlimited
	heap *heap.Heap

	// single reader/writer locks
	rmu sync.Mutex
	wmu sync.Mutex

	// channels to reader/writer to wait on
	readNotify  chan struct{}
	writeNotify chan struct{}

	// state
	mu     sync.Mutex
	closed bool
	head   *block // can be accessed atomically by reader
	more   []*block
}

func newQueue(heap *heap.Heap, cap int) *queue {
	return &queue{
		cap:  cap,
		heap: heap,

		readNotify:  make(chan struct{}, 1),
		writeNotify: make(chan struct{}, 1),
	}
}

// close closes the queue, it is still possible to read the existing message from it.
func (q *queue) close() {
	q.rmu.Lock()
	defer q.rmu.Unlock()

	q.wmu.Lock()
	defer q.wmu.Unlock()

	q.mu.Lock()
	defer q.mu.Unlock()

	q.closed = true
	close(q.readNotify)
	close(q.writeNotify)
}

// read reads the next message, the message is valid until the next call to read.
func (q *queue) read() ([]byte, bool, status.Status) {
	q.rmu.Lock()
	defer q.rmu.Unlock()

	block, ok, st := q.readBlock()
	switch {
	case !st.OK():
		return nil, false, st
	case !ok:
		return nil, false, status.OK
	}

	msg := block.read()
	return msg, true, status.OK
}

// readBlock returns a block to read from, or nil if the queue is empty.
func (q *queue) readBlock() (*block, bool, status.Status) {
	// Fast path: atomically load head and return if it is not empty
	{
		head := q.loadHead()
		if head == nil {
			return nil, false, status.OK
		}

		ri := head.readIndex
		wi := head.loadWriteIndex()
		if ri < wi {
			return head, true, status.OK
		}
	}

	// Slow path
	q.mu.Lock()
	defer q.mu.Unlock()
	defer q.notifyWrite()

	for {
		// Check head again
		head := q.head
		if head == nil {
			break
		}

		// Return if not empty
		ri := head.readIndex
		wi := head.writeIndex
		if ri < wi {
			return head, true, status.OK
		}

		// Head is empty, it may also be acquired by the writer.
		// We have to wait for the writer to release it.
		if head.acquired {
			return nil, false, status.OK
		}

		// Head is empty and not acquired by the writer,
		// so it can be reset or released.
		if len(q.more) == 0 {
			head.reset()
			break
		}

		// Release head to the heap,
		// and move the next block to head.
		q.heap.Free(head.b)
		q.head = q.more[0]

		copy(q.more, q.more[1:])
		q.more[len(q.more)-1] = nil
		q.more = q.more[:len(q.more)-1]
	}

	// Head is nil or (empty, not acquired and there are no more blocks).
	if q.closed {
		return nil, false, status.End
	}
	return nil, false, status.OK
}

// readWait returns a channel to wait for more messages or a close status.
func (q *queue) readWait() <-chan struct{} {
	q.rmu.Lock()
	defer q.rmu.Unlock()

	_, ok, st := q.readBlock()
	switch {
	case !st.OK():
		return closedChan
	case ok:
		return closedChan
	}

	return q.readNotify
}

// notifyRead notifies a waiting reader.
func (q *queue) notifyRead() {
	// Write mutex must be locked.
	if q.closed {
		return
	}

	select {
	case q.readNotify <- struct{}{}:
	default:
	}
}

// write

// write writes a message or returns an end status if the queue is closed.
func (q *queue) write(msg []byte) (bool, status.Status) {
	if len(msg) > math.MaxInt32 {
		panic("message too large")
	}

	q.wmu.Lock()
	defer q.wmu.Unlock()
	defer q.notifyRead()

	// Acquire a block to write to.
	block, ok, st := q.writeBlock(len(msg), true /* acquire block */)
	switch {
	case !st.OK():
		return false, st
	case !ok:
		return false, status.OK
	}

	// Copy message to the block outside of the lock.
	wi := block.copy(msg)

	// Commit the write.
	q.mu.Lock()
	defer q.mu.Unlock()

	block.acquired = false
	if q.closed {
		return false, status.End
	}

	block.storeWriteIndex(wi)
	return true, status.OK
}

// writeBlock acquires a block to write to.
func (q *queue) writeBlock(size int, acquire bool) (*block, bool, status.Status) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return nil, false, status.End
	}

	n := 4 + size

	// Acquire tail if it has enough free space.
	tail := q.tail()
	if tail != nil && tail.free() >= n {
		if acquire {
			tail.acquired = true
		}
		return tail, true, status.OK
	}

	// Check queue is not full. Also messages can be larger than the queue max capacity.
	// In this case we write it to a new block, but only if there are no more blocks yet.
	if q.cap > 0 {
		large := n > q.cap
		if large {
			if len(q.more) > 0 {
				return nil, false, status.OK
			}
		} else {
			total := q.occupied()
			if total >= q.cap {
				return nil, false, status.OK
			}
		}
	}

	// Allocate a new block.
	block := q.alloc(n)
	if acquire {
		block.acquired = true
	}
	return block, true, status.OK
}

// writeWait returns a channel to wait for more space.
func (q *queue) writeWait(size int) <-chan struct{} {
	q.wmu.Lock()
	defer q.wmu.Unlock()

	if q.closed {
		return closedChan
	}

	_, ok, st := q.writeBlock(size, false /* do not acquire block */)
	switch {
	case !st.OK():
		return closedChan
	case ok:
		return closedChan
	}

	return q.writeNotify
}

// notifyWrite notifies a waiting writer.
func (q *queue) notifyWrite() {
	// Read mutex must be locked.
	if q.closed {
		return
	}

	select {
	case q.writeNotify <- struct{}{}:
	default:
	}
}

// private

// occupied returns the total written bytes.
func (q *queue) occupied() int {
	if q.head == nil {
		return 0
	}

	var n int
	if len(q.more) == 0 {
		n = int(q.head.writeIndex)
	} else {
		n = q.head.cap()
	}

	for i, b := range q.more {
		last := i == len(q.more)-1
		if last {
			n += int(b.writeIndex)
		} else {
			n += b.cap()
		}
	}
	return n
}

// alloc allocates a new block.
func (q *queue) alloc(n int) *block {
	size := 0

	// Double tail block capacity if possible,
	// but no more than 1/4 of the queue capacity.
	tail := q.tail()
	if tail != nil {
		size = tail.cap() * 2

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

	// Use the requested size if larger.
	if n > size {
		size = n
	}

	// Allocate new block.
	b := q.heap.Alloc(size)
	block := newBlock(b)

	if q.head == nil {
		q.storeHead(block)
	} else {
		q.more = append(q.more, block)
	}
	return block
}

// head

func (q *queue) loadHead() *block {
	return (*block)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head))))
}

func (q *queue) storeHead(b *block) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)), unsafe.Pointer(b))
}

// tail

func (q *queue) tail() *block {
	// Get tail block
	if len(q.more) > 0 {
		return q.more[len(q.more)-1]
	}
	return q.head
}

// util

var closedChan = func() chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}()
