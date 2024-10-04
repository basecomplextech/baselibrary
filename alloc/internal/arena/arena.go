// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package arena

import (
	"unsafe"

	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
	"github.com/basecomplextech/baselibrary/buffer"
	"github.com/basecomplextech/baselibrary/pools"
)

// Arena is an arena allocator, which internally allocates memory in blocks.
//
// Arena is not thread-safe. If you need to use it in multiple goroutines, use NewMutexArena.
type Arena interface {
	// Cap returns the arena capacity.
	Cap() int64

	// Len calculates and returns the number of used bytes.
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

// New returns an empty non-thread-safe arena.
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

func newArena(h *heap.Heap) *arena {
	a := &arena{acquireState()}
	a.heap = h
	return a
}

// Cap returns the arena capacity.
func (a *arena) Cap() int64 {
	return a.cap
}

// Len calculates and returns the number of used bytes.
func (a *arena) Len() int64 {
	return a.len()
}

// Alloc allocates a memory block and returns a pointer to it.
func (a *arena) Alloc(size int) unsafe.Pointer {
	return a.alloc(size)
}

// Bytes allocates a byte slice.
func (a *arena) Bytes(size int) []byte {
	return a.bytes(size)
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
	a.pin(obj)
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
