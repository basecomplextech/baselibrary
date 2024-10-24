// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import (
	"crypto/rand"
	"runtime"
	"sync"
	"sync/atomic"
)

const (
	randomBuffer      = 8192
	randomConcurrency = 8
)

var random = newRandomPool2()

// pool

type randomPool struct {
	readers []randomReader
}

func newRandomPool2() *randomPool {
	n := runtime.NumCPU() * randomConcurrency
	p := &randomPool{readers: make([]randomReader, n)}

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
	pos atomic.Int32
	mu  sync.Mutex
	buf []byte

	_ [216]byte // pad to cache line
}

func newRandomReader2() *randomReader {
	r := &randomReader{}
	r.init(randomBuffer)
	return r
}

func (r *randomReader) init(size int) {
	r.buf = make([]byte, size)
	r.pos.Store(int32(size))
}

func (r *randomReader) read64() (v [8]byte) {
	// Fast path
	end := int(r.pos.Add(8))
	if end <= len(r.buf) {
		copy(v[:], r.buf[end-8:end])
		return v
	}

	// Slow path
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check again
	end = int(r.pos.Add(8))
	if end <= len(r.buf) {
		copy(v[:], r.buf[end-8:end])
		return v
	}

	// Refill buffer
	rand.Read(r.buf)
	r.pos.Store(8)

	copy(v[:], r.buf[0:8])
	return v
}

func (r *randomReader) read128() (v [16]byte) {
	// Fast path
	end := int(r.pos.Add(16))
	if end <= len(r.buf) {
		copy(v[:], r.buf[end-16:end])
		return v
	}

	// Slow path
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check again
	end = int(r.pos.Add(8))
	if end <= len(r.buf) {
		copy(v[:], r.buf[end-16:end])
		return v
	}

	// Refill buffer
	rand.Read(r.buf)
	r.pos.Store(16)

	copy(v[:], r.buf[0:16])
	return v
}

func (r *randomReader) read256() (v [32]byte) {
	// Fast path
	end := int(r.pos.Add(32))
	if end <= len(r.buf) {
		copy(v[:], r.buf[end-32:end])
		return v
	}

	// Slow path
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check again
	end = int(r.pos.Add(32))
	if end <= len(r.buf) {
		copy(v[:], r.buf[end-32:end])
		return v
	}

	// Refill buffer
	rand.Read(r.buf)
	r.pos.Store(32)

	copy(v[:], r.buf[0:32])
	return v
}
