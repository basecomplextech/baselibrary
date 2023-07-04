package ref

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBox(t *testing.T) {
	freed := false

	b := NewBoxedFunc(10, func() {
		freed = true
	})

	b.Retain()
	b.Release()

	v := b.Unwrap().Unwrap()
	require.Equal(t, 10, v)

	b.Release()
	assert.True(t, freed)
}
