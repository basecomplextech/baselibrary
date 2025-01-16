// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package context

import (
	"testing"
	"time"

	"github.com/basecomplextech/baselibrary/status"
	"github.com/stretchr/testify/assert"
)

// New/Free

func TestContext_New_Free__should_return_state_to_pool(t *testing.T) {
	allocs := testing.AllocsPerRun(1000, func() {
		ctx := New()
		ctx.Free()
	})

	assert.Equal(t, 1, int(allocs))
}

// Cancel/Free

func TestContext_Cancel_Free__should_not_race(t *testing.T) {
	n := 100_000

	for i := 0; i < n; i++ {
		ctx := Timeout(1)
		go func() {
			ctx.Free()
		}()
	}
}

// Cancel

func TestContext_Cancel__should_cancel_context_close_done_channel(t *testing.T) {
	ctx := New()

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
	parent := New()
	child := NextTimeout(parent, time.Second)

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
	ctx := Timeout(time.Millisecond * 5)

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
	parent := Timeout(time.Millisecond * 5)
	child := NextTimeout(parent, time.Second)

	select {
	case <-child.Wait():
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled")
	}

	st := child.Status()
	assert.Equal(t, status.Timeout, st)
}

// Next

func TestNextContext__should_not_deadlock_when_parent_cancelled_already(t *testing.T) {
	ctx := newContext(nil)
	ctx.cancel(status.None)

	ctx1 := Next(ctx)
	ctx1.Done()
}

// NextTimeout

func TestNextTimeout__should_not_deadlock(t *testing.T) {
	ctx := New()
	defer ctx.Free()

	ctx1 := NextTimeout(ctx, time.Second)
	defer ctx1.Free()

	ctx.Cancel()

	assert.True(t, ctx1.Done())
}
