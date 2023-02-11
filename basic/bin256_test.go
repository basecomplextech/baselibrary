package basic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStringBin256(t *testing.T) {
	u0 := RandomBin256()
	s := u0.String()

	u1, err := ParseStringBin256(s)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, u0, u1)
}
