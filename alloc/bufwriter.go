package alloc

import (
	"io"

	"github.com/basecomplextech/baselibrary/alloc/internal/bufwriter"
)

// Writer buffers small writes and flushes them to an underlying writer.
type BufferedWriter = bufwriter.Writer

// NewBufferedWriter returns a new buffered writer with the default buffer size.
func NewBufferedWriter(dst io.Writer) BufferedWriter {
	return bufwriter.New(dst)
}

// NewBufferedWriterSize returns a new buffered writer with the specified buffer size.
func NewBufferedWriterSize(dst io.Writer, size int) BufferedWriter {
	return bufwriter.NewSize(dst, size)
}
