package alloc

import (
	"github.com/basecomplextech/baselibrary/alloc/internal/buf"
)

// Buffer is a buffer allocated by an allocator.
// The buffer must be freed after usage.
type Buffer = buf.Buffer

// NewBuffer allocates a buffer in the global allocator.
func NewBuffer() *Buffer {
	return global.Buffer()
}

// NewBuffer allocates a buffer of a preallocated capacity in the global allocator.
func NewBufferSize(size int) *Buffer {
	return global.BufferSize(size)
}
