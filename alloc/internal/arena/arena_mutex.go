// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package arena

import (
	"sync"
	"unsafe"

	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
	"github.com/basecomplextech/baselibrary/buffer"
)

// MutexArena is a thread-safe arena which uses a mutex to synchronize access.
type MutexArena = Arena

// NewMutexArena returns a new thread-safe arena which uses a mutex to synchronize access.
func NewMutexArena() Arena {
	return newMutexArena(heap.Global)
}

// internal

var _ Arena = (*mutexArena)(nil)

type mutexArena struct {
	mu sync.Mutex
	*state
}

func newMutexArena(h *heap.Heap) *mutexArena {
	a := &mutexArena{state: acquireState()}
	a.heap = h
	return a
}

// Cap returns the arena capacity.
func (a *mutexArena) Cap() int64 {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.cap
}

// Len calculates and returns the number of used bytes.
func (a *mutexArena) Len() int64 {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.len()
}

// Alloc allocates a memory block and returns a pointer to it.
func (a *mutexArena) Alloc(size int) unsafe.Pointer {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.alloc(size)
}

// Bytes allocates a byte slice.
func (a *mutexArena) Bytes(size int) []byte {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.bytes(size)
}

// Buffer allocates a buffer in the arena, the buffer cannot be freed.
func (a *mutexArena) Buffer() buffer.Buffer {
	b := Alloc[arenaBuffer](a)
	b.init(a)
	return b
}

// Pin pins an external object to the arena.
// The method is used to prevent the object from being collected by the garbage collector.
func (a *mutexArena) Pin(obj any) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.pin(obj)
}

// Reset resets the arena.
func (a *mutexArena) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.reset()
}

// Internal

// Free frees the arena and releases its memory.
// The method is not thread-safe and must be called only once.
func (a *mutexArena) Free() {
	a.mu.Lock()
	defer a.mu.Unlock()

	s := a.state
	a.state = nil
	releaseState(s)
}
