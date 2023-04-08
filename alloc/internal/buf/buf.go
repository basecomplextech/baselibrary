package buf

import (
	"unicode/utf8"

	"github.com/complex1tech/baselibrary/alloc/internal/heap"
	"github.com/complex1tech/baselibrary/buffer"
	"github.com/complex1tech/baselibrary/collect/slices"
	"github.com/complex1tech/baselibrary/ref"
)

var (
	_ buffer.Buffer = (*Buffer)(nil)
	_ buffer.Writer = (*Buffer)(nil)
	_ ref.Freer     = (*Buffer)(nil)
)

// Buffer is a byte buffer, which internally allocates memory in blocks.
type Buffer struct {
	heap    *heap.Heap
	initCap int // initial capacity

	blocks []*heap.Block
	len    int // total length in bytes
}

// New returns a new buffer.
func New(heap *heap.Heap) *Buffer {
	return &Buffer{heap: heap}
}

// NewSize returns a new buffer with a preallocated memory storage.
func NewSize(heap *heap.Heap, size int) *Buffer {
	b := New(heap)
	if size > 0 {
		b.initCap = b.allocBlock(size).Cap()
	}
	return b
}

// Free frees the buffer and releases its memory to the heap.
func (b *Buffer) Free() {
	b.len = 0
	b.freeBlocks()
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
	if last == nil || last.Free() < n {
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
	b.len = 0
	if len(b.blocks) == 0 {
		return
	}

	// Maybe just reset the first block
	n := 0
	if f := b.blocks[0]; f.Cap() == b.initCap {
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

// allocBlock allocates the next block.
func (b *Buffer) allocBlock(n int) *heap.Block {
	// Use initial size or double last block capacity
	size := 0
	if len(b.blocks) == 0 {
		size = b.initCap
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
