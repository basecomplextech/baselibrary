package alloc

import (
	"unsafe"
)

const alignment = 8

type block struct {
	buf []byte
}

func newBlock(size int) *block {
	return &block{
		buf: make([]byte, 0, size),
	}
}

// len returns the number of used bytes in the block.
func (b *block) len() int {
	return len(b.buf)
}

// cap returns the block size in bytes.
func (b *block) cap() int {
	return cap(b.buf)
}

// free returns a free space in bytes.
func (b *block) free() int {
	return cap(b.buf) - len(b.buf)
}

// alloc grows an internal buffer and returns an aligned pointer to a memory of `size`.
func (b *block) alloc(size int) (unsafe.Pointer, bool) {
	// calc padding
	start := len(b.buf)
	start += (alignment - (start % alignment)) % alignment

	// return when not enough space
	free := cap(b.buf) - start
	if free < size {
		return nil, false
	}

	// grow buffer
	end := start + size
	b.buf = b.buf[:end]

	// slice buffer
	p := b.buf[start:end:end] // start:end:max (cap=max-start)
	ptr := unsafe.Pointer(&p[0])
	return ptr, true
}

// reset resets the block.
func (b *block) reset() {
	// zero out
	for i := range b.buf {
		b.buf[i] = 0
	}

	b.buf = b.buf[:0]
}
