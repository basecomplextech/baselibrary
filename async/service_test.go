package async

import (
	"testing"
	"time"

	"github.com/complex1tech/baselibrary/status"
	"github.com/stretchr/testify/assert"
)

func TestService_Start__should_run_function(t *testing.T) {
	fn := func(stop <-chan struct{}) status.Status {
		<-stop
		return status.OK
	}

	s := NewService(fn)
	s.Start()
	defer s.Stop()

	select {
	case <-s.Running():
	default:
		t.Error("service should be running")
	}

	select {
	case <-s.Stopped():
		t.Error("service should be running")
	default:
	}
}

func TestService_Start__should_return_when_already_running(t *testing.T) {
	fn := func(stop <-chan struct{}) status.Status {
		<-stop
		return status.OK
	}

	s := NewService(fn)
	s.Start()
	defer s.Stop()

	s.Start()

	select {
	case <-s.Running():
	default:
		t.Error("service should be running")
	}
}

func TestService_Start__should_restart_stopped_service(t *testing.T) {
	var i int
	fn := func(stop <-chan struct{}) status.Status {
		i++
		return status.OK
	}

	s := NewService(fn)
	s.Start()
	defer s.Stop()

	select {
	case <-s.Stopped():
	case <-time.After(time.Second):
		t.Error("result timeout")
	}

	st := s.Status()
	assert.Equal(t, 1, i)
	assert.Equal(t, status.CodeUnavailable, st.Code)

	s.Start()
	select {
	case <-s.Stopped():
	case <-time.After(time.Second):
		t.Error("result timeout")
	}

	st = s.Status()
	assert.Equal(t, 2, i)
	assert.Equal(t, status.CodeUnavailable, st.Code)
}

// Stop

func TestService_Stop__should_stop_service(t *testing.T) {
	fn := func(stop <-chan struct{}) status.Status {
		<-stop
		return status.OK
	}

	s := NewService(fn)
	s.Start()
	defer s.Stop()

	select {
	case <-s.Stop():
	case <-time.After(time.Second):
		t.Error("service should be stopped")
	}

	select {
	case <-s.Stopped():
	case <-time.After(time.Second):
		t.Error("service should be stopped")
	}

	select {
	case <-s.Running():
		t.Error("service should not be running")
	default:
	}

	st := s.Status()
	assert.Equal(t, status.CodeUnavailable, st.Code)
}

func TestService_Stop__should_stop_service_when_service_not_started(t *testing.T) {
	fn := func(stop <-chan struct{}) status.Status {
		<-stop
		return status.OK
	}

	s := NewService(fn)
	s.Stop()

	select {
	case <-s.Stopped():
	case <-time.After(time.Second):
		t.Error("result timeout")
	}

	st := s.Status()
	assert.Equal(t, status.CodeUnavailable, st.Code)
}
