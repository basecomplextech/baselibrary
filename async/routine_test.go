package async

import (
	"testing"
	"time"

	"github.com/epochtimeout/baselibrary/errors2"
	"github.com/epochtimeout/baselibrary/status"
	"github.com/epochtimeout/baselibrary/try"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Run

func TestRun__should_return_on_on_success(t *testing.T) {
	r := Run(func(stop <-chan struct{}) status.Status {
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

func TestRun__should_return_status_on_error(t *testing.T) {
	st := status.Error("test")
	r := Run(func(stop <-chan struct{}) status.Status {
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

func TestRun__should_return_recover_on_panic(t *testing.T) {
	r := Run(func(stop <-chan struct{}) status.Status {
		panic("test")
	})
	r.Stop()

	select {
	case <-r.Wait():
		_, st := r.Result()
		require.IsType(t, &try.Panic{}, st.Error)
		assert.EqualValues(t, "test", st.Error.(*try.Panic).E)

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestRun__should_stop_on_request(t *testing.T) {
	r := Run(func(stop <-chan struct{}) status.Status {
		<-stop
		return status.Cancelled
	})

	select {
	case <-r.Stop():
		st := r.Status()
		assert.Equal(t, status.Cancelled, st)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

// Call

func TestCall__should_return_result_on_success(t *testing.T) {
	r := Call(func(stop <-chan struct{}) (string, status.Status) {
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

func TestCall__should_return_status_on_error(t *testing.T) {
	st := status.Error("test")
	r := Call(func(stop <-chan struct{}) (string, status.Status) {
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

func TestCall__should_return_recover_on_panic(t *testing.T) {
	r := Call(func(stop <-chan struct{}) (string, status.Status) {
		panic("test")
	})
	r.Stop()

	select {
	case <-r.Wait():
		_, st := r.Result()
		require.IsType(t, &errors2.PanicError{}, st.Error)
		assert.EqualValues(t, "test", st.Error.(*errors2.PanicError).E)

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestCall__should_stop_on_request(t *testing.T) {
	r := Call(func(stop <-chan struct{}) (string, status.Status) {
		<-stop
		return "", status.Cancelled
	})

	select {
	case <-r.Stop():
		st := r.Status()
		assert.Equal(t, status.Cancelled, st)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}
