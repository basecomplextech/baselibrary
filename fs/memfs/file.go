package memfs

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/epochtimeout/baselibrary/buffer"
	"github.com/epochtimeout/baselibrary/fs"
)

var _ fs.File = (*memFile)(nil)

type memFile struct {
	fs *memFS

	name   string
	parent *memDir
	buffer buffer.Buffer

	refs    int
	offset  int64
	opened  bool
	deleted bool
}

func newMemFile(fs *memFS, name string, parent *memDir) *memFile {
	return &memFile{
		fs: fs,

		name:   name,
		parent: parent,
		buffer: buffer.New(),

		opened: true,
	}
}

// Filename returns a file name, not a path as in *os.File.
func (f *memFile) Filename() string {
	return f.name
}

// Name returns a file path as in *os.File.
func (f *memFile) Name() string {
	f.fs.mu.RLock()
	defer f.fs.mu.RUnlock()

	return f.getPath()
}

// Path returns *os.File.Name() as presented to Open.
func (f *memFile) Path() string {
	f.fs.mu.RLock()
	defer f.fs.mu.RUnlock()

	return f.getPath()
}

// Size returns the file size in bytes.
func (f *memFile) Size() (int64, error) {
	f.fs.mu.RUnlock()
	defer f.fs.mu.RUnlock()

	size := f.buffer.Len()
	return int64(size), nil
}

// Methods

// Close closes the file.
func (f *memFile) Close() error {
	f.fs.mu.Lock()
	defer f.fs.mu.Unlock()

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
	f.fs.mu.Lock()
	defer f.fs.mu.Unlock()

	if !f.opened {
		return nil, os.ErrClosed
	}

	b := f.buffer.Bytes()
	return b, nil
}

// Read reads data from the file into p.
func (f *memFile) Read(p []byte) (n int, err error) {
	f.fs.mu.Lock()
	defer f.fs.mu.Unlock()

	if !f.opened {
		return 0, os.ErrClosed
	}

	b := f.buffer.Bytes()
	b = b[f.offset:]
	if len(b) == 0 {
		return 0, io.EOF
	}

	n = copy(p, b)
	f.offset += int64(n)
	return n, nil
}

// ReadAt reads data from the file at an offset into p.
func (f *memFile) ReadAt(p []byte, off int64) (n int, err error) {
	f.fs.mu.Lock()
	defer f.fs.mu.Unlock()

	if !f.opened {
		return 0, os.ErrClosed
	}

	b := f.buffer.Bytes()
	if off < 0 || off >= int64(len(b)) {
		return 0, io.ErrUnexpectedEOF
	}

	n = copy(p, b[off:])
	return n, nil
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
	f.fs.mu.Lock()
	defer f.fs.mu.Unlock()

	if !f.opened {
		return 0, os.ErrClosed
	}
	if whence != 0 {
		panic("unsupported whence, only 0 is supported")
	}

	size := f.buffer.Len()
	if offset < 0 || offset > int64(size) {
		return 0, errors.New("seek out of range")
	}

	f.offset = offset
	return offset, nil
}

// Stat returns a file info.
func (f *memFile) Stat() (os.FileInfo, error) {
	f.fs.mu.Lock()
	defer f.fs.mu.Unlock()

	if !f.opened {
		return nil, os.ErrClosed
	}

	info := f.getInfo()
	return info, nil
}

// Sync syncs the file to disk.
func (f *memFile) Sync() error {
	return nil
}

// Truncate truncates the file to a given length.
func (f *memFile) Truncate(size int64) error {
	f.fs.mu.Lock()
	defer f.fs.mu.Unlock()

	if !f.opened {
		return os.ErrClosed
	}

	ln := f.buffer.Len()
	if size < 0 || size > int64(ln) {
		return errors.New("truncate out of range")
	}

	data := f.buffer.Bytes()
	data1 := make([]byte, size)
	copy(data1, data)

	f.buffer = buffer.NewBytes(data1)
	return nil
}

// Write writes data to the file at the current offset.
func (f *memFile) Write(p []byte) (n int, err error) {
	f.fs.mu.Lock()
	defer f.fs.mu.Unlock()

	if !f.opened {
		return 0, os.ErrClosed
	}

	return f.buffer.Write(p)
}

// WriteAt writes data to the file at an offset.
func (f *memFile) WriteAt(p []byte, off int64) (n int, err error) {
	f.fs.mu.Lock()
	defer f.fs.mu.Unlock()

	if !f.opened {
		return 0, os.ErrClosed
	}

	size := f.buffer.Len()
	if size < 0 || off > int64(size) {
		return 0, errors.New("write out of range")
	}

	data := f.buffer.Bytes()
	n = copy(data[off:], p)
	return n, nil
}

// WriteString writes a string to the file at the current offset.
func (f *memFile) WriteString(s string) (ret int, err error) {
	return f.Write([]byte(s))
}

// entry

func (f *memFile) getInfo() *memInfo {
	size := f.buffer.Len()
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
	f.deleted = true
	return nil
}

func (f *memFile) move(newName string, newParent *memDir) error {
	f.name = newName
	f.parent = newParent
	return nil
}
