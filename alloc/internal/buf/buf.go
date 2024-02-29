package buf

import (
	"sync"
	"unicode/utf8"

	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
	"github.com/basecomplextech/baselibrary/buffer"
	"github.com/basecomplextech/baselibrary/collect/slices"
	"github.com/basecomplextech/baselibrary/ref"
)

var (
	_ buffer.Buffer = (*Buffer)(nil)
	_ buffer.Writer = (*Buffer)(nil)
	_ ref.Freer     = (*Buffer)(nil)
)

// Buffer is a byte buffer, which internally allocates memory in blocks.
type Buffer struct {
	*state
}

type state struct {
	heap   *heap.Heap
	pooled bool

	init int // initial capacity
	len  int // total length in bytes

	blocks []*heap.Block
}

// New returns a new buffer.
func New() *Buffer {
	return newBufferSize(heap.Global, heap.MinBlockSize)
}

// NewSize returns a new buffer with a preallocated memory storage.
func NewSize(size int) *Buffer {
	return newBufferSize(heap.Global, size)
}

// Acquire returns a new buffer from the pool.
//
// The buffer must not be used or even referenced after Free.
// Use these method only when buffers do not escape an isolated scope.
//
// Typical usage:
//
//	buf := alloc.AcquireBuffer()
//	defer buf.Free() // free immediately
func Acquire() *Buffer {
	return acquireBuffer()
}

func newBuffer(h *heap.Heap) *Buffer {
	return newBufferSize(h, heap.MinBlockSize)
}

func newBufferSize(heap *heap.Heap, size int) *Buffer {
	b := &Buffer{acquireState()}
	b.heap = heap
	if size > 0 {
		b.init = b.allocBlock(size).Cap()
	}
	return b
}

// Free frees the buffer and releases its memory to the heap.
func (b *Buffer) Free() {
	if b.pooled {
		releaseBuffer(b)
		return
	}

	b.len = 0
	b.freeBlocks()

	s := b.state
	b.state = nil
	releaseState(s)
}

// Len returns the number of bytes in the buffer; b.Len() == len(b.Bytes()).
func (b *Buffer) Len() int {
	return b.len
}

// Bytes returns a byte slice with the buffer bytes.
// It is valid for use only until the next buffer mutation.
func (b *Buffer) Bytes() []byte {
	if len(b.blocks) == 0 {
		return nil
	}

	b.merge()
	return b.blocks[0].Bytes()
}

// Grow grows the buffer and returns an n-byte slice.
func (b *Buffer) Grow(n int) []byte {
	last := b.last()
	if last == nil || last.Rem() < n {
		last = b.allocBlock(n)
	}

	p := last.Grow(n)
	b.len += n
	return p
}

// Write appends bytes from p to the buffer.
func (b *Buffer) Write(p []byte) (n int, err error) {
	q := b.Grow(len(p))
	n = copy(q, p)
	return
}

// WriteByte writes a byte to the buffer.
func (b *Buffer) WriteByte(c byte) error {
	q := b.Grow(1)
	q[0] = c
	return nil
}

// WriteRune writes a rune to the buffer.
func (b *Buffer) WriteRune(r rune) (n int, err error) {
	p := [utf8.UTFMax]byte{}
	n = utf8.EncodeRune(p[:], r)

	q := b.Grow(n)
	copy(q, p[:n])
	return
}

// WriteString writes a string to the buffer.
func (b *Buffer) WriteString(s string) (n int, err error) {
	q := b.Grow(len(s))
	n = copy(q, s)
	return
}

// Reset resets the buffer to be empty.
func (b *Buffer) Reset() {
	b.reset()
}

// private

// last returns the last block or nil.
func (b *Buffer) last() *heap.Block {
	if len(b.blocks) == 0 {
		return nil
	}
	return b.blocks[len(b.blocks)-1]
}

// merge merges multiple blocks into a single one.
func (b *Buffer) merge() {
	if len(b.blocks) <= 1 {
		return
	}

	merged := b.heap.Alloc(b.len)
	for _, block := range b.blocks {
		b := block.Bytes()
		p := merged.Grow(len(b))
		copy(p, b)
	}

	b.freeBlocks()
	b.blocks = append(b.blocks, merged)
}

// reset resets the buffer, frees the blocks except the first one.
func (b *Buffer) reset() {
	b.len = 0
	if len(b.blocks) == 0 {
		return
	}

	// Maybe just reset the first block
	n := 0
	if f := b.blocks[0]; f.Cap() == b.init {
		n = 1
		f.Reset()

		if len(b.blocks) == 1 {
			return
		}
	}

	// Free other blocks
	b.heap.FreeMany(b.blocks[n:]...)
	slices.Zero(b.blocks[n:]) // for gc
	b.blocks = b.blocks[:n]
}

// blocks

// allocBlock allocates the next block.
func (b *Buffer) allocBlock(n int) *heap.Block {
	// Use initial size or double last block capacity
	size := 0
	if len(b.blocks) == 0 {
		size = b.init
	} else {
		last := b.blocks[len(b.blocks)-1]
		size = last.Cap() * 2
	}
	if n > size {
		size = n
	}

	block := b.heap.Alloc(size)
	b.blocks = append(b.blocks, block)
	return block
}

// freeBlocks clears and frees the blocks.
func (b *Buffer) freeBlocks() {
	b.heap.FreeMany(b.blocks...)

	slices.Zero(b.blocks) // for gc
	b.blocks = b.blocks[:0]
}

// pool

var pool = &sync.Pool{
	New: func() any {
		b := newBuffer(heap.Global)
		b.pooled = true
		return b
	},
}

func acquireBuffer() *Buffer {
	return pool.Get().(*Buffer)
}

func releaseBuffer(b *Buffer) {
	b.reset()
	pool.Put(b)
}

// state pool

var statePool = &sync.Pool{
	New: func() any {
		return &state{}
	},
}

func acquireState() *state {
	return statePool.Get().(*state)
}

func releaseState(s *state) {
	s.reset()
	statePool.Put(s)
}

func (s *state) reset() {
	blocks := slices.Clear(s.blocks)

	*s = state{}
	s.blocks = blocks
}
