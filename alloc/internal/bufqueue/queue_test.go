package bufqueue

import (
	"bytes"
	"testing"

	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
	"github.com/basecomplextech/baselibrary/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testWrite(t *testing.T, q *queue, msg []byte) {
	t.Helper()

	ok, _, st := q.Write(msg)
	if !st.OK() {
		t.Fatal(st)
	}
	if !ok {
		t.Fatal("write failed")
	}
}

func testRead(t *testing.T, q *queue) []byte {
	t.Helper()

	msg, ok, st := q.Read()
	if !st.OK() {
		t.Fatal(st)
	}
	if !ok {
		t.Fatal("read failed")
	}
	return msg
}

// Queue

func TestQueue__should_write_and_read_message(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)

	msg0 := []byte("hello, world")
	ok, _, st := q.Write(msg0)
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

// Read

func TestQueue_Read__should_read_message(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)

	msg := []byte("hello, world")
	testWrite(t, q, msg)

	msg1, ok, st := q.Read()
	if !st.OK() {
		t.Fatal(st)
	}
	assert.True(t, ok)
	assert.Equal(t, msg, msg1)
}

func TestQueue_Read__should_return_false_when_no_messages(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)

	msg, ok, st := q.Read()
	if !st.OK() {
		t.Fatal(st)
	}
	assert.False(t, ok)
	assert.Nil(t, msg)
}

func TestQueue_Read__should_return_false_when_no_unread_messages(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)

	msg := []byte("hello, world")
	testWrite(t, q, msg)
	testRead(t, q)

	msg1, ok, st := q.Read()
	if !st.OK() {
		t.Fatal(st)
	}
	assert.False(t, ok)
	assert.Nil(t, msg1)
}

func TestQueue_Read__should_return_false_when_no_unread_messages_in_all_blocks(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)

	msg := bytes.Repeat([]byte("a"), 4096-4)
	testWrite(t, q, msg)
	testWrite(t, q, msg)
	testWrite(t, q, msg)
	testWrite(t, q, msg)
	require.Equal(t, 3, len(q.blocks))

	testRead(t, q)
	testRead(t, q)
	testRead(t, q)
	testRead(t, q)

	_, ok, st := q.Read()
	if !st.OK() {
		t.Fatal(st)
	}
	assert.False(t, ok)
}

func TestQueue_Read__should_reset_only_block_when_all_messages_read(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)

	msg := []byte("hello, world")
	testWrite(t, q, msg)
	testRead(t, q)

	_, ok, st := q.Read()
	if !st.OK() {
		t.Fatal(st)
	}
	assert.False(t, ok)
	require.Equal(t, 1, len(q.blocks))

	block := q.blocks[0]
	assert.Equal(t, 0, block.read)
	assert.False(t, block.started)
	assert.Equal(t, 0, block.Len())
}

func TestQueue_Read__should_release_read_blocks(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)

	msg := bytes.Repeat([]byte("a"), 4096-4)
	testWrite(t, q, msg)
	testWrite(t, q, msg)
	require.Equal(t, 2, len(q.blocks))

	testRead(t, q)
	testRead(t, q)
	require.Equal(t, 1, len(q.blocks))
}

func TestQueue_Read__should_notify_waiting_writer(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 1024)

	msg := bytes.Repeat([]byte("a"), 1024-4)
	testWrite(t, q, msg)

	wait := q.WaitCanWrite(len(msg))
	testRead(t, q)

	_, ok, st := q.Read()
	if !st.OK() {
		t.Fatal(st)
	}
	assert.False(t, ok)

	select {
	case <-wait:
	default:
		t.Fatal("should notify waiting writer")
	}
}

func TestQueue_Read__should_read_existing_messages_when_queue_closed(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)
	msg := []byte("hello, world")

	testWrite(t, q, msg)
	testWrite(t, q, msg)
	q.Close()

	testRead(t, q)
	testRead(t, q)

	_, ok, st := q.Read()
	assert.Equal(t, status.End, st)
	assert.False(t, ok)
}

// Write

func TestQueue_Write__should_write_message(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)
	msg := []byte("hello, world")

	ok, wasEmpty, st := q.Write(msg)
	if !st.OK() {
		t.Fatal(st)
	}
	require.True(t, ok)
	assert.True(t, wasEmpty)

	b := q.blocks[0].Bytes()
	assert.Equal(t, msg, b[4:])

	ok, wasEmpty, st = q.Write(msg)
	if !st.OK() {
		t.Fatal(st)
	}
	assert.True(t, ok)
	assert.False(t, wasEmpty)
}

func TestQueue_Write__should_alloc_next_block(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)

	msg := bytes.Repeat([]byte("a"), 1024-4)
	testWrite(t, q, msg)
	testWrite(t, q, msg)

	assert.Equal(t, 2, len(q.blocks))
}

func TestQueue_Write__should_return_false_when_queue_full(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 1024)

	msg := bytes.Repeat([]byte("a"), 1024-4)
	testWrite(t, q, msg)

	ok, _, st := q.Write(msg)
	if !st.OK() {
		t.Fatal(st)
	}
	assert.False(t, ok)
}

func TestQueue_Write__should_return_false_when_all_blocks_full(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 1024+4096)

	msg := bytes.Repeat([]byte("a"), 1024-4)
	testWrite(t, q, msg)
	testWrite(t, q, msg)
	testWrite(t, q, msg)
	testWrite(t, q, msg)
	testWrite(t, q, msg)

	ok, _, st := q.Write(msg)
	if !st.OK() {
		t.Fatal(st)
	}
	assert.False(t, ok)
}

func TestQueue_Write__should_notify_waiting_reader(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)

	msg := []byte("hello, world")
	wait := q.Wait()
	testWrite(t, q, msg)

	select {
	case <-wait:
	default:
		t.Fatal("should notify waiting reader")
	}
}

func TestQueue_Write__should_return_error_when_queue_closed(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)
	q.Close()

	msg := []byte("hello, world")
	ok, _, st := q.Write(msg)
	assert.Equal(t, status.End, st)
	assert.False(t, ok)
}
