// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestRandomReader__should_have_size_equal_to_cache_line(t *testing.T) {
	r := randomReader{}
	size := unsafe.Sizeof(r)

	assert.Equal(t, 256, int(size))
}
