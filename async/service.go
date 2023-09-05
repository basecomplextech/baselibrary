package async

import (
	"sync"

	"github.com/basecomplextech/baselibrary/status"
)

// Service is a service which can be started and stopped.
type Service interface {
	// Running indicates that the service is running.
	Running() <-chan struct{}

	// Stopped indicates that the service is stopped.
	Stopped() <-chan struct{}

	// Routine returns the service routine or a stopped routine if the service is not running.
	Routine() Routine[struct{}]

	// Methods

	// Start starts the service if not running and returns its routine.
	Start() (Routine[struct{}], status.Status)

	// Stop requests the service to stop.
	Stop()
}

// NewService returns a new stopped service.
func NewService(fn func(cancel <-chan struct{}) status.Status) Service {
	return newService(fn)
}

// internal

var _ Service = (*service)(nil)

type service struct {
	fn func(cancel <-chan struct{}) status.Status

	running *Flag
	stopped *Flag

	mu      sync.Mutex
	routine Routine[struct{}]
}

func newService(fn func(cancel <-chan struct{}) status.Status) *service {
	return &service{
		fn:      fn,
		running: UnsetFlag(),
		stopped: SetFlag(),
	}
}

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
func (s *service) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.routine == nil {
		return
	}

	s.routine.Cancel()
}

// private

func (s *service) run(cancel <-chan struct{}) status.Status {
	defer s.stop()

	return s.fn(cancel)
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
