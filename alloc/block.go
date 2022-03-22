package alloc

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

// reset resets the block.
func (b *block) reset() {
	// zero out
	for i := range b.buf {
		b.buf[i] = 0
	}

	b.buf = b.buf[:0]
}
