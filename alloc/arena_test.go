package alloc

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArena_Alloc__should_allocate_data(t *testing.T) {
	a := newArena()
	d := a.Alloc(8)

	v := (*int64)(d)
	*v = math.MaxInt64

	assert.Equal(t, int64(math.MaxInt64), *v)
}
