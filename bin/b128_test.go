package bin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStringBin128(t *testing.T) {
	v0 := RandomBin128()
	s := v0.String()

	v1, err := ParseStringBin128(s)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, v0, v1)
}
