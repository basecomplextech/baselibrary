// Copyright 2021 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"testing"
	"time"

	"github.com/basecomplextech/baselibrary/async/internal/context"
	"github.com/basecomplextech/baselibrary/async/internal/flag"
	"github.com/basecomplextech/baselibrary/panics"
	"github.com/basecomplextech/baselibrary/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Run

func TestRun__should_return_result_on_success(t *testing.T) {
	r := Run(func(context.Context) (string, status.Status) {
		return "hello, world", status.OK
	})
	r.Stop()

	select {
	case <-r.Wait():
		v, st := r.Result()
		if !st.OK() {
			t.Fatal(st)
		}
		assert.Equal(t, "hello, world", v)

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestRun__should_return_status_on_error(t *testing.T) {
	st := status.Test("test")
	r := Run(func(context.Context) (string, status.Status) {
		return "", st
	})
	r.Stop()

	select {
	case <-r.Wait():
		_, st1 := r.Result()
		assert.Equal(t, st, st1)

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestRun__should_return_recover_on_panic(t *testing.T) {
	r := Run(func(context.Context) (string, status.Status) {
		panic("test")
	})
	r.Stop()

	select {
	case <-r.Wait():
		_, st := r.Result()
		require.IsType(t, &panics.Error{}, st.Error)
		assert.EqualValues(t, "test", st.Error.(*panics.Error).E)

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

// RunVoid

func TestRunVoid__should_return_on_success(t *testing.T) {
	r := RunVoid(func(context.Context) status.Status {
		return status.OK
	})
	defer r.Stop()

	select {
	case <-r.Wait():
		_, st := r.Result()
		if !st.OK() {
			t.Fatal(st)
		}

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestRunVoid__should_return_status_on_error(t *testing.T) {
	st := status.Test("test")
	r := RunVoid(func(context.Context) status.Status {
		return st
	})
	defer r.Stop()

	select {
	case <-r.Wait():
		_, st1 := r.Result()
		assert.Equal(t, st, st1)

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestRunVoid__should_return_recover_on_panic(t *testing.T) {
	r := RunVoid(func(context.Context) status.Status {
		panic("test")
	})
	defer r.Stop()

	select {
	case <-r.Wait():
		_, st := r.Result()
		require.IsType(t, &panics.Error{}, st.Error)
		assert.EqualValues(t, "test", st.Error.(*panics.Error).E)

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

// Start

func TestRoutine_Start__should_start_routine(t *testing.T) {
	done := flag.UnsetFlag()

	r := NewRoutineVoid(func(ctx context.Context) status.Status {
		done.Set()
		return status.OK
	})
	r.Start()

	select {
	case <-done.Wait():
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

func TestRoutine_Start__should_not_start_when_already_stopped(t *testing.T) {
	r := NewRoutineVoid(func(ctx context.Context) status.Status {
		return status.OK
	})
	r.Stop()
	r.Start()

	st := r.Status()
	assert.Equal(t, status.Cancelled, st)
}

// Stop

func TestRoutine_Stop__should_request_stop_cancel_context(t *testing.T) {
	r := RunVoid(func(ctx context.Context) status.Status {
		<-ctx.Wait()
		return ctx.Status()
	})
	r.Stop()

	select {
	case <-r.Wait():
		st := r.Status()
		assert.Equal(t, status.Cancelled, st)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestRoutine_Stop__should_reject_not_started_routine(t *testing.T) {
	r := NewRoutineVoid(func(context.Context) status.Status {
		return status.OK
	})
	r.Stop()

	select {
	case <-r.Wait():
		st := r.Status()
		assert.Equal(t, status.Cancelled, st)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

// OnStop

func TestRoutine__should_call_stop_callbacks_on_stop(t *testing.T) {
	done0 := flag.UnsetFlag()
	done1 := flag.UnsetFlag()

	r := NewRoutineVoid(func(ctx context.Context) status.Status {
		return status.OK
	})

	r.OnStop(func(RoutineVoid) {
		done0.Set()
	})
	r.OnStop(func(RoutineVoid) {
		done1.Set()
	})

	r.Start()

	select {
	case <-done0.Wait():
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}

	select {
	case <-done1.Wait():
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}
