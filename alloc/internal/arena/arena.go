// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package arena

import (
	"unsafe"

	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
	"github.com/basecomplextech/baselibrary/buffer"
	"github.com/basecomplextech/baselibrary/collect/sets"
	"github.com/basecomplextech/baselibrary/pools"
)

// Arena is an arena allocator, which internally allocates memory in blocks.
//
// Arena is not thread-safe. If you need to use it in multiple goroutines,
// you must synchronize access or you may consider adding an AtomicArena wrapper.
type Arena interface {
	// Cap returns the arena capacity.
	Cap() int64

	// Len calculates and returns the arena used size.
	Len() int64

	// Methods

	// Alloc allocates a memory block and returns a pointer to it.
	Alloc(size int) unsafe.Pointer

	// Bytes allocates a byte slice.
	Bytes(size int) []byte

	// Buffer allocates a buffer in the arena, the buffer cannot be freed.
	Buffer() buffer.Buffer

	// Pin pins an external object to the arena.
	// The method is used to prevent the object from being collected by the garbage collector.
	Pin(obj any)

	// Reset resets the arena.
	Reset()

	// Internal

	// Free frees the arena and releases its memory.
	// The method is not thread-safe and must be called only once.
	Free()
}

// New returns an empty arena.
func New() Arena {
	return newArena(heap.Global)
}

// AcquireArena returns a pooled arena, which is released to the pool on Free.
//
// The arena must not be used or even referenced after Free.
// Use these method only when arenas do not escape an isolated scope.
//
// Typical usage:
//
//	arena := alloc.AcquireArena()
//	defer arena.Free() // free immediately
func AcquireArena() Arena {
	return acquireArena()
}

// internal

type arena struct {
	*state
}

type state struct {
	heap   *heap.Heap
	pooled bool
	cap    int64 // total allocated capacity

	blocks []*heap.Block
	pinned sets.Set[any]
}

func newArena(h *heap.Heap) *arena {
	a := &arena{acquireState()}
	a.heap = h
	return a
}

// Cap returns the arena capacity.
func (a *arena) Cap() int64 {
	return a.cap
}

// Len calculates and returns the arena used size.
func (a *arena) Len() int64 {
	n := int64(0)
	for _, block := range a.blocks {
		n += int64(block.Len())
	}
	return n
}

// Alloc allocates a memory block and returns a pointer to it.
func (a *arena) Alloc(size int) unsafe.Pointer {
	if len(a.blocks) > 0 {
		b := a.blocks[len(a.blocks)-1]

		ptr := b.Alloc(size)
		if ptr != nil {
			return ptr
		}
	}

	b := a.allocBlock(size)
	return b.Alloc(size)
}

// Bytes allocates a byte slice.
func (a *arena) Bytes(size int) []byte {
	if size == 0 {
		return nil
	}

	ptr := a.Alloc(size)
	return unsafe.Slice((*byte)(ptr), size)
}

// Buffer allocates a buffer in the arena, the buffer cannot be freed.
func (a *arena) Buffer() buffer.Buffer {
	b := Alloc[arenaBuffer](a)
	b.init(a)
	return b
}

// Pin pins an external object to the arena.
// The method is used to prevent the object from being collected by the garbage collector.
func (a *arena) Pin(obj any) {
	if a.pinned == nil {
		a.pinned = sets.New[any]()
	}

	a.pinned.Add(obj)
}

// Reset resets the arena.
func (a *arena) Reset() {
	a.reset()
}

// Internal

// Free frees the arena and releases its memory.
// The method is not thread-safe and must be called only once.
func (a *arena) Free() {
	if a.pooled {
		releaseArena(a)
		return
	}

	s := a.state
	a.state = nil
	releaseState(s)
}

// private

func (s *state) reset() {
	// Clear pinned objects
	if s.pinned != nil {
		clear(s.pinned)
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

// blocks

func (a *arena) allocBlock(n int) *heap.Block {
	// Double last block capacity
	size := 0
	if len(a.blocks) > 0 {
		last := a.blocks[len(a.blocks)-1]
		size = last.Cap() * 2
	}
	if n > size {
		size = n
	}

	// Alloc next block
	b := a.heap.Alloc(size)
	a.blocks = append(a.blocks, b)
	a.cap += int64(b.Cap())
	return b
}

// arena pool

var arenaPool = pools.NewPoolFunc(
	func() *arena {
		a := newArena(heap.Global)
		a.pooled = true
		return a
	},
)

func acquireArena() *arena {
	return arenaPool.New()
}

func releaseArena(a *arena) {
	a.reset()
	arenaPool.Put(a)
}

// state pool

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
