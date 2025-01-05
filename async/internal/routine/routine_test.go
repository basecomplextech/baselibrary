// Copyright 2021 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package routine

import (
	"testing"
	"time"

	"github.com/basecomplextech/baselibrary/async/internal/context"
	"github.com/basecomplextech/baselibrary/panics"
	"github.com/basecomplextech/baselibrary/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Go

func TestGo__should_return_on_on_success(t *testing.T) {
	r := Go(func(context.Context) status.Status {
		return status.OK
	})
	r.Stop()

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

func TestGo__should_return_status_on_error(t *testing.T) {
	st := status.Test("test")
	r := Go(func(context.Context) status.Status {
		return st
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

func TestGo__should_return_recover_on_panic(t *testing.T) {
	r := Go(func(context.Context) status.Status {
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

func TestGo__should_stop_on_request(t *testing.T) {
	r := Go(func(ctx context.Context) status.Status {
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

func TestRun__should_stop_on_request(t *testing.T) {
	r := Run(func(ctx context.Context) (string, status.Status) {
		<-ctx.Wait()
		return "", ctx.Status()
	})

	select {
	case <-r.Stop():
		st := r.Status()
		assert.Equal(t, status.Cancelled, st)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}
