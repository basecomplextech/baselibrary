// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package alloc

import (
	"github.com/basecomplextech/baselibrary/alloc/internal/arena"
	"github.com/basecomplextech/baselibrary/buffer"
)

// BufferPool is a pool of buffers allocated in the arena.
// It is thread-safe but only if backed by [MutexArena].
// The pool itself is allocated in the arena.
type BufferPool interface {
	// Get returns an empty pool.
	Get() buffer.Buffer

	// Put reset and puts a buffer back into the pool.
	// The buffer must be allocated in this pool.
	Put(buf buffer.Buffer)
}

// NewBufferPool returns a new buffer pool which allocates buffers in the given arena.
func NewBufferPool(a Arena) BufferPool {
	return arena.NewBufferPool(a)
}
