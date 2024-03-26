package async

import (
	"testing"
	"time"

	"github.com/basecomplextech/baselibrary/status"
	"github.com/stretchr/testify/assert"
)

func TestContext_Cancel_Free__should_not_race(t *testing.T) {
	n := 10_000

	for i := 0; i < n; i++ {
		ctx := TimeoutContext(1)
		go func() {
			ctx.Free()
		}()
	}
}

// Cancel

func TestContext_Cancel__should_cancel_context_close_done_channel(t *testing.T) {
	ctx := NewContext()

	go func() {
		time.Sleep(time.Millisecond * 5)
		ctx.Cancel()
	}()

	select {
	case <-ctx.Wait():
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled")
	}
	st := ctx.Status()
	assert.Equal(t, status.Cancelled, st)

	select {
	case <-ctx.Wait():
	case <-time.After(time.Millisecond * 5):
		t.Fatal("done channel was not closed")
	}
	st = ctx.Status()
	assert.Equal(t, status.Cancelled, st)
}

func TestContext_Cancel__should_cancel_child_context(t *testing.T) {
	parent := NewContext()
	child := NextTimeoutContext(parent, time.Second)

	go func() {
		time.Sleep(time.Millisecond * 5)
		parent.Cancel()
	}()

	select {
	case <-child.Wait():
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled")
	}
	st := child.Status()
	assert.Equal(t, status.Cancelled, st)
}

// Timeout

func TestContext_Timeout__should_timeout_context(t *testing.T) {
	ctx := TimeoutContext(time.Millisecond * 5)

	select {
	case <-ctx.Wait():
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled")
	}
	st := ctx.Status()
	assert.Equal(t, status.Timeout, st)

	select {
	case <-ctx.Wait():
	case <-time.After(time.Millisecond * 5):
		t.Fatal("done channel was not closed")
	}
	st = ctx.Status()
	assert.Equal(t, status.Timeout, st)
}

func TestContext_Timeout__should_timeout_child_context(t *testing.T) {
	parent := TimeoutContext(time.Millisecond * 5)
	child := NextTimeoutContext(parent, time.Second)

	select {
	case <-child.Wait():
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled")
	}

	st := child.Status()
	assert.Equal(t, status.Timeout, st)
}
