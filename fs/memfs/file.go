package memfs

import (
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/epochtimeout/baselibrary/fs"
)

// NewFile returns a new in-memory file detached from a file system.
func NewFile() fs.File {
	return newMemFile(nil, "", nil)
}

var _ fs.File = (*memFile)(nil)

type memFile struct {
	mu     *sync.RWMutex // shared fs.mu or new mutex when detached
	parent *memDir       // nil when detached

	name   string
	buffer *memBuffer

	refs    int
	offset  int
	opened  bool
	deleted bool
}

func newMemFile(fs *memFS, name string, parent *memDir) *memFile {
	return &memFile{
		mu:     fs.mu,
		parent: parent,

		name:   name,
		buffer: newMemBuffer(),

		opened: true,
	}
}

func newDetachedFile() *memFile {
	return &memFile{
		mu:     new(sync.RWMutex),
		parent: nil,

		name:   "",
		buffer: newMemBuffer(),

		opened: true,
	}
}

// Filename returns a file name, not a path as in *os.File.
func (f *memFile) Filename() string {
	return f.name
}

// Name returns a file path as in *os.File.
func (f *memFile) Name() string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.getPath()
}

// Path returns *os.File.Name() as presented to Open.
func (f *memFile) Path() string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.getPath()
}

// Size returns the file size in bytes.
func (f *memFile) Size() (int64, error) {
	f.mu.RUnlock()
	defer f.mu.RUnlock()

	size := f.buffer.size()
	return int64(size), nil
}

// Methods

// Close closes the file.
func (f *memFile) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.opened {
		return os.ErrClosed
	}

	f.refs--
	if f.refs > 0 {
		return nil
	}

	f.opened = false
	return nil
}

// Map maps the file into memory and returns its data.
func (f *memFile) Map() ([]byte, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.opened {
		return nil, os.ErrClosed
	}

	b := f.buffer.bytes()
	return b, nil
}

// Read reads data from the file into p.
func (f *memFile) Read(p []byte) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.opened {
		return 0, os.ErrClosed
	}

	n, err = f.buffer.read(p, f.offset)
	f.offset += n
	return n, err
}

// ReadAt reads data from the file at an offset into p.
func (f *memFile) ReadAt(p []byte, off int64) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.opened {
		return 0, os.ErrClosed
	}

	return f.buffer.read(p, int(off))
}

// Readdir reads and returns the directory entries, upto n entries if n > 0.
func (f *memFile) Readdir(count int) ([]os.FileInfo, error) {
	panic("not implemented")
}

// Readdirnames reads and returns the directory entries, upto n entries if n > 0.
func (f *memFile) Readdirnames(n int) ([]string, error) {
	panic("not implemented")
}

// Seek sets the file offset.
func (f *memFile) Seek(offset int64, whence int) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.opened {
		return 0, os.ErrClosed
	}
	if whence != 0 {
		panic("unsupported whence, only 0 is supported")
	}

	size := f.buffer.size()
	if offset < 0 || offset > int64(size) {
		return 0, errors.New("seek out of range")
	}

	f.offset = int(offset)
	return offset, nil
}

// Stat returns a file info.
func (f *memFile) Stat() (os.FileInfo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.opened {
		return nil, os.ErrClosed
	}

	info := f.getInfo()
	return info, nil
}

// Sync syncs the file to disk.
func (f *memFile) Sync() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.opened {
		return os.ErrClosed
	}
	return nil
}

// Truncate truncates the file to a given length.
func (f *memFile) Truncate(length int64) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.opened {
		return os.ErrClosed
	}

	f.buffer.truncate(int(length))
	if f.offset > int(length) {
		f.offset = int(length)
	}
	return nil
}

// Write writes data to the file at the current offset.
func (f *memFile) Write(p []byte) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.opened {
		return 0, os.ErrClosed
	}

	return f.buffer.write(p)
}

// WriteAt writes data to the file at an offset.
func (f *memFile) WriteAt(p []byte, off int64) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.opened {
		return 0, os.ErrClosed
	}

	return f.buffer.writeAt(p, int(off))
}

// WriteString writes a string to the file at the current offset.
func (f *memFile) WriteString(s string) (ret int, err error) {
	return f.Write([]byte(s))
}

// entry

func (f *memFile) isDir() bool {
	return false
}

func (f *memFile) isEmpty() bool {
	return f.buffer.size() == 0
}

func (f *memFile) getInfo() *memInfo {
	size := f.buffer.size()
	return &memInfo{
		name: f.name,
		size: int64(size),
		dir:  false,
	}
}

func (f *memFile) getName() string {
	return f.name
}

func (f *memFile) getPath() string {
	if f.parent == nil {
		return f.name
	}

	names := []string{f.name}
	for parent := f.parent; parent != nil; parent = parent.parent {
		names = append([]string{parent.name}, names...)
	}

	return filepath.Join(names...)
}

func (f *memFile) getParent() *memDir {
	return f.parent
}

func (f *memFile) open() error {
	if f.deleted {
		return os.ErrNotExist
	}

	f.refs++
	f.opened = true
	return nil
}

func (f *memFile) delete() error {
	f.parent = nil
	f.buffer = newMemBuffer()
	f.opened = false
	f.deleted = true
	return nil
}

func (f *memFile) move(newName string, newParent *memDir) error {
	f.name = newName
	f.parent = newParent
	return nil
}
