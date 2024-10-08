// Copyright 2021 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package memfs

import (
	"strings"
	"testing"

	"github.com/basecomplextech/baselibrary/filesys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemFS_Create__should_create_file(t *testing.T) {
	fs := newMemFS()

	f, err := fs.Create("file")
	if err != nil {
		t.Fatal(err)
	}

	name := f.Name()
	assert.Equal(t, "file", name)
}

func TestMemFS_Create__should_create_file_from_path(t *testing.T) {
	fs := newMemFS()
	testDir(t, fs, "dir")

	f, err := fs.Create("dir/file")
	if err != nil {
		t.Fatal(err)
	}

	name := f.Name()
	assert.Equal(t, "dir/file", name)
}

// Exists

func TestMemFS_Exists__should_return_true_when_entry_exists(t *testing.T) {
	fs := newMemFS()
	testDir(t, fs, "dir")
	testFile(t, fs, "dir/file").Close()

	ok, err := fs.Exists("dir/file")
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, ok)
}

func TestMemFS_Exists__should_return_false_when_entry_not_found(t *testing.T) {
	fs := newMemFS()
	testDir(t, fs, "dir")

	ok, err := fs.Exists("dir/file")
	if err != nil {
		t.Fatal(err)
	}

	assert.False(t, ok)
}

// MakeDir

func TestMemFS_MakeDir__should_create_directory(t *testing.T) {
	fs := newMemFS()

	if err := fs.MakeDir("dir", 0); err != nil {
		t.Fatal(err)
	}

	info, err := fs.Stat("dir")
	if err != nil {
		t.Fatal(err)
	}

	name := info.Name()
	assert.Equal(t, "dir", name)
	assert.True(t, info.IsDir())
}

func TestMemFS_MakeDir__should_create_directory_from_path(t *testing.T) {
	fs := newMemFS()
	testDir(t, fs, "dir")

	if err := fs.MakeDir("dir/subdir", 0); err != nil {
		t.Fatal(err)
	}

	info, err := fs.Stat("dir/subdir")
	if err != nil {
		t.Fatal(err)
	}

	name := info.Name()
	assert.Equal(t, "subdir", name)
	assert.True(t, info.IsDir())
}

// MakePath

func TestMemFS_MakePath__should_make_directories(t *testing.T) {
	fs := newMemFS()

	if err := fs.MakePath("dir/subdir/hello", 0); err != nil {
		t.Fatal(err)
	}

	info, err := fs.Stat("dir/subdir/hello")
	if err != nil {
		t.Fatal(err)
	}

	name := info.Name()
	assert.Equal(t, "hello", name)
	assert.True(t, info.IsDir())
}

func TestMemFS_MakePath__should_skip_existing_directories(t *testing.T) {
	fs := newMemFS()

	if err := fs.MakePath("dir", 0); err != nil {
		t.Fatal(err)
	}
	if err := fs.MakePath("dir/hello", 0); err != nil {
		t.Fatal(err)
	}
	if err := fs.MakePath("dir/subdir/hello/", 0); err != nil {
		t.Fatal(err)
	}

	info, err := fs.Stat("dir/subdir/hello")
	if err != nil {
		t.Fatal(err)
	}

	name := info.Name()
	assert.Equal(t, "hello", name)
	assert.True(t, info.IsDir())
}

// Open

func TestMemFS_Open__should_open_file(t *testing.T) {
	fs := newMemFS()
	testDir(t, fs, "dir/subdir")
	testFile(t, fs, "dir/subdir/file").Close()

	f, err := fs.Open("dir/subdir/file")
	if err != nil {
		t.Fatal(err)
	}

	name := f.Name()
	assert.Equal(t, "dir/subdir/file", name)
}

func TestMemFS_Open__should_open_directory(t *testing.T) {
	fs := newMemFS()
	testDir(t, fs, "dir/subdir")

	d, err := fs.Open("dir/subdir")
	if err != nil {
		t.Fatal(err)
	}

	info, err := d.Stat()
	if err != nil {
		t.Fatal(err)
	}
	name := info.Name()
	assert.Equal(t, "subdir", name)
	assert.True(t, info.IsDir())
}

// OpenFile

