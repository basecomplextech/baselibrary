// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package alloc

import (
	"fmt"

	"github.com/basecomplextech/baselibrary/alloc/internal/arena"
	"github.com/basecomplextech/baselibrary/ref"
)

// Arena is an arena memory allocator.
// The arena must be freed after usage.
type (
	Arena      = arena.Arena
	MutexArena = arena.MutexArena
)

// NewArena returns a new non-thread-safe arena.
func NewArena() Arena {
	return arena.New()
}

// NewArenaRef returns a reference to a new non-thread-safe arena.
func NewArenaRef() ref.R[Arena] {
	return ref.New(NewArena())
}

// NewMutexArena returns a new thread-safe arena which uses a mutex to synchronize access.
func NewMutexArena() Arena {
	return arena.NewMutexArena()
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
	return arena.AcquireArena()
}

// Pin

// Pinned is a wrapper for an object pinned to an arena.
type Pinned[T any] struct {
	Set bool
	Obj T
}

// Pin pins an object to an arena.
func Pin[T any](arena Arena, obj T) Pinned[T] {
	arena.Pin(obj)
	return Pinned[T]{
		Set: true,
		Obj: obj,
	}
}

// Reset clears the pinned object and the set flag.
func (p *Pinned[T]) Reset() {
	var zero T
	p.Set = false
	p.Obj = zero
}

// Unwrap returns the pinned object and panics if the object is not pinned.
func (p Pinned[T]) Unwrap() T {
	if !p.Set {
		panic(fmt.Sprintf("no pinned object %T", p.Obj))
	}
	return p.Obj
}
