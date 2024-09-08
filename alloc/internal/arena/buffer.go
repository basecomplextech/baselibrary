// Copyright 2024 Ivan Korobkov. All rights reserved.

package arena

import "github.com/basecomplextech/baselibrary/buffer"

var _ buffer.Buffer = (*arenaBuffer)(nil)

type arenaBuffer struct {
	arena *arena
	buf   []byte
}

func (b *arenaBuffer) init(a *arena) {
	b.arena = a
}

// Len returns the number of bytes in the buffer; b.Len() == len(b.Bytes()).
func (b *arenaBuffer) Len() int {
	return len(b.buf)
}

// Bytes returns a byte slice with the buffer bytes.
// It is valid for use only until the next buffer mutation.
func (b *arenaBuffer) Bytes() []byte {
	return b.buf
}

// Grow grows the buffer and returns an n-byte slice.
// It be should be used directly and is only valid until the next buffer mutation.
//
// Usage:
//
//	p := b.Grow(8)
//	binary.BigEndian.PutUint64(p, 1234)
func (b *arenaBuffer) Grow(n int) []byte {
	cp := cap(b.buf)
	ln := len(b.buf)

	// Realloc
	free := cp - ln
	if free < n {
		size := (cp * 2) + n
		buf := b.arena.Bytes(size)[:ln:size]
		copy(buf, b.buf)
		b.buf = buf
	}

	// Grow buffer
	size := ln + n
	b.buf = b.buf[:size]

	// Return slice
	return b.buf[ln:size]
}

// Write appends bytes from p to the buffer.
//
// Equivalent to:
//
//	buf := b.Grow(n)
//	copy(buf, p)
func (b *arenaBuffer) Write(p []byte) (n int, err error) {
	buf := b.Grow(len(p))

	n = copy(buf, p)
	return
}

// WriteByte writes a byte to the buffer.
func (b *arenaBuffer) WriteByte(v byte) error {
	buf := b.Grow(1)
	buf[0] = v
	return nil
}

// Reset resets the buffer to be empty, but can retain its internal buffer.
func (b *arenaBuffer) Reset() {
	b.buf = b.buf[:0]
}
