// Copyright 2022 Ivan Korobkov. All rights reserved.

package memfs

import "io"

type memBuffer struct {
	bufs [][]byte
}

func newMemBuffer() *memBuffer {
	return &memBuffer{}
}

func (b *memBuffer) bytes() []byte {
	if len(b.bufs) == 0 {
		return nil
	}
	if len(b.bufs) > 1 {
		b._merge()
	}

	return b.bufs[0]
}

func (b *memBuffer) read(p []byte, offset int) (int, error) {
	if len(b.bufs) == 0 {
		return 0, io.EOF
	}

	b._merge()
	buf := b.bufs[0]
	if offset >= len(buf) {
		return 0, io.EOF
	}

	n := copy(p, buf[offset:])
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}

func (b *memBuffer) write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	var last *[]byte
	var free int = 0
	if len(b.bufs) > 0 {
		last = &b.bufs[len(b.bufs)-1]
		free = cap(*last) - len(*last)
	}

	if free < len(p) {
		length := b.size() + len(p)
		buf := make([]byte, 0, length)
		b.bufs = append(b.bufs, buf)
		last = &b.bufs[len(b.bufs)-1]
	}

	*last = append(*last, p...)
	return len(p), nil
}

func (b *memBuffer) writeAt(p []byte, off int) (int, error) {
	size := b.size()
	if off >= size {
		return 0, io.EOF
	}

	b._merge()
	buf := b.bufs[0]
	n := copy(buf[off:], p)
	return n, nil
}

func (b *memBuffer) truncate(length int) {
	size := b.size()
	if length >= size {
		return
	}

	b._merge()

	buf := b.bufs[0]
	buf = buf[:length]
	b.bufs[0] = buf
}

func (b *memBuffer) size() int {
	var size int
	for _, buf := range b.bufs {
		size += len(buf)
	}
	return size
}

func (b *memBuffer) _merge() {
	if len(b.bufs) <= 1 {
		return
	}

	// Make buf
	size := b.size()
	data := make([]byte, 0, size)

	// Copy bufs
	for _, buf := range b.bufs {
		data = append(data, buf...)
	}

	// Done
	b.bufs = [][]byte{data}
}
