package mq

import (
	"math"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
	"github.com/basecomplextech/baselibrary/collect/slices"
	"github.com/basecomplextech/baselibrary/status"
)

var _ MessageQueue = (*queue)(nil)

const (
	// maxBlockSize tries to keep blocks from growing too large.
	// larger blocks can still be allocated to fit large messages.
	maxBlockSize = 1 << 17 // 128K
)

type queue struct {
	cap  int // maximum queue capacity, it is a soft limit, 0 means unlimited
	heap *heap.Heap

	// single reader lock
	rmu sync.Mutex

	// channels for reader/writer to wait on
	readChan  chan struct{}
	writeChan chan struct{}

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

		readChan:  make(chan struct{}, 1),
		writeChan: make(chan struct{}, 1),
	}
}

// clear releases all unread messages.
func (q *queue) clear() {
	q.rmu.Lock()
	defer q.rmu.Unlock()

	q.mu.Lock()
	defer q.mu.Unlock()
	defer q.notifyWrite()

	q.freeBlocks()
}

// close closes the queue, it is still possible to read the existing message from it.
func (q *queue) close() {
	q.rmu.Lock()
	defer q.rmu.Unlock()

	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return
	}

	q.closed = true
	close(q.readChan)
	close(q.writeChan)
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
	// Fast path: atomically load head and return if it is not empty.
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
		// Check head again.
		head := q.head
		if head == nil {
			break
		}

		// Return if not empty.
		ri := head.readIndex
		wi := head.writeIndex
		if ri < wi {
			return head, true, status.OK
		}

		// Head is empty, it can be reset or released.
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

	// Head is nil or empty, and there are no more blocks.
	if q.closed {
		return nil, false, status.End
	}
	return nil, false, status.OK
}

// readWait returns a channel to wait for more messages or an end.
func (q *queue) readWait() <-chan struct{} {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return closedChan
	}

	// Check head not empty.
	head := q.head
	if head != nil {
		ri := head.readIndex
		wi := head.loadWriteIndex()
		if ri < wi {
			return closedChan
		}
	}

	select {
	case <-q.readChan:
	default:
	}

	return q.readChan
}

// notifyRead notifies a waiting reader.
func (q *queue) notifyRead() {
	// Write mutex must be locked.
	if q.closed {
		return
	}

	select {
	case q.readChan <- struct{}{}:
	default:
	}
}

// write

// write writes a message or returns an end status if the queue is closed.
func (q *queue) write(msg []byte) (bool, status.Status) {
	if len(msg) > math.MaxInt32 {
		panic("message too large")
	}

	q.mu.Lock()
	defer q.mu.Unlock()
	defer q.notifyRead()

	if q.closed {
		return false, status.End
	}

	size := len(msg)

	// Get a block to write to.
	block, ok := q._writeBlock(size)
	if !ok {
		return false, status.OK
	}

	// Write message to the block.
	wi := block.copy(msg)
	block.storeWriteIndex(wi)
	return true, status.OK
}

// _writeBlock returns or allocates a block to write to.
func (q *queue) _writeBlock(size int) (*block, bool) {
	n := 4 + size

	// Return tail if it has enough free space.
	tail := q.tail()
	if tail != nil && tail.free() >= n {
		return tail, true
	}

	// Check queue is not full. Messages can be larger than the queue max capacity.
	// In this case we write it to a new block, but only if there are no more blocks yet.
	if q.cap > 0 {
		large := n > q.cap
		if large {
			if len(q.more) > 0 {
				return nil, false
			}
		} else {
			total := q.occupied()
			if total >= q.cap {
				return nil, false
			}
		}
	}

	// Allocate a new block.
	block := q.alloc(n)
	return block, true
}

// writeWait returns a channel to wait for more space.
func (q *queue) writeWait(size int) <-chan struct{} {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return closedChan
	}

	_, ok := q._writeBlock(size)
	if ok {
		return closedChan
	}

	select {
	case <-q.writeChan:
	default:
	}

	return q.writeChan
}

// notifyWrite notifies a waiting writer.
func (q *queue) notifyWrite() {
	// Read mutex must be locked.
	if q.closed {
		return
	}

	select {
	case q.writeChan <- struct{}{}:
	default:
	}
}

// more

// reset resets the queue, releases all unread messages, the queue can be used again.
func (q *queue) reset() {
	q.rmu.Lock()
	defer q.rmu.Unlock()

	q.mu.Lock()
	defer q.mu.Unlock()

	q.freeBlocks()
	q.notifyRead()
	q.notifyWrite()

	if q.closed {
		q.closed = false
		q.readChan = make(chan struct{}, 1)
		q.writeChan = make(chan struct{}, 1)
	}
}

// free releases the queue and its iternal resources.
func (q *queue) free() {
	q.rmu.Lock()
	defer q.rmu.Unlock()

	q.mu.Lock()
	defer q.mu.Unlock()

	q.freeBlocks()
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

func (q *queue) freeBlocks() {
	if q.head != nil {
		q.heap.Free(q.head.b)
		q.head = nil
	}

	for _, b := range q.more {
		q.heap.Free(b.b)
	}

	slices.Clear(q.more)
	q.more = q.more[:0]
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
