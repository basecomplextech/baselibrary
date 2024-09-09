// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package byteq

import (
	"encoding/binary"
	"sync/atomic"

	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
	"github.com/basecomplextech/baselibrary/pools"
)

type block struct {
	b *heap.Block

	readIndex  int32 // next read start, mutated by reader
	writeIndex int32 // last write end, mutated by writer, can be loaded atomically by reader
}

func newBlock(b *heap.Block) *block {
	b.Grow(b.Cap()) // use all available space

	bb := acquireBlock()
	bb.b = b
	return bb
}

func (b *block) loadWriteIndex() int32 {
	return atomic.LoadInt32(&b.writeIndex)
}

func (b *block) storeWriteIndex(wi int32) {
	// Paranoid check
	if wi < b.writeIndex {
		panic("write index is less than current write index")
	}
	atomic.StoreInt32(&b.writeIndex, wi)
}

// read

// read reads the next message and increments the read index.
// the method is called by a single reader inside the read lock.
func (b *block) read() []byte {
	p := b.b.Bytes()
	p = p[b.readIndex:]

	size := binary.BigEndian.Uint32(p)
	msg := p[4 : 4+size]

	b.readIndex += 4 + int32(size)

	// Paranoid check
	wi := b.loadWriteIndex()
	if b.readIndex > wi {
		panic("read index is greater than write index")
	}
	return msg
}

// copy

// copy copies the message to the block, and returns the next write index.
// the method is called by a single writer inside the write lock.
func (b *block) copy(msg []byte) int32 {
	size := len(msg)
	n := 4 + size
	wi := b.writeIndex + int32(n)

	p := b.b.Bytes()
	p = p[b.writeIndex:wi]

	binary.BigEndian.PutUint32(p, uint32(size))
	copy(p[4:], msg)

	return wi
}

// guarded by queue.mu

// cap returns the block capacity.
func (b *block) cap() int {
	return b.b.Cap()
}

// rem returns the remaining space in bytes.
func (b *block) rem() int {
	cp := b.b.Cap()
	wi := int(b.writeIndex)
	return cp - wi
}

// reset resets the block.
func (b *block) reset() {
	b.readIndex = 0
	b.writeIndex = 0
}

func (b *block) free(heap *heap.Heap) {
	heap.Free(b.b)
	b.b = nil
}

// pool

var pool = pools.NewPoolFunc(func() *block {
	return &block{}
})

func acquireBlock() *block {
	return pool.New()
}

func releaseBlock(b *block) {
	*b = block{}
	pool.Put(b)
}
