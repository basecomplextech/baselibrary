package memfs

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemDir_Readdir__should_return_directory_entries(t *testing.T) {
	fs := newMemFS()

	testDir(t, fs, "dir/subdir")
	testFile(t, fs, "dir/file")

	dir, err := fs.Open("dir")
	if err != nil {
		t.Fatal(err)
	}

	infos, err := dir.Readdir(-1)
	if err != nil {
		t.Fatal(err)
	}

	names := []string{}
	for _, info := range infos {
		names = append(names, info.Name())
	}
	slices.Sort(names)

	exp := []string{"file", "subdir"}
	assert.Equal(t, exp, names)
}

func TestMemDir_Readdirnames__should_return_directory_entries(t *testing.T) {
	fs := newMemFS()

	testDir(t, fs, "dir/subdir")
	testFile(t, fs, "dir/file")

	dir, err := fs.Open("dir")
	if err != nil {
		t.Fatal(err)
	}

	names, err := dir.Readdirnames(-1)
	if err != nil {
		t.Fatal(err)
	}
	slices.Sort(names)

	exp := []string{"file", "subdir"}
	assert.Equal(t, exp, names)
}
