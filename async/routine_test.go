package async

import (
	"testing"
	"time"

	"github.com/epochtimeout/baselibrary/status"
	"github.com/epochtimeout/baselibrary/try"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun__should_run_and_stop(t *testing.T) {
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

func TestRun__should_run_and_return_error(t *testing.T) {
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

func TestRun__should_run_and_recover(t *testing.T) {
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
