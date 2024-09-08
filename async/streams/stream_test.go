// Copyright 2024 Ivan Korobkov. All rights reserved.

package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStream__should_send_messages_to_listeners(t *testing.T) {
	src := NewSource[string]()
	q := src.Subscribe()

	src.Send("hello")
	src.Send("world")

	s, ok := q.Pop()
	assert.True(t, ok)
	assert.Equal(t, s, "hello")

	s, ok = q.Pop()
	assert.True(t, ok)
	assert.Equal(t, s, "world")

	_, ok = q.Pop()
	assert.False(t, ok)

	q.Free()
	src.Send("goodbye")
}

func TestStreamFilter__should_filter_messages(t *testing.T) {
	src := NewSource[string]()
	stream := src.Filter(func(s string) bool {
		return s == "hello"
	})

	q := stream.Subscribe()

	src.Send("hello")
	src.Send("world")

	s, ok := q.Pop()
	assert.True(t, ok)
	assert.Equal(t, s, "hello")

	_, ok = q.Pop()
	assert.False(t, ok)

	q.Free()
	src.Send("goodbye")
}

func TestStreamParallel__should_send_messages_to_multiple_listeners(t *testing.T) {
	src := NewSource[string]()
	q0 := src.Subscribe()
	q1 := src.Filter(func(s string) bool {
		return s != "world"
	}).Subscribe()

	src.Send("hello")
	src.Send("world")

	s, ok := q0.Pop()
	assert.True(t, ok)
	assert.Equal(t, s, "hello")

	s, ok = q1.Pop()
	assert.True(t, ok)
	assert.Equal(t, s, "hello")

	s, ok = q0.Pop()
	assert.True(t, ok)
	assert.Equal(t, s, "world")

	_, ok = q1.Pop()
	assert.False(t, ok)

	q0.Free()
	src.Send("goodbye")

	s, ok = q1.Pop()
	assert.True(t, ok)
	assert.Equal(t, s, "goodbye")
}
