package buffer

// Buffer is a general purpose byte buffer interface.
type Buffer interface {
	// Len returns the number of bytes in the buffer; b.Len() == len(b.Bytes()).
	Len() int

	// Bytes returns a byte slice with the buffer bytes.
	// It is valid for use only until the next buffer mutation.
	Bytes() []byte

	// Grow grows the buffer and returns an n-byte slice.
	// It be should be used directly and is only valid until the next buffer mutation.
	//
	// Usage:
	//
	//	p := b.Grow(8)
	//	binary.BigEndian.PutUint64(p, 1234)
	//
	Grow(n int) []byte

	// Write appends bytes from p to the buffer.
	//
	// Equivalent to:
	//
	//	buf := b.Grow(n)
	//	copy(buf, p)
	//
	Write(p []byte) (n int, err error)

	// Reset resets the buffer to be empty, but can retain its internal buffer.
	Reset()
}

// New returns a new buffer and initializes it with a byte slice.
// The new buffer takes the ownership of the slice.
func New(buf []byte) Buffer {
	return newBuffer(buf)
}

// NewSize returns a new buffer and initializes it with a byte slice.
func NewSize(size int) Buffer {
	buf := make([]byte, 0, size)
	return newBuffer(buf)
}

type buffer struct {
	buf []byte
}

func newBuffer(buf []byte) *buffer {
	return &buffer{buf: buf}
}

// Len returns the number of bytes in the buffer; b.Len() == len(b.Bytes()).
func (b *buffer) Len() int {
	return len(b.buf)
}

// Bytes returns a byte slice with the buffer bytes.
func (b *buffer) Bytes() []byte {
	return b.buf
}

// Grow grows the buffer and returns an n-byte slice.
func (b *buffer) Grow(n int) []byte {
	cp := cap(b.buf)
	ln := len(b.buf)

	// realloc
	free := cp - ln
	if free < n {
		size := (cp * 2) + n
		buf := make([]byte, ln, size)
		copy(buf, b.buf)
		b.buf = buf
	}

	// grow buffer
	size := ln + n
	b.buf = b.buf[:size]

	// return slice
	return b.buf[ln:size]
}

// Write appends bytes from p to the buffer.
func (b *buffer) Write(p []byte) (n int, err error) {
	buf := b.Grow(len(p))

	n = copy(buf, p)
	return
}

// Reset resets the buffer to be empty, but can retain its internal buffer.
func (b *buffer) Reset() {
	b.buf = b.buf[:0]
}
