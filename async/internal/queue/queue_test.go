// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	q := New[int]()
	q.Push(1)
	q.Push(2)

	v, ok := q.Poll()
	require.True(t, ok)
	assert.Equal(t, 1, v)

	v, ok = q.Poll()
	require.True(t, ok)
	assert.Equal(t, 2, v)

	_, ok = q.Poll()
	require.False(t, ok)

	q.Clear()
}
