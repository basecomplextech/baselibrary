package async

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStream__should_send_messages_to_listeners(t *testing.T) {
	src := NewStreamSource[string]()
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

	_, ok = q.Pop()
	assert.False(t, ok)
}

func TestStreamFilter__should_filter_messages(t *testing.T) {
	src := NewStreamSource[string]()
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

	_, ok = q.Pop()
	assert.False(t, ok)
}

func TestStreamParallel__should_send_messages_to_multiple_listeners(t *testing.T) {
	src := NewStreamSource[string]()
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

	_, ok = q0.Pop()
	assert.False(t, ok)

	s, ok = q1.Pop()
	assert.True(t, ok)
	assert.Equal(t, s, "goodbye")
}
