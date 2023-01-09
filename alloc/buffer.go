package alloc

import (
	"github.com/complex1tech/baselibrary/buffer"
	"github.com/complex1tech/baselibrary/ref"
)

var (
	_ buffer.Buffer = (*Buffer)(nil)
	_ ref.Freer     = (*Buffer)(nil)
)

type Buffer struct {
	a *allocator

	len    int
	blocks []*block
}

// NewBuffer returns a new empty buffer.
func NewBuffer() *Buffer {
	return newBuffer(global)
}

// NewBuffer returns a new empty buffer with a preallocated memory storage.
func NewBufferSize(size int) *Buffer {
	b := newBuffer(global)
	b.allocBlock(size)
	return b
}

func newBuffer(a *allocator) *Buffer {
	return &Buffer{a: a}
}

// Free frees the buffer and releases its memory to the allocator.
func (b *Buffer) Free() {
	b.len = 0
	b.clearBlocks()
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
	return b.blocks[0].buf
}

// Grow grows the buffer and returns an n-byte slice.
func (b *Buffer) Grow(n int) []byte {
	// maybe allocate block
	last := b.last()
	if last == nil || last.free() < n {
		last = b.allocBlock(n)
	}

	// grow buffer
	start := len(last.buf)
	end := start + n
	last.buf = last.buf[:end]
	b.len += n

	// slice buffer
	p := last.buf[start:end:end] // start:end:max, cap=max-start
	return p
}

// Write appends bytes from p to the buffer.
func (b *Buffer) Write(p []byte) (n int, err error) {
	q := b.Grow(len(p))
	n = copy(q, p)
	return
}

// Reset resets the buffer to be empty, releases its memory storage to the allocator.
func (b *Buffer) Reset() {
	b.len = 0
	b.clearBlocks()
}

// private

// last returns the last block or nil.
func (b *Buffer) last() *block {
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

	// alloc block
	one, _ := b.a.allocBlock(b.len)

	// merge data
	for _, block := range b.blocks {

		// grow buffer
		start := len(one.buf)
		end := start + len(block.buf)
		one.buf = one.buf[:end]

		// copy data
		p := one.buf[start:end]
		copy(p, block.buf)
	}

	// replace blocks
	b.clearBlocks()
	b.blocks = append(b.blocks, one)
}

// allocBlock allocates the next block.
func (b *Buffer) allocBlock(n int) *block {
	// double last block size
	// limit it to maxBlockSize
	size := 0
	last := b.last()

	if last != nil {
		size = last.cap() * 2
	}
	if size > maxBlockSize {
		size = maxBlockSize
	}
	if n > size {
		size = n
	}

	block, _ := b.a.allocBlock(size)
	b.blocks = append(b.blocks, block)
	return block
}

// clearBlocks clears and frees blocks.
func (b *Buffer) clearBlocks() {
	// free blocks
	b.a.freeBlocks(b.blocks...)

	// clear blocks for gc
	for i := range b.blocks {
		b.blocks[i] = nil
	}
	b.blocks = b.blocks[:0]
}
