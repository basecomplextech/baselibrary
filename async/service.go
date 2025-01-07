// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"sync"

	"github.com/basecomplextech/baselibrary/collect/chans"
	"github.com/basecomplextech/baselibrary/opt"
	"github.com/basecomplextech/baselibrary/status"
)

// Service is a service which can be started and stopped.
type Service interface {
	// Status returns the stop status or none.
	Status() status.Status

	// Flags

	// Running indicates that the service routine is running.
	Running() Flag

	// Stopped indicates that the service is stopped.
	Stopped() Flag

	// Methods

	// Start starts the service if not running.
	Start() status.Status

	// Stop requests the service to stop and returns its routine or a stopped routine.
	Stop() <-chan struct{}

	// Wait returns a channel which is closed when the service is stopped.
	Wait() <-chan struct{}
}

// NewService returns a new stopped service.
func NewService(fn func(ctx Context) status.Status) Service {
	return newService(fn)
}

// internal

var _ Service = (*service)(nil)

type service struct {
	fn func(ctx Context) status.Status

	// flags
	running MutFlag
	stopped MutFlag

	// routine
	mu      sync.Mutex
	routine opt.Opt[RoutineVoid]
}

func newService(fn func(ctx Context) status.Status) *service {
	return &service{
		fn:      fn,
		running: UnsetFlag(),
		stopped: SetFlag(),
	}
}

// Status returns the stop status or none.
func (s *service) Status() status.Status {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.routine.Unwrap()
	if !ok {
		return status.None
	}
	return r.Status()
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

// Start starts the service if not running.
func (s *service) Start() status.Status {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Return if running
	r, ok := s.routine.Unwrap()
	if ok {
		if !r.Done() {
			return status.OK
		}
	}

	// Reset stopped
	s.stopped.Unset()

	// Make routine
	r = NewRoutineVoid(s.run)
	r.OnStop(s.onStop)
	s.routine.Set(r)

	// Start routine
	r.Start()
	return status.OK
}

// Stop requests the service to stop.
func (s *service) Stop() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.routine.Unwrap()
	if !ok {
		return chans.Closed()
	}
	return r.Stop()
}

// Wait returns a channel which is closed when the service is stopped.
func (s *service) Wait() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.routine.Unwrap()
	if !ok {
		return chans.Closed()
	}
	return r.Wait()
}

// private

func (s *service) run(ctx Context) status.Status {
	defer s.running.Unset()
	s.running.Set()

	return s.fn(ctx)
}

func (s *service) onStop(r RoutineVoid) {
	s.stopped.Set()
}
