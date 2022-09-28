package memfs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/complex1tech/baselibrary/fs"
)

// New returns a new in-memory file system.
func New() fs.FileSystem {
	return newMemFS()
}

var _ fs.FileSystem = (*memFS)(nil)

type memFS struct {
	mu   *sync.RWMutex
	root *memDir
}

func newMemFS() *memFS {
	fs := &memFS{
		mu: new(sync.RWMutex),
	}
	fs.root = newMemDir(fs, "", nil)
	return fs
}

// Create creates a file in the file system.
func (fs *memFS) Create(path string) (fs.File, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	f, err := fs.create(path)
	if err != nil {
		return nil, err
	}

	h := newMemHandle(path, f)
	return h, nil
}

// Exists returns true if the file/directory exists.
func (fs *memFS) Exists(path string) (bool, error) {
	_, err := fs.Stat(path)
	switch {
	case os.IsNotExist(err):
		return false, nil
	case err != nil:
		return false, nil
	}
	return true, nil
}

// MakeDir creates a directory in the file system.
func (fs *memFS) MakeDir(path string, perm os.FileMode) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	names := splitPath(path)
	if len(names) == 0 {
		return os.ErrInvalid
	}

	parent, ok, err := fs.findParent(names...)
	if err != nil {
		return err
	}
	if !ok {
		return os.ErrNotExist
	}

	last := names[len(names)-1]
	_, err = parent.makeDir(last)
	return err
}

// MakePath creates a directory path and all parents that does not exist.
func (fs *memFS) MakePath(path string, perm os.FileMode) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	names := splitPath(path)
	_, err := fs.root.makePath(names...)
	return err
}

// Open opens the named file for reading, see os.Open.
func (fs *memFS) Open(path string) (fs.File, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	names := splitPath(path)
	f, ok := fs.root.findPath(names...)
	if !ok {
		return nil, os.ErrNotExist
	}
	if err := f.open(); err != nil {
		return nil, err
	}

	h := newMemHandle(path, f)
	return h, nil
}

// OpenFile is the generalized open call, see os.OpenFile.
func (fs *memFS) OpenFile(path string, flag int, perm os.FileMode) (fs.File, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	names := splitPath(path)
	f, ok := fs.root.findPath(names...)
	if ok {
		if err := f.open(); err != nil {
			return nil, err
		}

		h := newMemHandle(path, f)
		return h, nil
	}

	if flag&os.O_CREATE == 0 {
		return nil, os.ErrNotExist
	}

	f, err := fs.create(path)
	if err != nil {
		return nil, err
	}

	h := newMemHandle(path, f)
	return h, nil
}

// Remove removes a file or an empty directory.
func (fs *memFS) Remove(path string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	names := splitPath(path)
	parent, ok, err := fs.findParent(names...)
	switch {
	case err != nil:
		return err
	case !ok:
		return nil
	}

	name := names[len(names)-1]
	e, ok := parent.find(name)
	if !ok {
		return nil
	}

	if e.isDir() {
		dir := e.(*memDir)
		if !dir.isEmpty() {
			return errors.New("directory not empty")
		}
	}

	return parent.remove(name)
}

// RemoveAll removes a path and any children it contains.
func (fs *memFS) RemoveAll(path string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	names := splitPath(path)
	parent, ok, err := fs.findParent(names...)
	switch {
	case err != nil:
		return err
	case !ok:
		return nil
	}

	name := names[len(names)-1]
	return parent.remove(name)
}

// Rename renames a file or a directory.
func (fs *memFS) Rename(srcPath string, dstPath string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	// get source and its parent
	srcNames := splitPath(srcPath)
	if len(srcNames) == 0 {
		return os.ErrInvalid
	}
	src, ok := fs.root.findPath(srcNames...)
	if !ok {
		return os.ErrNotExist
	}
	srcParent := src.getParent()

	// get dest parent
	dstNames := splitPath(dstPath)
	if len(dstNames) == 0 {
		return os.ErrInvalid
	}
	dstParent, ok, err := fs.findParent(dstNames...)
	switch {
	case err != nil:
		return err
	case !ok:
		return os.ErrNotExist
	}
	dstName := dstNames[len(dstNames)-1]

	// move source to dest
	srcParent.removeEntry(src)
	dstParent.addEntry(src, dstName)
	return nil
}

// Stat returns a file info.
func (fs *memFS) Stat(path string) (fs.FileInfo, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	names := splitPath(path)
	f, ok := fs.root.findPath(names...)
	if !ok {
		return nil, os.ErrNotExist
	}

	info := f.getInfo()
	return info, nil
}

// MkdirTemp creates a new temporary directory in the directory dir
// and returns the pathname of the new directory.
func (fs *memFS) TempDir(dir, pattern string) (path string, err error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if len(splitPath(pattern)) != 1 {
		return "", os.ErrInvalid
	}
	dir = cleanPath(dir)

	// find dir
	parent := fs.root
	if dir != "" {
		names := splitPath(dir)
		entry, ok := fs.root.findPath(names...)
		if !ok {
			return "", os.ErrNotExist
		}

		parent, ok = entry.(*memDir)
		if !ok {
			return "", errors.New("path is not a directory")
		}
	}

	// create file
	time := time.Now().UnixNano()
	for i := 0; i < 100; i++ {
		name := strings.Replace(pattern, "*", fmt.Sprintf("%v", time), -1)

		_, ok := parent.find(name)
		if ok {
			time++
			continue
		}

		_, err := parent.makeDir(name)
		if err != nil {
			return "", err
		}
		return filepath.Join(dir, name), nil
	}

	return "", errors.New("failed to create temp directory, too many attempts")
}

// TempFile creates a new temporary file in the directory dir,
// opens the file for reading and writing, and returns the resulting *os.File.
func (fs *memFS) TempFile(dir, pattern string) (fs.File, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if len(splitPath(pattern)) != 1 {
		return nil, os.ErrInvalid
	}
	dir = cleanPath(dir)

	// find dir
	parent := fs.root
	if dir != "" {
		names := splitPath(dir)
		entry, ok := fs.root.findPath(names...)
		if !ok {
			return nil, os.ErrNotExist
		}

		parent, ok = entry.(*memDir)
		if !ok {
			return nil, errors.New("path is not a directory")
		}
	}

	// create file
	time := time.Now().UnixNano()
	for i := 0; i < 100; i++ {
		name := strings.Replace(pattern, "*", fmt.Sprintf("%v", time), -1)

		_, ok := parent.find(name)
		if ok {
			time++
			continue
		}

		f, err := parent.create(name)
		if err != nil {
			return nil, err
		}

		path := filepath.Join(dir, name)
		h := newMemHandle(path, f)
		return h, nil
	}

	return nil, errors.New("failed to create temp file, too many attempts")
}

// internal

func (fs *memFS) create(path string) (memEntry, error) {
	names := splitPath(path)
	if len(names) == 0 {
		return nil, os.ErrInvalid
	}

	parent, ok, err := fs.findParent(names...)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, os.ErrNotExist
	}

	last := names[len(names)-1]
	return parent.create(last)
}

func (fs *memFS) findParent(names ...string) (*memDir, bool, error) {
	if len(names) <= 1 {
		return fs.root, true, nil
	}

	last, ok := fs.root.findPath(names[:len(names)-1]...)
	if !ok {
		return nil, false, nil
	}

	dir, ok := last.(*memDir)
	if !ok {
		return nil, false, errors.New("path is not a directory")
	}

	return dir, true, nil
}
