package memfs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	fs := New()
	f, err := fs.Create("test")
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("hello, world")
	_, err = f.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	b, err := f.Map()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, data, b)
}
