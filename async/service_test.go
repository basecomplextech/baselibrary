// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"testing"
	"time"

	"github.com/basecomplextech/baselibrary/async/internal/context"
	"github.com/basecomplextech/baselibrary/status"
	"github.com/stretchr/testify/assert"
)

// Start

func TestService_Start__should_set_running_after_start(t *testing.T) {
	s := NewService(func(ctx context.Context) status.Status {
		<-ctx.Wait()
		return status.OK
	})

	s.Start()
	defer s.Stop()

	select {
	case <-s.Running().Wait():
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

func TestService_Start__should_restart_service(t *testing.T) {
	done := make(chan struct{}, 1)
	s := NewService(func(ctx context.Context) status.Status {
		done <- struct{}{}
		<-ctx.Wait()
		return status.OK
	})
	defer s.Stop()

	s.Start()
	<-done
	<-s.Stop()

	s.Start()
	<-done
}

// Stop

func TestService_Stop__should_stop_service(t *testing.T) {
	s := NewService(func(ctx context.Context) status.Status {
		<-ctx.Wait()
		return status.OK
	})

	s.Start()
	s.Stop()

	select {
	case <-s.Wait():
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

func TestService_Stop__should_set_stopped_after_status(t *testing.T) {
	s := NewService(func(ctx context.Context) status.Status {
		return status.Error("test")
	})

	for i := 0; i < 100; i++ {
		s.Start()
		<-s.Stop()

		st := s.Status()
		assert.Equal(t, status.Error("test"), st)
	}
}
