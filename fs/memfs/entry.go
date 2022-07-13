package memfs

import "os"

type memEntry interface {
	Size() (int64, error)
	Map() ([]byte, error)
	Read(p []byte) (n int, err error)
	ReadAt(p []byte, off int64) (n int, err error)
	Readdir(count int) ([]os.FileInfo, error)
	Readdirnames(n int) ([]string, error)
	Seek(offset int64, whence int) (int64, error)
	Stat() (os.FileInfo, error)
	Sync() error
	Truncate(size int64) error
	Write(p []byte) (n int, err error)
	WriteAt(p []byte, off int64) (n int, err error)
	WriteString(s string) (ret int, err error)

	// internal

	isDir() bool
	isEmpty() bool

	getInfo() *memInfo
	getName() string
	getParent() *memDir

	move(newName string, newParent *memDir) error
}
