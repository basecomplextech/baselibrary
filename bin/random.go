// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import (
	"crypto/rand"
	"runtime"
	"sync"

	_ "unsafe"
)

var random = newRandomPool()

const (
	randomBuffer = 8192
	randomMulti  = 8
)

type randomPool struct {
	readers []randomReader
}

func newRandomPool() *randomPool {
	n := runtime.NumCPU() * randomMulti

	p := &randomPool{
		readers: make([]randomReader, n),
	}
	for i := range p.readers {
		p.readers[i].init(randomBuffer)
	}
	return p
}

func (p *randomPool) read64() [8]byte {
	i := int(fastrand()) % len(p.readers)
	return p.readers[i].read64()
}

func (p *randomPool) read128() [16]byte {
	i := int(fastrand()) % len(p.readers)
	return p.readers[i].read128()
}

func (p *randomPool) read256() [32]byte {
	i := int(fastrand()) % len(p.readers)
	return p.readers[i].read256()
}

// reader

type randomReader struct {
	mu  sync.Mutex
	buf []byte
	pos int

	_pad [216]byte // pad to cache line
}

func newRandomReader() *randomReader {
	r := &randomReader{}
	r.init(randomBuffer)
	return r
}

func (r *randomReader) init(size int) {
	r.buf = make([]byte, size)
	r.pos = size
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

// util

//go:linkname fastrand runtime.fastrand
func fastrand() uint32
