package memfs

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/complex1tech/baselibrary/filesys"
)

var _ filesys.File = (*memHandle)(nil)

type memHandle struct {
	path string

	mu    sync.RWMutex
	open  bool
	entry memEntry
}

func newMemHandle(path string, entry memEntry) *memHandle {
	return &memHandle{
		path:  path,
		open:  true,
		entry: entry,
	}
}

// Filename returns a file name, not a path as in *os.File.
func (h *memHandle) Filename() string {
	_, base := filepath.Split(h.path)
	return base
}

// Name returns a file path as in *os.File.
func (h *memHandle) Name() string {
	return h.path
}

// Path returns *os.File.Name() as presented to Open.
func (h *memHandle) Path() string {
	return h.path
}

// Size returns the file size in bytes.
func (h *memHandle) Size() (int64, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if !h.open {
		return 0, os.ErrClosed
	}

	return h.entry.Size()
}

// Methods

// Close closes the file.
func (h *memHandle) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.open {
		return nil
	}

	h.open = false
	h.entry = nil
	return nil
}

// Map maps the file into memory and returns its data.
func (h *memHandle) Map() ([]byte, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.open {
		return nil, os.ErrClosed
	}

	return h.entry.Map()
}

// Read reads data from the file into p.
func (h *memHandle) Read(p []byte) (n int, err error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if !h.open {
		return 0, os.ErrClosed
	}

	return h.entry.Read(p)
}

// ReadAt reads data from the file at an offset into p.
func (h *memHandle) ReadAt(p []byte, off int64) (n int, err error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if !h.open {
		return 0, os.ErrClosed
	}

	return h.entry.ReadAt(p, off)
}

// Readdir reads and returns the directory entries, upto n entries if n > 0.
func (h *memHandle) Readdir(count int) ([]os.FileInfo, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if !h.open {
		return nil, os.ErrClosed
	}

	return h.entry.Readdir(count)
}

// Readdirnames reads and returns the directory entries, upto n entries if n > 0.
func (h *memHandle) Readdirnames(n int) ([]string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if !h.open {
		return nil, os.ErrClosed
	}

	return h.entry.Readdirnames(n)
}

// Seek sets the file offset.
func (h *memHandle) Seek(offset int64, whence int) (int64, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.open {
		return 0, os.ErrClosed
	}

	return h.entry.Seek(offset, whence)
}

// Stat returns a file info.
func (h *memHandle) Stat() (os.FileInfo, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if !h.open {
		return nil, os.ErrClosed
	}

	return h.entry.Stat()
}

// Sync syncs the file to disk.
func (h *memHandle) Sync() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.open {
		return os.ErrClosed
	}

	return h.entry.Sync()
}

// Truncate truncates the file to a given length.
func (h *memHandle) Truncate(size int64) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.open {
		return os.ErrClosed
	}

	return h.entry.Truncate(size)
}

// Write writes data to the file at the current offset.
func (h *memHandle) Write(p []byte) (n int, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.open {
		return 0, os.ErrClosed
	}

	return h.entry.Write(p)
}

// WriteAt writes data to the file at an offset.
func (h *memHandle) WriteAt(p []byte, off int64) (n int, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.open {
		return 0, os.ErrClosed
	}

	return h.entry.WriteAt(p, off)
}

// WriteString writes a string to the file at the current offset.
func (h *memHandle) WriteString(s string) (ret int, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.open {
		return 0, os.ErrClosed
	}

	return h.entry.WriteString(s)
}
