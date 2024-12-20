// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"testing"
	"time"

	"github.com/basecomplextech/baselibrary/status"
	"github.com/stretchr/testify/assert"
)

func TestThread_Start__should_start_thread(t *testing.T) {
	done := UnsetFlag()

	th := NewThreadDyn(func(ctx Context) status.Status {
		done.Set()
		return status.OK
	})
	th.Start()

	select {
	case <-done.Wait():
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

func TestThread_Start__should_stop_start_thread_when_stopped(t *testing.T) {
	th := NewThreadDyn(func(ctx Context) status.Status {
		return status.OK
	})
	th.Stop()
	th.Start()

	st := th.Status()
	assert.Equal(t, status.Cancelled, st)
}
