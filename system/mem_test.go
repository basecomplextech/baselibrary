// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemory__should_return_memory_info(t *testing.T) {
	info, err := Memory()
	if err != nil {
		t.Fatal(err)
	}

	assert.NotZero(t, info.Total)
	assert.NotZero(t, info.Free)
	assert.NotZero(t, info.Used)

	assert.Equal(t, info.Total-info.Free, info.Used)
}
