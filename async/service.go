// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"sync"

	"github.com/basecomplextech/baselibrary/collect/chans"
	"github.com/basecomplextech/baselibrary/status"
)

// Service is a service which can be started and stopped.
type Service interface {
	// Wait returns a channel which is closed when the service is stopped.
	Wait() <-chan struct{}

	// Status returns the exit status or none.
	Status() status.Status

	// Flags

	// Running indicates that the service is running.
	Running() Flag

	// Stopped indicates that the service is stopped.
	Stopped() Flag

	// Methods

	// Start starts the service if not running and returns its routine.
	Start() (Routine[struct{}], status.Status)

	// Stop requests the service to stop and returns its routine or a stopped routine.
	Stop() <-chan struct{}
}

// NewService returns a new stopped service.
func NewService(fn func(ctx Context) status.Status) Service {
	return newService(fn)
}

// internal

var _ Service = (*service)(nil)

type service struct {
	fn func(ctx Context) status.Status

	running MutFlag
	stopped MutFlag

	mu      sync.Mutex
	routine Routine[struct{}]
}

func newService(fn func(ctx Context) status.Status) *service {
	return &service{
		fn:      fn,
		running: UnsetFlag(),
		stopped: SetFlag(),
	}
}

// Wait returns a channel which is closed when the service is stopped.
func (s *service) Wait() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.routine == nil {
		return chans.Closed()
	}
	return s.routine.Wait()
}

// Status returns the exit status or none.
func (s *service) Status() status.Status {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.routine == nil {
		return status.None
	}
	return s.routine.Status()
}

// Flags

// Running indicates that the service is running.
func (s *service) Running() Flag {
	return s.running
}

// Stopped indicates that the service is stopped.
func (s *service) Stopped() Flag {
	return s.stopped
}

// Methods

// Start starts the service and returns its routine.
func (s *service) Start() (Routine[struct{}], status.Status) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.routine != nil {
		select {
		case <-s.routine.Wait():
		default:
			return s.routine, status.OK
		}
	}

	s.routine = Go(s.run)
	s.running.Set()
	s.stopped.Unset()
	return s.routine, status.OK
}

// Stop requests the service to stop.
func (s *service) Stop() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.routine == nil {
		return chans.Closed()
	}
	return s.routine.Stop()
}

// private

func (s *service) run(ctx Context) status.Status {
	defer s.stop()

	return s.fn(ctx)
}

func (s *service) stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.running.Unset()
	s.stopped.Set()
}

var stoppedRoutine = func() Routine[struct{}] {
	return Exited(struct{}{}, status.Cancelled)
}()
