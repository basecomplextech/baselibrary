package bin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseString128(t *testing.T) {
	v0 := Random128()
	s := v0.String()

	v1, err := ParseString128(s)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, v0, v1)
}
