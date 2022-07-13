package memfs

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemFile_Map__should_return_file_data(t *testing.T) {
	fs := newMemFS()
	f := testFile(t, fs, "file")

	data := []byte("hello, world")
	if _, err := f.Write(data); err != nil {
		t.Fatal(err)
	}

	data1, err := f.Map()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, data, data1)
}

// Read

func TestMemFile_Read__should_read_data_increment_offset(t *testing.T) {
	fs := newMemFS()
	f := testFile(t, fs, "file")

	data := []byte("hello, world")
	if _, err := f.Write(data); err != nil {
		t.Fatal(err)
	}

	var data1 []byte
	p := make([]byte, 1)
loop:
	for {
		_, err := f.Read(p)
		switch {
		case err == io.EOF:
			break loop
		case err != nil:
			t.Fatal(err)
		}

		data1 = append(data1, p...)
	}

	assert.Equal(t, data, data1)
}

// ReadAt

func TestMemFile_ReadAt__should_read_data_at_offset(t *testing.T) {
	fs := newMemFS()
	f := testFile(t, fs, "file")

	data := []byte("hello, world")
	if _, err := f.Write(data); err != nil {
		t.Fatal(err)
	}

	p := make([]byte, 3)
	for i := 0; i < len(data); i++ {
		n, err := f.ReadAt(p[:], int64(i))
		switch {
		case err == io.EOF:
		case err != nil:
			t.Fatal(err)
		}

		d := data[i : i+n]
		p1 := p[:n]
		require.Equal(t, d, p1, i)
	}
}

// Seek

func TestMemFile_Seek__should_set_file_offset(t *testing.T) {
	fs := newMemFS()
	f := testFile(t, fs, "file")

	data := []byte("hello, world")
	if _, err := f.Write(data); err != nil {
		t.Fatal(err)
	}

	if _, err := f.Seek(7, 0); err != nil {
		t.Fatal(err)
	}

	data1 := make([]byte, 5)
	_, err := f.Read(data1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []byte("world"), data1)
}

// Truncate

func TestMemFile_Truncate__should_truncate_file_data(t *testing.T) {
	fs := newMemFS()
	f := testFile(t, fs, "file")

	data := []byte("hello, world")
	if _, err := f.Write(data); err != nil {
		t.Fatal(err)
	}

	if err := f.Truncate(5); err != nil {
		t.Fatal(err)
	}

	data1, err := f.Map()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []byte("hello"), data1)
}

// Write

func TestMemFile_Write__should_append_data(t *testing.T) {
	fs := newMemFS()
	f := testFile(t, fs, "file")

	data := []byte("hello, world;")
	if _, err := f.Write(data); err != nil {
		t.Fatal(err)
	}

	data1 := []byte("hello, world")
	if _, err := f.Write(data1); err != nil {
		t.Fatal(err)
	}

	data2, err := f.Map()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []byte("hello, world;hello, world"), data2)
}

// WriteAt

func TestMemFile_WriteAt__should_write_data_at_offset(t *testing.T) {
	fs := newMemFS()
	f := testFile(t, fs, "file")

	data := []byte("hello, world")
	if _, err := f.Write(data); err != nil {
		t.Fatal(err)
	}

	data1 := []byte("hello")
	if _, err := f.WriteAt(data1, 7); err != nil {
		t.Fatal(err)
	}

	data2, err := f.Map()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []byte("hello, hello"), data2)
}
