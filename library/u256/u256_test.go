package u256

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseString(t *testing.T) {
	u0 := Random()
	s := u0.String()

	u1, err := ParseString(s)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, u0, u1)
}
