package disk

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInfo(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	info, err := GetInfo(path)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, info.Total > 0)
	assert.True(t, info.Free > 0)
	assert.True(t, info.Used > 0)
	assert.Equal(t, info.Used, info.Total-info.Free)
}
