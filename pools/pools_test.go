// Copyright 2024 Ivan Korobkov. All rights reserved.

package pools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPool(t *testing.T) {
	pp := NewPools()

	// Same types
	a0 := GetPool[int32](pp)
	a1 := GetPool[int32](pp)
	b0 := GetPool[int64](pp)

	assert.Same(t, a0, a1)
	assert.NotSame(t, a0, b0)
}
