// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bufwriter

import (
	"io"

	"github.com/basecomplextech/baselibrary/alloc/internal/buffer"
)

// Writer buffers small writes and flushes them to an underlying writer.
type Writer interface {
	io.Writer

	// Len returns the number of buffered bytes.
	Len() int

	// Flush writes any buffered data to the underlying writer.
	Flush() error

	// Reset discards the unwritten data, clears the error and sets a new destination.
	Reset(io.Writer)

	// Internal

	// Free frees the writer, releases its internal resources.
	Free()
}

// New returns a new buffered writer with the default buffer size.
func New(dst io.Writer) Writer {
	buf := buffer.NewSize(defaultSize)
	return newWriter(dst, buf)
}

// NewSize returns a new buffered writer with the specified buffer size.
func NewSize(dst io.Writer, size int) Writer {
	buf := buffer.NewSize(size)
	return newWriter(dst, buf)
}

// internal

const defaultSize = 4096

type writer struct {
	dst io.Writer
	buf buffer.Buffer
	err error
}

func newWriter(dst io.Writer, buf buffer.Buffer) *writer {
	return &writer{
		dst: dst,
		buf: buf,
	}
}

// Len returns the number of buffered bytes.
func (w *writer) Len() int {
	return w.buf.Len()
}

// Flush writes any buffered data to the underlying writer.
func (w *writer) Flush() error {
	return w.flush()
}

// Reset discards the unwritten data, clears the error and sets a new destination.
func (w *writer) Reset(dst io.Writer) {
	w.dst = dst
	w.buf.Reset()
	w.err = nil
}

// Write writes len(p) bytes from p to the buffer, flushes the buffer if required.
func (w *writer) Write(p []byte) (int, error) {
	var n int

	for len(p) > 0 {
		rem := w.buf.Rem()
		if rem > 0 {
			m := min(rem, len(p))
			b := p[:m]
			p = p[m:]

			n1, err := w.buf.Write(b)
			n += n1
			if err != nil {
				w.err = err
				return n, err
			}
		}

		rem = w.buf.Rem()
		if rem == 0 {
			if err := w.flush(); err != nil {
				w.err = err
				return n, err
			}
		}
	}

	return n, nil
}

// Internal

// Free frees the writer, releases its internal resources.
func (w *writer) Free() {
	if w.buf != nil {
		b := w.buf
		w.buf = nil
		b.Free()
	}
}

// private

func (w *writer) flush() error {
	if w.err != nil {
		return w.err
	}
	if w.buf.Len() == 0 {
		return nil
	}

	b := w.buf.Bytes()
	_, err := w.dst.Write(b)
	if err != nil {
		w.err = err
		return err
	}

	w.buf.Reset()
	return nil
}
