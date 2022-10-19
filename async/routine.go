package async

import (
	"github.com/complex1tech/baselibrary/errors2"
	"github.com/complex1tech/baselibrary/status"
)

var _ CancelWaiter = (Routine[any])(nil)

// Routine is a goroutine wrapper which returns the result as a future, recovers on panics,
// and can be cancelled.
type Routine[T any] interface {
	Future[T]

	// Cancel requests the routine to cancel.
	Cancel() <-chan struct{}

	// CancelWait cancels the routine and awaits its result.
	CancelWait()
}

// Run runs a function in a routine.
func Run(fn func(stop <-chan struct{}) status.Status) Routine[struct{}] {
	r := newRoutine[struct{}]()

	go func() {
		defer func() {
			e := recover()
			if e == nil {
				return
			}

			err := errors2.Recover(e)
			st := status.WrapError(err)
			r.Complete(struct{}{}, st)
		}()

		st := fn(r.cancel)
		r.Complete(struct{}{}, st)
	}()

	return r
}

// Execute executes a function in a routine.
func Execute[T any](fn func(stop <-chan struct{}) (T, status.Status)) Routine[T] {
	r := newRoutine[T]()

	go func() {
		defer func() {
			e := recover()
			if e == nil {
				return
			}

			err := errors2.Recover(e)
			st := status.WrapError(err)

			var zero T
			r.Complete(zero, st)
		}()

		result, st := fn(r.cancel)
		r.Complete(result, st)
	}()

	return r
}

// internal

type routine[T any] struct {
	promise[T]

	cancel    chan struct{}
	cancelled bool
}

func newRoutine[T any]() *routine[T] {
	return &routine[T]{
		promise: promise[T]{
			wait: make(chan struct{}),
		},
		cancel: make(chan struct{}),
	}
}

// Cancel requests the routine to cancel.
func (r *routine[T]) Cancel() <-chan struct{} {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.done || r.cancelled {
		return r.wait
	}

	close(r.cancel)
	r.cancelled = true
	return r.wait
}

// CancelWait cancels the routine and awaits its result.
func (r *routine[T]) CancelWait() {
	r.Cancel()
	<-r.Wait()
}
