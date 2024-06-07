package byteq

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

	ok, st := q.Write(msg)
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

// Clear

func TestQueue_Clear__should_release_all_blocks(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)

	msg := bytes.Repeat([]byte("a"), 4096-4)
	testWrite(t, q, msg)
	testWrite(t, q, msg)
	testWrite(t, q, msg)
	testWrite(t, q, msg)

	q.Clear()
	assert.False(t, q.closed)
	assert.Nil(t, q.head)
	assert.Equal(t, 0, len(q.more))
}

// Close

func TestQueue_Close__should_close_queue(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)

	q.Close()
	assert.True(t, q.closed)
}

func TestQueue_Close__should_allow_reading_pending_messages_til_end(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)

	msg := []byte("hello, world")
	testWrite(t, q, msg)
	testWrite(t, q, msg)
	testWrite(t, q, msg)
	testWrite(t, q, msg)
	q.Close()

	testRead(t, q)
	testRead(t, q)
	testRead(t, q)
	testRead(t, q)

	_, ok, st := q.Read()
	assert.False(t, ok)
	assert.Equal(t, status.End, st)
}

func TestQueue_Close__should_notify_waiting_reader(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)

	wait := q.ReadWait()
	q.Close()

	select {
	case <-wait:
	default:
		t.Fatal("should notify waiting reader")
	}
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
	require.Equal(t, 2, len(q.more))

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
	require.Equal(t, 0, len(q.more))

	block := q.head
	assert.Equal(t, int32(0), block.readIndex)
	assert.Equal(t, int32(0), block.writeIndex)
}

func TestQueue_Read__should_release_read_blocks(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)

	msg := bytes.Repeat([]byte("a"), 4096-4)
	testWrite(t, q, msg)
	testWrite(t, q, msg)
	require.Equal(t, 1, len(q.more))

	testRead(t, q)
	testRead(t, q)
	require.Equal(t, 0, len(q.more))
}

func TestQueue_Read__should_notify_waiting_writer(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 1024)

	msg0 := bytes.Repeat([]byte("a"), 1024-4)
	msg1 := bytes.Repeat([]byte("a"), 4096-4)
	testWrite(t, q, msg0)
	testWrite(t, q, msg1)

	wait := q.WriteWait(len(msg1))
	testRead(t, q)

	_, ok, st := q.Read()
	if !st.OK() {
		t.Fatal(st)
	}
	assert.True(t, ok)

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

// ReadWait

func TestQueue_ReadWait__should_return_closed_chan_when_head_not_empty(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)

	msg := []byte("hello, world")
	testWrite(t, q, msg)

	wait := q.ReadWait()
	select {
	case <-wait:
	default:
		t.Fatal("should return closed chan")
	}
}

// Write

func TestQueue_Write__should_write_message(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)
	msg := []byte("hello, world")

	ok, st := q.Write(msg)
	if !st.OK() {
		t.Fatal(st)
	}
	require.True(t, ok)

	b := q.head.b.Bytes()[4 : 4+len(msg)]
	assert.Equal(t, msg, b)

	ok, st = q.Write(msg)
	if !st.OK() {
		t.Fatal(st)
	}
	assert.True(t, ok)
}

func TestQueue_Write__should_alloc_next_block(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)

	msg := bytes.Repeat([]byte("a"), 1024-4)
	testWrite(t, q, msg)
	testWrite(t, q, msg)

	assert.Equal(t, 1, len(q.more))
}

func TestQueue_Write__should_return_false_when_queue_full(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 1024)

	msg0 := bytes.Repeat([]byte("a"), 1024-4)
	msg1 := bytes.Repeat([]byte("a"), 4096-4)
	testWrite(t, q, msg0)
	testWrite(t, q, msg1)

	ok, st := q.Write(msg1)
	if !st.OK() {
		t.Fatal(st)
	}
	assert.False(t, ok)
}

func TestQueue_Write__should_return_false_when_all_blocks_full(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 4096)

	msg := bytes.Repeat([]byte("a"), 1024-4)
	testWrite(t, q, msg)
	testWrite(t, q, msg)
	testWrite(t, q, msg)
	testWrite(t, q, msg)

	ok, st := q.Write(msg)
	if !st.OK() {
		t.Fatal(st)
	}
	assert.False(t, ok)
}

func TestQueue_Write__should_notify_waiting_reader(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)

	msg := []byte("hello, world")
	wait := q.ReadWait()
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
	ok, st := q.Write(msg)
	assert.Equal(t, status.End, st)
	assert.False(t, ok)
}

// WriteWait

func TestQueue_WriteWait__should_return_closed_chan_when_space_available(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 1024)

	wait := q.WriteWait(1024 - 4)
	select {
	case <-wait:
	default:
		t.Fatal("should return closed chan")
	}
}

// Reset

func TestQueue_Reset__should_reset_queue(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)

	msg := []byte("hello, world")
	testWrite(t, q, msg)
	testWrite(t, q, msg)

	q.Close()
	q.Reset()

	assert.False(t, q.closed)
	assert.Nil(t, q.head)
	assert.Equal(t, 0, len(q.more))

	select {
	case <-q.readChan:
	default:
	}

	select {
	case <-q.writeChan:
	default:
	}
}

// Free

func TestQueue_Free__should_free_blocks(t *testing.T) {
	h := heap.New()
	q := newQueue(h, 0)

	msg := bytes.Repeat([]byte("a"), 1024-4)
	testWrite(t, q, msg)
	testWrite(t, q, msg)
	testWrite(t, q, msg)
	testWrite(t, q, msg)

	q.Free()

	assert.Nil(t, q.head)
	assert.Equal(t, 0, len(q.more))
}
