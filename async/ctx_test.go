package async

import (
	"testing"
	"time"

	"github.com/basecomplextech/baselibrary/status"
	"github.com/stretchr/testify/assert"
)

// Cancel

func TestContext_Cancel__should_cancel_context_close_done_channel(t *testing.T) {
	ctx := NewContext()

	go func() {
		time.Sleep(time.Millisecond * 5)
		ctx.Cancel()
	}()

	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled")
	}
	st := ctx.Status()
	assert.Equal(t, status.Cancelled, st)

	select {
	case <-ctx.Done():
	case <-time.After(time.Millisecond * 5):
		t.Fatal("done channel was not closed")
	}
	st = ctx.Status()
	assert.Equal(t, status.Cancelled, st)
}

func TestContext_Cancel__should_cancel_child_context(t *testing.T) {
	parent := NewContext()
	child := NextContextTimeout(parent, time.Second)

	go func() {
		time.Sleep(time.Millisecond * 5)
		parent.Cancel()
	}()

	select {
	case <-child.Done():
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled")
	}
	st := child.Status()
	assert.Equal(t, status.Cancelled, st)
}

// Timeout

func TestContext_Timeout__should_timeout_context(t *testing.T) {
	ctx := NewContextTimeout(time.Millisecond * 5)

	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled")
	}
	st := ctx.Status()
	assert.Equal(t, status.Timeout, st)

	select {
	case <-ctx.Done():
	case <-time.After(time.Millisecond * 5):
		t.Fatal("done channel was not closed")
	}
	st = ctx.Status()
	assert.Equal(t, status.Timeout, st)
}

func TestContext_Timeout__should_timeout_child_context(t *testing.T) {
	parent := NewContextTimeout(time.Millisecond * 5)
	child := NextContextTimeout(parent, time.Second)

	select {
	case <-child.Done():
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled")
	}

	st := child.Status()
	assert.Equal(t, status.Timeout, st)
}
