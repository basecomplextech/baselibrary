// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAtomic64(t *testing.T) {
	r := Atomic64{}
	r.Init(1)

	acquired := r.Acquire()
	require.True(t, acquired)

	released := r.Release()
	require.False(t, released)

	released = r.Release()
	require.True(t, released)

	acquired = r.Acquire()
	require.False(t, acquired)
}
