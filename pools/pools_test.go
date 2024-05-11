package pools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPool(t *testing.T) {
	pp := New()

	// Same types
	a0 := GetPool[int32, int32](pp)
	a1 := GetPool[int32, int32](pp)

	// Same key types, different value types
	b0 := GetPool[int64, int32](pp)
	b1 := GetPool[int64, int64](pp)

	assert.Same(t, a0, a1)
	assert.NotSame(t, a0, b0)
	assert.NotSame(t, b0, b1)
}
