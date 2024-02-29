package async

import (
	"sync"

	"github.com/basecomplextech/baselibrary/status"
)

// Service is a service which can be started and stopped.
type Service interface {
	// IsRunning returns true if the service is running.
	IsRunning() bool

	// Async

	// Running indicates that the service is running.
	Running() <-chan struct{}

	// Stopped indicates that the service is stopped.
	Stopped() <-chan struct{}

	// Methods

	// Start starts the service if not running and returns its routine.
	Start() (Routine[struct{}], status.Status)

	// Stop requests the service to stop and returns its routine or a stopped routine.
	Stop() Routine[struct{}]

	// Routine returns the service routine or a stopped routine if the service is not running.
	Routine() Routine[struct{}]
}

// NewService returns a new stopped service.
func NewService(fn func(ctx Context) status.Status) Service {
	return newService(fn)
}

// internal

var _ Service = (*service)(nil)

type service struct {
	fn func(ctx Context) status.Status

	running *Flag
	stopped *Flag

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

// IsRunning returns true if the service is running.
func (s *service) IsRunning() bool {
	return s.running.IsSet()
}

// Async

// Running indicates that the service is running.
func (s *service) Running() <-chan struct{} {
	return s.running.Wait()
}

// Stopped indicates that the service is stopped.
func (s *service) Stopped() <-chan struct{} {
	return s.stopped.Wait()
}

// Start starts the service if not running and returns its routine.
func (s *service) Routine() Routine[struct{}] {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.routine == nil {
		return stoppedRoutine
	}
	return s.routine
}

// Methods

// Start starts the service and returns its routine.
func (s *service) Start() (Routine[struct{}], status.Status) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.routine != nil && s.running.IsSet() {
		return s.routine, status.OK
	}

	s.routine = Go(s.run)
	s.running.Set()
	s.stopped.Unset()
	return s.routine, status.OK
}

// Stop requests the service to stop.
func (s *service) Stop() Routine[struct{}] {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.routine == nil {
		return stoppedRoutine
	}

	r := s.routine
	r.Cancel()
	return r
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
	s.routine = nil
}

var stoppedRoutine = func() Routine[struct{}] {
	return Exited(struct{}{}, status.Cancelled)
}()
