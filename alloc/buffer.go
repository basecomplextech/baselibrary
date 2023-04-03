package alloc

import (
	"unicode/utf8"

	"github.com/complex1tech/baselibrary/buffer"
	"github.com/complex1tech/baselibrary/ref"
)

var (
	_ buffer.Buffer = (*Buffer)(nil)
	_ buffer.Writer = (*Buffer)(nil)
	_ ref.Freer     = (*Buffer)(nil)
)

type Buffer struct {
	a *allocator

	len    int
	blocks []*block
}

// NewBuffer returns a new empty buffer with the global allocator.
func NewBuffer() *Buffer {
	return newBuffer(global)
}

// NewBuffer returns a new empty buffer with a preallocated memory storage in the global allocator.
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
	// Maybe allocate block
	last := b.last()
	if last == nil || last.free() < n {
		last = b.allocBlock(n)
	}

	// Grow buffer
	start := len(last.buf)
	end := start + n
	last.buf = last.buf[:end]
	b.len += n

	// Slice buffer
	p := last.buf[start:end:end] // start:end:max, cap=max-start
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

	// Alloc block
	one, _ := b.a.allocBlock(b.len)

	// Merge data
	for _, block := range b.blocks {

		// Grow buffer
		start := len(one.buf)
		end := start + len(block.buf)
		one.buf = one.buf[:end]

		// Copy data
		p := one.buf[start:end]
		copy(p, block.buf)
	}

	// Replace blocks
	b.clearBlocks()
	b.blocks = append(b.blocks, one)
}

// allocBlock allocates the next block.
func (b *Buffer) allocBlock(n int) *block {
	// Double last block size
	// Limit it to maxBlockSize
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
	// Free blocks
	b.a.freeBlocks(b.blocks...)

	// Clear blocks for gc
	for i := range b.blocks {
		b.blocks[i] = nil
	}
	b.blocks = b.blocks[:0]
}
