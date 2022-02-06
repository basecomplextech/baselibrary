package async

import (
	"errors"
	"sync"

	"github.com/baseone-run/library/try"
)

// Stopped is the error which can be used to indicate that a routine has been stopped.
var Stopped = errors.New("routine stopped")

// Routine runs a goroutine and returns its result as a future,
// recovers on a panic and returns a *try.Panic error with a stack trace.
type Routine[T any] interface {
	Future[T]

	// Stop requests the routine to stop and returns its wait channel.
	Stop() <-chan struct{}
}

// Run runs a function in a routine.
func Run(fn func(stop <-chan struct{}) error) Routine[Void] {
	r := newRoutine[Void]()

	go func() {
		defer func() {
			if e := recover(); e != nil {
				err := try.Recover(e)
				r.Reject(err)
			}
		}()

		err := fn(r.stopCh)
		r.Reject(err)
	}()

	return r
}

// Call calls a function in a routine and returns its result.
func Call[T any](fn func(stop <-chan struct{}) (T, error)) Routine[T] {
	r := newRoutine[T]()

	go func() {
		defer func() {
			if e := recover(); e != nil {
				err := try.Recover(e)
				r.Reject(err)
			}
		}()

		result, err := fn(r.stopCh)
		r.Complete(result, err)
	}()

	return r
}

type routine[T any] struct {
	Promise[T]

	mu     sync.Mutex
	stop   bool
	stopCh chan struct{}
}

func newRoutine[T any]() *routine[T] {
	return &routine[T]{
		Promise: Pending[T](),
		stopCh:  make(chan struct{}),
	}
}

// Stop requests the routine to stop and returns its wait channel.
func (r *routine[T]) Stop() <-chan struct{} {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.stop {
		return r.Wait()
	}

	r.stop = true
	close(r.stopCh)
	return r.Wait()
}
