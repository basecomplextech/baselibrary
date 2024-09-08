// Copyright 2024 Ivan Korobkov. All rights reserved.

package alloc

import (
	"testing"

	"github.com/basecomplextech/baselibrary/alloc/internal/arena"
	"github.com/stretchr/testify/assert"
)

func TestString__should_return_string_copy(t *testing.T) {
	a := arena.Test()
	s0 := "hello, world"
	s1 := String(a, s0)

	assert.Equal(t, s0, s1)
	assert.NotSame(t, s0, s1)
}

func TestString__should_alloc_string_copy(t *testing.T) {
	a := NewArena()
	s0 := "hello"
	s1 := String(a, s0)
	assert.Equal(t, s0, s1)
}

func TestStringBytes__should_alloc_string_copy(t *testing.T) {
	a := NewArena()
	s0 := []byte("hello")
	s1 := StringBytes(a, s0)
	assert.Equal(t, "hello", s1)
}

func TestStringRunes__should_alloc_string_copy(t *testing.T) {
	a := NewArena()
	s0 := []rune("hello")
	s1 := StringRunes(a, s0)
	assert.Equal(t, "hello", s1)
}

func TestStringJoin__should_join_strings_using_separator(t *testing.T) {
	a := NewArena()
	s0 := []string{"hello", "world"}
	s1 := StringJoin(a, s0, " ")
	assert.Equal(t, "hello world", s1)
}

func TestStringJoin__should_return_empty_string_when_src_is_empty(t *testing.T) {
	a := NewArena()
	s0 := []string{}
	s1 := StringJoin(a, s0, " ")
	assert.Equal(t, "", s1)
}

func TestStringJoin2__should_join_strings_using_separator(t *testing.T) {
	a := NewArena()
	s1 := StringJoin2(a, "hello", "world", " ")
	assert.Equal(t, "hello world", s1)
}

func TestStringFormat__should_format_string(t *testing.T) {
	a := NewArena()
	s1 := StringFormat(a, "%s %s", "hello", "world")
	assert.Equal(t, "hello world", s1)
}
