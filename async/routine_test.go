// Copyright 2021 Ivan Korobkov. All rights reserved.

package async

import (
	"testing"
	"time"

	"github.com/basecomplextech/baselibrary/panics"
	"github.com/basecomplextech/baselibrary/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Go

func TestGo__should_return_on_on_success(t *testing.T) {
	p := Go(func(Context) status.Status {
		return status.OK
	})
	p.Stop()

	select {
	case <-p.Wait():
		_, st := p.Result()
		if !st.OK() {
			t.Fatal(st)
		}

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestGo__should_return_status_on_error(t *testing.T) {
	st := status.Test("test")
	p := Go(func(Context) status.Status {
		return st
	})
	p.Stop()

	select {
	case <-p.Wait():
		_, st1 := p.Result()
		assert.Equal(t, st, st1)

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestGo__should_return_recover_on_panic(t *testing.T) {
	p := Go(func(Context) status.Status {
		panic("test")
	})
	p.Stop()

	select {
	case <-p.Wait():
		_, st := p.Result()
		require.IsType(t, &panics.Error{}, st.Error)
		assert.EqualValues(t, "test", st.Error.(*panics.Error).E)

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestGo__should_stop_on_request(t *testing.T) {
	p := Go(func(ctx Context) status.Status {
		<-ctx.Wait()
		return ctx.Status()
	})

	p.Stop()
	select {
	case <-p.Wait():
		st := p.Status()
		assert.Equal(t, status.Cancelled, st)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

// Call

func TestCall__should_return_result_on_success(t *testing.T) {
	p := Call(func(Context) (string, status.Status) {
		return "hello, world", status.OK
	})
	p.Stop()

	select {
	case <-p.Wait():
		v, st := p.Result()
		if !st.OK() {
			t.Fatal(st)
		}
		assert.Equal(t, "hello, world", v)

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestCall__should_return_status_on_error(t *testing.T) {
	st := status.Test("test")
	p := Call(func(Context) (string, status.Status) {
		return "", st
	})
	p.Stop()

	select {
	case <-p.Wait():
		_, st1 := p.Result()
		assert.Equal(t, st, st1)

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestCall__should_return_recover_on_panic(t *testing.T) {
	p := Call(func(Context) (string, status.Status) {
		panic("test")
	})
	p.Stop()

	select {
	case <-p.Wait():
		_, st := p.Result()
		require.IsType(t, &panics.Error{}, st.Error)
		assert.EqualValues(t, "test", st.Error.(*panics.Error).E)

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestCall__should_stop_on_request(t *testing.T) {
	p := Call(func(ctx Context) (string, status.Status) {
		<-ctx.Wait()
		return "", ctx.Status()
	})

	select {
	case <-p.Stop():
		st := p.Status()
		assert.Equal(t, status.Cancelled, st)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}
