// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package arena

import (
	"unsafe"

	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
	"github.com/basecomplextech/baselibrary/collect/sets"
	"github.com/basecomplextech/baselibrary/opt"
	"github.com/basecomplextech/baselibrary/pools"
)

type state struct {
	heap   *heap.Heap
	pooled bool
	cap    int64 // total allocated capacity

	blocks []*heap.Block
	pinned opt.Opt[sets.Set[any]]
}

// len calculates and returns the number of used bytes.
func (s *state) len() int64 {
	n := int64(0)
	for _, block := range s.blocks {
		n += int64(block.Len())
	}
	return n
}

// alloc allocates a memory block and returns a pointer to it.
func (s *state) alloc(size int) unsafe.Pointer {
	if len(s.blocks) > 0 {
		b := s.blocks[len(s.blocks)-1]

		ptr := b.Alloc(size)
		if ptr != nil {
			return ptr
		}
	}

	b := s.allocBlock(size)
	return b.Alloc(size)
}

// bytes allocates a byte slice.
func (s *state) bytes(size int) []byte {
	if size == 0 {
		return nil
	}

	ptr := s.alloc(size)
	return unsafe.Slice((*byte)(ptr), size)
}

// pin pins an external object to the arena.
func (s *state) pin(obj any) {
	set, ok := s.pinned.Unwrap()
	if !ok {
		set = sets.New[any]()
		s.pinned.Set(set)
	}

	set.Add(obj)
}

// private

func (s *state) allocBlock(n int) *heap.Block {
	// Double last block capacity
	size := 0
	if len(s.blocks) > 0 {
		last := s.blocks[len(s.blocks)-1]
		size = last.Cap() * 2
	}
	if n > size {
		size = n
	}

	// Alloc next block
	b := s.heap.Alloc(size)
	s.blocks = append(s.blocks, b)
	s.cap += int64(b.Cap())
	return b
}

func (s *state) reset() {
	// Clear pinned objects
	if set, ok := s.pinned.Unwrap(); ok {
		clear(set)
	}

	// Return if no blocks
	if len(s.blocks) == 0 {
		return
	}

	// Reset capacity
	s.cap = 0

	// Reset the first block if small
	n := 0
	if b := s.blocks[0]; b.Cap() == heap.MinBlockSize {
		n = 1

		b.Reset()
		s.cap = int64(b.Cap())

		if len(s.blocks) == 1 {
			return
		}
	}

	// Free other blocks
	s.heap.FreeMany(s.blocks[n:]...)
	clear(s.blocks[n:]) // for gc
	s.blocks = s.blocks[:n]
}

// pool

var statePool = pools.NewPoolFunc(
	func() *state {
		return &state{}
	},
)

func acquireState() *state {
	return statePool.New()
}

func releaseState(s *state) {
	s.reset()
	statePool.Put(s)
}
