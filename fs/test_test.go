package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestFile_Map__should_return_data(t *testing.T) {
	f := TestFile()

	data := []byte("hello, world")
	if _, err := f.Write(data); err != nil {
		t.Fatal(err)
	}

	b, err := f.Map()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, data, b)
}
