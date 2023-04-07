package heap

import "unsafe"

type Block struct {
	buf []byte
}

func newBlock(size int) *Block {
	return &Block{
		buf: make([]byte, 0, size),
	}
}

// Cap returns the block capacity in bytes.
func (b *Block) Cap() int {
	return cap(b.buf)
}

// Len returns the number of used bytes in the block.
func (b *Block) Len() int {
	return len(b.buf)
}

// Free returns the free space in bytes.
func (b *Block) Free() int {
	return cap(b.buf) - len(b.buf)
}

// Bytes returns the underlying block bytes.
func (b *Block) Bytes() []byte {
	return b.buf
}

// Reset resets the block.
func (b *Block) Reset() {
	b.buf = b.buf[:0]
}

// Alloc

const alignment = 8

// Alloc returns an aligned byte slice or nil.
func (b *Block) Alloc(size int) unsafe.Pointer {
	// Align start
	start := len(b.buf)
	start += (alignment - (start % alignment)) % alignment
	end := start + size
	if end > cap(b.buf) {
		return nil
	}

	// Grow buffer
	b.buf = b.buf[:end]

	// Slice buffer
	p := b.buf[start:end:end] // start:end:max, cap=max-start
	ptr := unsafe.Pointer(&p[0])
	return ptr
}

// Grow grows the buffer and allocates a byte slice.
func (b *Block) Grow(size int) []byte {
	free := cap(b.buf) - len(b.buf)
	if free < size {
		return nil
	}

	// Grow buffer
	start := len(b.buf)
	end := start + size
	b.buf = b.buf[:end]

	// Slice buffer
	return b.buf[start:end:end] // start:end:max, cap=max-start
}

// internal

func (b *Block) reset() {
	// Zero out
	for i := range b.buf {
		b.buf[i] = 0
	}

	b.buf = b.buf[:0]
}
