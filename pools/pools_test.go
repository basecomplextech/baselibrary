package pools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPools(t *testing.T) {
	pp := New()

	a0 := Get[int32](pp)
	a1 := Get[int32](pp)
	b0 := Get[int64](pp)

	assert.Same(t, a0, a1)
	assert.NotSame(t, a0, b0)
}
