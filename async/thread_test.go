package async

import (
	"testing"
	"time"

	"github.com/complex1tech/baselibrary/panics"
	"github.com/complex1tech/baselibrary/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Run

func TestRun__should_return_on_on_success(t *testing.T) {
	p := Run(func(cancel <-chan struct{}) status.Status {
		return status.OK
	})
	p.Cancel()

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

func TestRun__should_return_status_on_error(t *testing.T) {
	st := status.Error("test")
	p := Run(func(cancel <-chan struct{}) status.Status {
		return st
	})
	p.Cancel()

	select {
	case <-p.Wait():
		_, st1 := p.Result()
		assert.Equal(t, st, st1)

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestRun__should_return_recover_on_panic(t *testing.T) {
	p := Run(func(cancel <-chan struct{}) status.Status {
		panic("test")
	})
	p.Cancel()

	select {
	case <-p.Wait():
		_, st := p.Result()
		require.IsType(t, &panics.Error{}, st.Error)
		assert.EqualValues(t, "test", st.Error.(*panics.Error).E)

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestRun__should_stop_on_request(t *testing.T) {
	p := Run(func(cancel <-chan struct{}) status.Status {
		<-cancel
		return status.Cancelled
	})

	p.Cancel()
	select {
	case <-p.Wait():
		st := p.Status()
		assert.Equal(t, status.Cancelled, st)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

// Execute

func TestExecute__should_return_result_on_success(t *testing.T) {
	p := Execute(func(cancel <-chan struct{}) (string, status.Status) {
		return "hello, world", status.OK
	})
	p.Cancel()

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

func TestExecute__should_return_status_on_error(t *testing.T) {
	st := status.Error("test")
	p := Execute(func(cancel <-chan struct{}) (string, status.Status) {
		return "", st
	})
	p.Cancel()

	select {
	case <-p.Wait():
		_, st1 := p.Result()
		assert.Equal(t, st, st1)

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestExecute__should_return_recover_on_panic(t *testing.T) {
	p := Execute(func(cancel <-chan struct{}) (string, status.Status) {
		panic("test")
	})
	p.Cancel()

	select {
	case <-p.Wait():
		_, st := p.Result()
		require.IsType(t, &panics.Error{}, st.Error)
		assert.EqualValues(t, "test", st.Error.(*panics.Error).E)

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestExecute__should_stop_on_request(t *testing.T) {
	p := Execute(func(cancel <-chan struct{}) (string, status.Status) {
		<-cancel
		return "", status.Cancelled
	})

	select {
	case <-p.Cancel():
		st := p.Status()
		assert.Equal(t, status.Cancelled, st)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}
