package fs

import (
	"testing"

	"github.com/epochtimeout/baselibrary/logging"
)

func TestFileSystem_TempFile__should_create_temp_file(t *testing.T) {
	logger := logging.Stderr
	fs := newFileSystem(logger)

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
