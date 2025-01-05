// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package service

import (
	"sync"

	"github.com/basecomplextech/baselibrary/async/internal/context"
	"github.com/basecomplextech/baselibrary/async/internal/flag"
	"github.com/basecomplextech/baselibrary/async/internal/routine"
	"github.com/basecomplextech/baselibrary/collect/chans"
	"github.com/basecomplextech/baselibrary/opt"
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
	Running() flag.Flag

	// Stopped indicates that the service is stopped.
	Stopped() flag.Flag

	// Methods

	// Start starts the service if not running and returns its routine.
	Start() (routine.RoutineVoid, status.Status)

	// Stop requests the service to stop and returns its routine or a stopped routine.
	Stop() <-chan struct{}
}

// New returns a new stopped service.
func New(fn func(ctx context.Context) status.Status) Service {
	return newService(fn)
}

// internal

var _ Service = (*service)(nil)

type service struct {
	fn func(ctx context.Context) status.Status

	running flag.MutFlag
	stopped flag.MutFlag

	mu      sync.Mutex
	routine opt.Opt[routine.RoutineVoid]
}

func newService(fn func(ctx context.Context) status.Status) *service {
	return &service{
		fn:      fn,
		running: flag.UnsetFlag(),
		stopped: flag.SetFlag(),
	}
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

// Status returns the exit status or none.
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
func (s *service) Running() flag.Flag {
	return s.running
}

// Stopped indicates that the service is stopped.
func (s *service) Stopped() flag.Flag {
	return s.stopped
}

// Methods

// Start starts the service and returns its routine.
func (s *service) Start() (routine.RoutineVoid, status.Status) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.routine.Unwrap()
	if ok {
		if !r.Done() {
			return r, status.OK
		}
	}

	s.running.Set()
	s.stopped.Unset()

	r = routine.Go(s.run)
	s.routine.Set(r)
	return r, status.OK
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

// private

func (s *service) run(ctx context.Context) status.Status {
	defer s.running.Unset()
	defer s.stopped.Set()

	return s.fn(ctx)
}
