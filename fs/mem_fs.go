package fs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var _ FileSystem = (*memFS)(nil)

type memFS struct {
	mu   sync.RWMutex
	root *memDir
}

func newMemFS() *memFS {
	fs := &memFS{}
	fs.root = newMemDir(fs, "", nil)
	return fs
}

// Create creates a file in the file system.
func (fs *memFS) Create(path string) (File, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	return fs.create(path)
}

// MakeDir creates a directory in the file system.
func (fs *memFS) MakeDir(path string, perm os.FileMode) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	names := filepath.SplitList(path)
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

	names := filepath.SplitList(path)
	_, err := fs.root.makePath(names...)
	return err
}

// Open opens the named file for reading, see os.Open.
func (fs *memFS) Open(path string) (File, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	names := filepath.SplitList(path)
	f, ok := fs.root.findPath(names...)
	if !ok {
		return nil, os.ErrNotExist
	}
	if err := f.open(); err != nil {
		return nil, err
	}
	return f, nil
}

// OpenFile is the generalized open call, see os.OpenFile.
func (fs *memFS) OpenFile(path string, flag int, perm os.FileMode) (File, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	f, ok := fs.root.find(path)
	if !ok {
		if flag&os.O_CREATE != 0 {
			return fs.create(path)
		}
		return nil, os.ErrNotExist
	}
	if err := f.open(); err != nil {
		return nil, err
	}
	return f, nil
}

// Remove removes a file or an empty directory.
func (fs *memFS) Remove(path string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	return fs.root.remove(path)
}

// RemoveAll removes a path and any children it contains.
func (fs *memFS) RemoveAll(path string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	names := filepath.SplitList(path)
	return fs.root.removePath(names...)
}

// Rename renames a file or a directory.
func (fs *memFS) Rename(srcPath string, dstPath string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	// get source and its parent
	srcNames := filepath.SplitList(srcPath)
	if len(srcNames) == 0 {
		return os.ErrInvalid
	}
	src, ok := fs.root.find(srcPath)
	if !ok {
		return os.ErrNotExist
	}
	srcParent := src.getParent()

	// get dest parent
	dstNames := filepath.SplitList(dstPath)
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
func (fs *memFS) Stat(path string) (FileInfo, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	names := filepath.SplitList(path)
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

	if len(filepath.SplitList(pattern)) != 1 {
		return "", os.ErrInvalid
	}

	// find dir
	parent := fs.root
	if dir != "" {
		names := filepath.SplitList(dir)
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
	time := time.Now().UnixMilli()
	for i := 0; i < 10; i++ {
		name := strings.Replace(pattern, "*", fmt.Sprintf("%v", time), -1)
		_, ok := parent.find(name)
		if ok {
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
func (fs *memFS) TempFile(dir, pattern string) (File, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if len(filepath.SplitList(pattern)) != 1 {
		return nil, os.ErrInvalid
	}

	// find dir
	parent := fs.root
	if dir != "" {
		names := filepath.SplitList(dir)
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
	time := time.Now().UnixMilli()
	for i := 0; i < 10; i++ {
		name := strings.Replace(pattern, "*", fmt.Sprintf("%v", time), -1)
		_, ok := parent.find(name)
		if ok {
			continue
		}

		return parent.create(name)
	}

	return nil, errors.New("failed to create temp file, too many attempts")
}

// internal

func (fs *memFS) create(path string) (File, error) {
	names := filepath.SplitList(path)
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