func TestMemFS_OpenFile__should_open_file(t *testing.T) {
	fs := newMemFS()
	testDir(t, fs, "dir")
	testFile(t, fs, "dir/file").Close()

	f, err := fs.OpenFile("dir/file", 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	name := f.Name()
	assert.Equal(t, "dir/file", name)
}

func TestMemFS_OpenFile__should_create_absent_file(t *testing.T) {
	fs := newMemFS()
	testDir(t, fs, "dir")

	f, err := fs.OpenFile("dir/file", filesys.O_CREATE, 0)
	if err != nil {
		t.Fatal(err)
	}

	name := f.Name()
	assert.Equal(t, "dir/file", name)
}

func TestMemFS_OpenFile__should_open_directory(t *testing.T) {
	fs := newMemFS()
	testDir(t, fs, "dir")

	d, err := fs.OpenFile("dir", 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	name := d.Name()
	assert.Equal(t, "dir", name)
}

// Remove

func TestMemFS_Remove__should_delete_file(t *testing.T) {
	fs := newMemFS()
	testDir(t, fs, "dir")
	testFile(t, fs, "dir/file").Close()

	if err := fs.Remove("dir/file"); err != nil {
		t.Fatal(err)
	}

	_, err := fs.Open("dir/file")
	assert.Equal(t, filesys.ErrNotExist, err)
}

func TestMemFS_Remove__should_delete_directory(t *testing.T) {
	fs := newMemFS()
	testDir(t, fs, "dir")
	testFile(t, fs, "dir/file").Close()

	if err := fs.Remove("dir/file"); err != nil {
		t.Fatal(err)
	}

	_, err := fs.Open("dir/file")
	assert.Equal(t, filesys.ErrNotExist, err)
}

// RemoveAll

func TestMemFS_RemoveAll__should_remove_directory_with_files(t *testing.T) {
	fs := newMemFS()
	testDir(t, fs, "dir/subdir")
	testFile(t, fs, "dir/subdir/file").Close()

	if err := fs.RemoveAll("dir/subdir"); err != nil {
		t.Fatal(err)
	}

	_, err := fs.Open("dir/subdir/file")
	assert.Equal(t, filesys.ErrNotExist, err)
}

// Rename

func TestMemFS_Rename__should_rename_file(t *testing.T) {
	fs := newMemFS()
	testDir(t, fs, "dir")
	testFile(t, fs, "dir/file").Close()

	if err := fs.Rename("dir/file", "dir/file2"); err != nil {
		t.Fatal(err)
	}

	_, err := fs.Open("dir/file")
	assert.Equal(t, filesys.ErrNotExist, err)

	f, err := fs.Open("dir/file2")
	if err != nil {
		t.Fatal(err)
	}

	path := f.Path()
	assert.Equal(t, "dir/file2", path)
}

func TestMemFS_Rename__should_replace_file(t *testing.T) {
	fs := newMemFS()
	testDir(t, fs, "dir")
	testFile(t, fs, "dir/file").Close()
	testFile(t, fs, "dir/file2").Close()

	if err := fs.Rename("dir/file", "dir/file2"); err != nil {
		t.Fatal(err)
	}

	_, err := fs.Open("dir/file")
	assert.Equal(t, filesys.ErrNotExist, err)

	f, err := fs.Open("dir/file2")
	if err != nil {
		t.Fatal(err)
	}

	path := f.Path()
	assert.Equal(t, "dir/file2", path)
}

func TestMemFS_Rename__should_move_file(t *testing.T) {
	fs := newMemFS()
	testDir(t, fs, "dir1")
	testDir(t, fs, "dir2")
	testFile(t, fs, "dir1/file").Close()

	if err := fs.Rename("dir1/file", "dir2/file"); err != nil {
		t.Fatal(err)
	}

	_, err := fs.Open("dir1/file")
	assert.Equal(t, filesys.ErrNotExist, err)

	f, err := fs.Open("dir2/file")
	if err != nil {
		t.Fatal(err)
	}

	path := f.Path()
	assert.Equal(t, "dir2/file", path)
}

// Stat

func TestMemFS_Stat__should_return_file_info(t *testing.T) {
	fs := newMemFS()
	testDir(t, fs, "dir")
	testFile(t, fs, "dir/file").Close()

	info, err := fs.Stat("dir/file")
	if err != nil {
		t.Fatal(err)
	}

	name := info.Name()
	assert.Equal(t, "file", name)
	assert.False(t, info.IsDir())
}

// TempDir

func TestMemFS_TempDir__should_create_temp_directory(t *testing.T) {
	fs := newMemFS()

	for i := 0; i < 10; i++ {
		dir, err := fs.TempDir("", "tmp-*")
		if err != nil {
			t.Fatal(err)
		}

		require.True(t, strings.HasPrefix(dir, "tmp-"), i)
	}
}

func TestMemFS_TempDir__should_create_temp_directory_in_custom_directory(t *testing.T) {
	fs := newMemFS()
	testDir(t, fs, "tmp")

	dir, err := fs.TempDir("tmp", "tmp-*")
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, strings.HasPrefix(dir, "tmp/tmp-"))
}

// TempFile

func TestMemFS_TempFile__should_create_temp_file(t *testing.T) {
	fs := newMemFS()

	for i := 0; i < 10; i++ {
		f, err := fs.TempFile("", "tmp-*")
		if err != nil {
			t.Fatal(err)
		}

		require.True(t, strings.HasPrefix(f.Path(), "tmp-"), i)
	}
}

func TestMemFS_TempFile__should_create_temp_file_in_custom_directory(t *testing.T) {
	fs := newMemFS()
	testDir(t, fs, "tmp")

	f, err := fs.TempFile("tmp", "tmp-*")
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, strings.HasPrefix(f.Path(), "tmp/tmp-"))
}
