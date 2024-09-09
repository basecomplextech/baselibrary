// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package memfs

import (
	"errors"
	"sync"

	"github.com/basecomplextech/baselibrary/filesys"
)

// NewFile returns a new in-memory file detached from a file system.
func NewFile() filesys.File {
	f := newDetachedFile()
	return newMemHandle("", f)
}

type memFile struct {
	mu     *sync.RWMutex // shared fs.mu or new mutex when detached
	parent *memDir       // nil when detached

	name   string
	buffer *memBuffer
	offset int
}

func newMemFile(fs *memFS, name string, parent *memDir) *memFile {
	return &memFile{
		mu:     fs.mu,
		parent: parent,

		name:   name,
		buffer: newMemBuffer(),
	}
}

func newDetachedFile() *memFile {
	return &memFile{
		mu:     new(sync.RWMutex),
		parent: nil,

		name:   "",
		buffer: newMemBuffer(),
	}
}

// Size returns the file size in bytes.
func (f *memFile) Size() (int64, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	size := f.buffer.size()
	return int64(size), nil
}

// Methods

// Close closes the file.
func (f *memFile) Close() error {
	panic("not implemented")
}

// Map maps the file into memory and returns its data.
func (f *memFile) Map() ([]byte, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	b := f.buffer.bytes()
	return b, nil
}

// Read reads data from the file into p.
func (f *memFile) Read(p []byte) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	n, err = f.buffer.read(p, f.offset)
	f.offset += n
	return n, err
}

// ReadAt reads data from the file at an offset into p.
func (f *memFile) ReadAt(p []byte, off int64) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.buffer.read(p, int(off))
}

// Readdir reads and returns the directory entries, upto n entries if n > 0.
func (f *memFile) Readdir(count int) ([]filesys.FileInfo, error) {
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
func (f *memFile) Stat() (filesys.FileInfo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	info := f.getInfo()
	return info, nil
}

// Sync syncs the file to disk.
func (f *memFile) Sync() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	return nil
}

// Truncate truncates the file to a given length.
func (f *memFile) Truncate(length int64) error {
	f.mu.Lock()
	defer f.mu.Unlock()

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

	return f.buffer.write(p)
}

// WriteAt writes data to the file at an offset.
func (f *memFile) WriteAt(p []byte, off int64) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

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

func (f *memFile) getParent() *memDir {
	return f.parent
}

func (f *memFile) open() error {
	f.offset = 0
	return nil
}

func (f *memFile) move(newName string, newParent *memDir) error {
	f.name = newName
	f.parent = newParent
	return nil
}
