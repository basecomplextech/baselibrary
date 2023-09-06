package ref

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRef(t *testing.T) {
	freed := false

	ref := NewFree[int](10, func() {
		freed = true
	})
	ref.Retain()

	v := ref.Unwrap()
	require.Equal(t, 10, v)

	ref.Release()
	require.False(t, freed)

	ref.Release()
	require.True(t, freed)
}
