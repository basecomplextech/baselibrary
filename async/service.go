package async

import (
	"sync"

	"github.com/complex1tech/baselibrary/status"
)

// Service executes a function in a goroutine, recovers on panics, and returns an exit status.
// Service can be restarted multiple times.
type Service interface {
	// Status returns none if running, unavailable if stopped successfully, or an error status.
	Status() status.Status

	// Flags

	// Running returns a channel which is closed when the service is running.
	Running() <-chan struct{}

	// Stopped returns a channel which is closed when the service is stopped.
	Stopped() <-chan struct{}

	// Start/stop

	// Start starts/restarts the service if not already running.
	Start()

	// Stop stops the service and returns a channel which is closed when the service stops.
	Stop() <-chan struct{}
}

// NewService returns a new service.
func NewService(fn func(stop <-chan struct{}) status.Status) Service {
	return newService(fn)
}

// internal

var stopped = status.Unavailable("stopped")

type service struct {
	fn func(stop <-chan struct{}) status.Status

	running *Flag
	stopped *Flag

	mu      sync.Mutex
	status  status.Status
	routine Routine[struct{}]
}

func newService(fn func(stop <-chan struct{}) status.Status) *service {
	return &service{
		fn: fn,

		running: NewFlag(),
		stopped: SetFlag(),

		status: stopped,
	}
}

// Status returns none if running, unavailable if stopped successfully, or an error status.
func (th *service) Status() status.Status {
	th.mu.Lock()
	defer th.mu.Unlock()

	return th.status
}

// Flags

// Running returns a channel which is closed when the service is running.
func (th *service) Running() <-chan struct{} {
	return th.running.Wait()
}

// Stopped returns a channel which is closed when the service is stopped.
func (th *service) Stopped() <-chan struct{} {
	return th.stopped.Wait()
}

// Start/stop

// Start starts/restarts the service if not already running.
func (th *service) Start() {
	th.mu.Lock()
	defer th.mu.Unlock()

	if th.routine != nil {
		return
	}

	th.status = status.None
	th.routine = Run(th.run)

	th.running.Signal()
	th.stopped.Reset()
}

// Stop stops the service and returns a channel which is closed when the service stops.
func (th *service) Stop() <-chan struct{} {
	th.mu.Lock()
	defer th.mu.Unlock()

	if th.routine == nil {
		th.status = stopped
		return closedChan
	}

	return th.routine.Cancel()
}

// private

func (th *service) run(stop <-chan struct{}) (st status.Status) {
	defer func() {
		th.mu.Lock()
		defer th.mu.Unlock()

		if st.OK() {
			th.status = stopped
		} else {
			th.status = st
		}

		th.routine = nil
		th.running.Reset()
		th.stopped.Signal()
	}()

	st = th.fn(stop)
	return
}
