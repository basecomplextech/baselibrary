// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package testfs

import (
	"os"
	"strings"

	"github.com/basecomplextech/baselibrary/filesys"
	"github.com/basecomplextech/baselibrary/filesys/memfs"
	"github.com/basecomplextech/baselibrary/tests"
)

const (
	TEST_FS   = "TEST_FS"
	TEST_DISK = "disk"
)

// Test returns a test file system, removes it on cleanup.
// The kind of file system is determined by the TEST_FS environment variable.
//
// To use a disk file system, use:
//
//	env TEST_FS=disk go test ...
func Test(t tests.T) (fs filesys.FileSystem, path string) {
	kind := os.Getenv(TEST_FS)
	kind = strings.ToLower(kind)

	if kind == TEST_DISK {
		return testDisk(t)
	}
	return testMemory(t)
}

func TestNoPath(t tests.T) filesys.FileSystem {
	fs, _ := Test(t)
	return fs
}

func testDisk(t tests.T) (filesys.FileSystem, string) {
	fs := filesys.New()
	path, err := fs.TempDir(".", "var_*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := fs.RemoveAll(path); err != nil {
			t.Error(err)
		}
	})
	return fs, path
}

func testMemory(t tests.T) (filesys.FileSystem, string) {
	fs := memfs.New()
	if err := fs.MakeDir("var", 0); err != nil {
		t.Fatal(err)
	}
	return fs, "var"
}
