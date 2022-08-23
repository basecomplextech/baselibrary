package bin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseString256(t *testing.T) {
	u0 := Random256()
	s := u0.String()

	u1, err := ParseString256(s)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, u0, u1)
}
