// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestMapShard__size_must_be_256(t *testing.T) {
	size := unsafe.Sizeof(mapShard[int, int]{})
	assert.Equal(t, 256, int(size))
}
