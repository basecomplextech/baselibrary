package fs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var _ File = (*testFile)(nil)

func TestFile() File {
	return newTestFile("", nil)
}

func TestFileBytes(b []byte) File {
	return newTestFile("", b)
}

func TestFileName(name string) File {
	return newTestFile(name, nil)
}

func TestFileNameBytes(name string, b []byte) File {
	return newTestFile(name, b)
}

type testFile struct {
	path string
	buf  *testBuf

	mu     sync.RWMutex
	closed bool
}

func newTestFile(path string, b []byte) *testFile {
	return &testFile{
		path: path,
		buf:  newTestBuf(b),
	}
}

func (f *testFile) Released() {
	f.Close()
}

// Filename returns a file name, not a path as in *os.File.
func (f *testFile) Filename() string {
	_, name := filepath.Split(f.path)
	return name
}

// Ext returns a file extension.
func (f *testFile) Ext() string {
	return filepath.Ext(f.path)
}

// Path returns *os.File.Name() as presented to Open.
func (f *testFile) Path() string {
	return f.path
}

// Map maps the file into memory and returns its data.
func (f *testFile) Map() ([]byte, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.closed {
		return nil, os.ErrClosed
	}
	return f.buf.buf, nil
}

// Size returns the file size in bytes.
func (f *testFile) Size() (int64, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	size := int64(f.buf.size())
	return size, nil
}

// os.File methods

func (f *testFile) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.closed {
		return nil
	}

	f.closed = true
	return nil
}

func (f *testFile) Name() string {
	return f.path
}

func (f *testFile) Read(b []byte) (n int, err error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.buf.read(b)
}

func (f *testFile) ReadAt(b []byte, off int64) (n int, err error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.buf.readAt(b, off)
}

func (f *testFile) Readdir(count int) ([]os.FileInfo, error) {
	panic("not implemented")
}

func (f *testFile) Readdirnames(n int) ([]string, error) {
	panic("not implemented")
}

func (f *testFile) Seek(offset int64, whence int) (ret int64, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.buf.seek(offset, whence)
}

func (f *testFile) SetDeadline(t time.Time) error {
	return nil
}

func (f *testFile) SetReadDeadline(t time.Time) error {
	return nil
}

func (f *testFile) SetWriteDeadline(t time.Time) error {
	return nil
}

func (f *testFile) Stat() (os.FileInfo, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	info := &testInfo{
		name: f.path,
		size: int64(f.buf.size()),
	}
	return info, nil
}

func (f *testFile) Sync() error {
	return nil
}

func (f *testFile) Truncate(size int64) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.buf.truncate(size)
}

func (f *testFile) Write(b []byte) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.buf.write(b)
}

func (f *testFile) WriteAt(b []byte, off int64) (n int, err error) {
	panic("not implemented")
}

func (f *testFile) WriteString(s string) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.buf.write([]byte(s))
}

// File info

var _ FileInfo = (*testInfo)(nil)

type testInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modtime time.Time
}

func (i *testInfo) Name() string       { return i.name }
func (i *testInfo) Size() int64        { return i.size }
func (i *testInfo) Mode() os.FileMode  { return i.mode }
func (i *testInfo) ModTime() time.Time { return i.modtime }
func (i *testInfo) IsDir() bool        { return false }
func (i *testInfo) Sys() interface{}   { return nil }

// File buffer

type testBuf struct {
	buf []byte
	off int
}

func newTestBuf(b []byte) *testBuf {
	return &testBuf{buf: b}
}

func (b *testBuf) read(p []byte) (n int, err error) {
	n = copy(p, b.buf[b.off:])
	b.off += n

	if n == 0 {
		return 0, io.EOF
	}
	return n, nil
}

func (b *testBuf) readAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, fmt.Errorf("fs: invalid read offset, off=%d", off)
	}

	unread := len(b.buf) - int(off)
	switch {
	case unread == 0:
		return 0, io.EOF
	case unread < len(p):
		return 0, io.ErrUnexpectedEOF
	}

	copy(p, b.buf[off:])
	return len(p), nil
}

func (b *testBuf) seek(offset int64, whence int) (ret int64, err error) {
	if whence != 0 {
		panic("unsupported whence")
	}

	if offset > int64(len(b.buf)) {
		return 0, errors.New("seek offset is out of range")
	}

	b.off = int(offset)
	return offset, nil
}

func (b *testBuf) size() int {
	return len(b.buf)
}

func (b *testBuf) truncate(size int64) error {
	if int64(len(b.buf)) < size {
		return errors.New("truncate size > file size")
	}

	b.buf = b.buf[:size]
	return nil
}

func (b *testBuf) write(p []byte) (n int, err error) {
	// return len(p), nil
	length := len(b.buf)
	newLength := length + len(p)

	if cap(b.buf) < newLength {
		size := cap(b.buf)*2 + len(p)
		buf := make([]byte, length, size)

		copy(buf, b.buf)
		b.buf = buf
	}

	b.buf = b.buf[:newLength]
	n = copy(b.buf[length:], p)
	return n, nil
}
