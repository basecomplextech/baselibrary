package memfs

import (
	"github.com/complex1tech/baselibrary/fs"
	"github.com/complex1tech/baselibrary/tests"
)

func testDir(t tests.T, fs fs.FileSystem, path string) {
	if err := fs.MakePath(path, 0); err != nil {
		t.Fatal(err)
	}
}

func testFile(t tests.T, fs fs.FileSystem, path string) fs.File {
	f, err := fs.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	return f
}
