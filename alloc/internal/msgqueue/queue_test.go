package msgqueue

import (
	"testing"

	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueue_Write_Read__should_write_and_read_message(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)

	msg0 := []byte("hello, world")
	ok, st := q.Write(msg0)
	if !st.OK() {
		t.Fatal(st)
	}
	require.True(t, ok)

	msg1, ok, st := q.Read()
	if !st.OK() {
		t.Fatal(st)
	}
	assert.True(t, ok)
	assert.Equal(t, msg0, msg1)
}
