// Copyright 2022 Ivan Korobkov. All rights reserved.

package memfs

import (
	"github.com/basecomplextech/baselibrary/filesys"
	"github.com/basecomplextech/baselibrary/tests"
)

func testDir(t tests.T, fs filesys.FileSystem, path string) {
	if err := fs.MakePath(path, 0); err != nil {
		t.Fatal(err)
	}
}

func testFile(t tests.T, fs filesys.FileSystem, path string) filesys.File {
	f, err := fs.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	return f
}
