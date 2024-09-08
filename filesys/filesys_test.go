// Copyright 2021 Ivan Korobkov. All rights reserved.

package filesys

import (
	"testing"
)

func TestFileSystem_TempFile__should_create_temp_file(t *testing.T) {
	fs := newFS()

	f, err := fs.TempFile("", "tmp-*")
	if err != nil {
		t.Fatal(err)
	}

	name := f.Name()
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	if err := fs.Remove(name); err != nil {
		t.Fatal(err)
	}
}
