package alloc

import (
	"github.com/basecomplextech/baselibrary/alloc/internal/buf"
)

// Buffer is a buffer allocated by an allocator.
// The buffer must be freed after usage.
type Buffer = buf.Buffer

// NewBuffer allocates a buffer.
func NewBuffer() *Buffer {
	return buf.New()
}

// NewBuffer allocates a buffer of a preallocated capacity.
func NewBufferSize(size int) *Buffer {
	return buf.NewSize(size)
}

// AcquireBuffer returns a new buffer from the pool.
//
// The buffer must not be used or even referenced after Free.
// Use these method only when buffers do not escape an isolated scope.
//
// Typical usage:
//
//	buf := alloc.AcquireBuffer()
//	defer buf.Free() // free immediately
func AcquireBuffer() *Buffer {
	return buf.Acquire()
}
