package async

import (
	"testing"
	"time"

	"github.com/epochtimeout/baselibrary/status"
	"github.com/stretchr/testify/assert"
)

func TestThread_Start__should_run_function(t *testing.T) {
	fn := func(stop <-chan struct{}) status.Status {
		<-stop
		return status.OK
	}

	th := NewVoidThread(fn)
	th.Start()
	defer th.Stop()

	select {
	case <-th.Running():
	default:
		t.Error("thread should be running")
	}

	select {
	case <-th.Stopped():
		t.Error("thread should be running")
	default:
	}
}

func TestThread_Start__should_return_when_already_running(t *testing.T) {
	fn := func(stop <-chan struct{}) status.Status {
		<-stop
		return status.OK
	}

	th := NewVoidThread(fn)
	th.Start()
	defer th.Stop()

	th.Start()

	select {
	case <-th.Running():
	default:
		t.Error("thread should be running")
	}
}

func TestThread_Start__should_restart_stopped_thread(t *testing.T) {
	var i int
	fn := func(stop <-chan struct{}) (int, status.Status) {
		i++
		return i, status.OK
	}

	th := NewThread(fn)
	th.Start()
	defer th.Stop()

	select {
	case <-th.Wait():
	case <-time.After(time.Second):
		t.Error("result timeout")
	}

	result, st := th.Result()
	assert.Equal(t, 1, result)
	assert.Equal(t, status.CodeOK, st.Code)

	th.Start()
	select {
	case <-th.Wait():
	case <-time.After(time.Second):
		t.Error("result timeout")
	}

	result, st = th.Result()
	assert.Equal(t, 2, result)
	assert.Equal(t, status.CodeOK, st.Code)
}

// Stop

func TestThread_Stop__should_stop_thread(t *testing.T) {
	fn := func(stop <-chan struct{}) status.Status {
		<-stop
		return status.Cancelled
	}

	th := NewVoidThread(fn)
	th.Start()
	defer th.Stop()

	select {
	case <-th.Stop():
	case <-time.After(time.Second):
		t.Error("thread should be stopped")
	}

	select {
	case <-th.Stopped():
	case <-time.After(time.Second):
		t.Error("thread should be stopped")
	}

	select {
	case <-th.Running():
		t.Error("thread should not be running")
	default:
	}

	_, st := th.Result()
	assert.Equal(t, status.CodeCancelled, st.Code)
}

func TestThread_Stop__should_cancel_result_when_thread_not_started(t *testing.T) {
	fn := func(stop <-chan struct{}) status.Status {
		<-stop
		return status.OK
	}

	th := NewVoidThread(fn)
	go func() {
		time.Sleep(50 * time.Millisecond)
		th.Stop()
	}()

	select {
	case <-th.Wait():
	case <-time.After(time.Second):
		t.Error("result timeout")
	}

	_, st := th.Result()
	assert.Equal(t, status.CodeCancelled, st.Code)
}

// Wait

func TestThread_Wait__should_await_result_even_before_thread_is_started(t *testing.T) {
	var i int
	fn := func(stop <-chan struct{}) (int, status.Status) {
		i++
		return i, status.OK
	}

	th := NewThread(fn)
	go func() {
		time.Sleep(50 * time.Millisecond)
		th.Start()
	}()

	select {
	case <-th.Wait():
	case <-time.After(time.Second):
		t.Error("result timeout")
	}

	i, st := th.Result()
	assert.Equal(t, 1, i)
	assert.Equal(t, status.CodeOK, st.Code)
}
