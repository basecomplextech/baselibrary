// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import (
	"crypto/rand"
	"sync"
)

var random = newRandomReader()

type randomReader struct {
	mu  sync.Mutex
	buf []byte
	pos int
}

func newRandomReader() *randomReader {
	return &randomReader{
		buf: make([]byte, 4096),
		pos: 4096,
	}
}

func (r *randomReader) read64() [8]byte {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.pos+8 > len(r.buf) {
		r.fill()
	}

	var v [8]byte
	copy(v[:], r.buf[r.pos:])
	r.pos += 8
	return v
}

func (r *randomReader) read128() [16]byte {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.pos+16 > len(r.buf) {
		r.fill()
	}

	var v [16]byte
	copy(v[:], r.buf[r.pos:])
	r.pos += 16
	return v
}

func (r *randomReader) read256() [32]byte {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.pos+32 > len(r.buf) {
		r.fill()
	}

	var v [32]byte
	copy(v[:], r.buf[r.pos:])
	r.pos += 32
	return v

}

func (r *randomReader) fill() {
	rand.Read(r.buf)
	r.pos = 0
}
