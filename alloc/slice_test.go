// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package alloc

import (
	"testing"

	"github.com/basecomplextech/baselibrary/alloc/internal/arena"
	"github.com/stretchr/testify/assert"
)

// Append

func TestAppend__should_append_item_to_slice(t *testing.T) {
	a := arena.Test()

	var s []int
	s = Append[[]int](a, s, 1)
	s = Append[[]int](a, s, 2)
	s = Append[[]int](a, s, 3)

	assert.Equal(t, []int{1, 2, 3}, s)
}

func TestAppendN__should_append_items_to_slice(t *testing.T) {
	a := arena.Test()

	var s []int
	s = AppendN[[]int](a, s, 1, 2, 3)

	assert.Equal(t, []int{1, 2, 3}, s)
}

func TestAppend__should_grow_slice_capacity(t *testing.T) {
	a := arena.Test()

	var s []int
	for i := range 4096 {
		s = Append[[]int](a, s, i)
	}

	assert.Equal(t, 4096, len(s))
}

// Copy

func TestCopy__should_copy_existing_slice_into_arena(t *testing.T) {
	type Value struct {
		A int64
		B int64
		C int64
	}

	v0 := []Value{
		{1, 2, 3},
		{10, 20, 30},
		{100, 200, 300},
	}

	a := arena.Test()
	v1 := Copy[[]Value](a, v0)

	assert.Equal(t, v0, v1)
}

// Grow

func TestGrow__should_grow_slice_capacity(t *testing.T) {
	a := arena.Test()

	s0 := []int{1, 2, 3}
	s1 := Grow[[]int](a, s0, 10)

	assert.Equal(t, 3, cap(s0))
	assert.Equal(t, 10, cap(s1))
	assert.Equal(t, []int{1, 2, 3}, s1)
	assert.NotSame(t, s0, s1)
}

// Slice

func TestSlice__should_allocate_slice(t *testing.T) {
	a := arena.Test()

	s := Slice[[]int](a, 16, 16)
	s[0] = 1
	s[1] = 2
	s[2] = 3
	s[3] = 4
	s[4] = 5

	assert.Equal(t, 1, s[0])
	assert.Equal(t, 2, s[1])
	assert.Equal(t, 3, s[2])
	assert.Equal(t, 4, s[3])
	assert.Equal(t, 5, s[4])

	assert.Equal(t, 16, len(s))
	assert.Equal(t, 16, cap(s))
}

func TestSlice1__should_allocate_slice_with_single_item(t *testing.T) {
	a := arena.Test()

	s := Slice1[[]int](a, 123)

	assert.Equal(t, 123, s[0])
	assert.Equal(t, 1, len(s))
	assert.Equal(t, 1, cap(s))
}
