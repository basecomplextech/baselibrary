package async

import (
	"github.com/complex1tech/baselibrary/status"
)

// Routine is a generic async routine which returns the result as a future, and can be cancelled.
type Routine[T any] interface {
	Future[T]

	// Cancel requests the routine to cancel and returns a wait channel.
	Cancel() <-chan struct{}
}

// Methods

// Run runs a function in a new routine, recovers on panics.
func Run(fn func(cancel <-chan struct{}) status.Status) Routine[struct{}] {
	r := newRoutine[struct{}]()

	go func() {
		defer func() {
			e := recover()
			if e == nil {
				return
			}

			st := status.Recover(e)
			r.Complete(struct{}{}, st)
		}()

		st := fn(r.cancel)
		r.Complete(struct{}{}, st)
	}()

	return r
}

// Execute executes a function in a new routine, recovers on panics.
func Execute[T any](fn func(cancel <-chan struct{}) (T, status.Status)) Routine[T] {
	r := newRoutine[T]()

	go func() {
		defer func() {
			e := recover()
			if e == nil {
				return
			}

			var zero T
			st := status.Recover(e)
			r.Complete(zero, st)
		}()

		result, st := fn(r.cancel)
		r.Complete(result, st)
	}()

	return r
}

// Join joins all routines into a single routine.
// The routine returns all the results and the first non-OK status.
func Join[T any](routines ...Routine[T]) Routine[[]T] {
	return Execute(func(cancel <-chan struct{}) ([]T, status.Status) {
		// Await all or cancel
		st := status.OK
	loop:
		for _, r := range routines {
			select {
			case <-r.Wait():
			case <-cancel:
				st = status.Cancelled
				break loop
			}
		}

		// Cancel all
		for _, r := range routines {
			r.Cancel()
		}

		// Collect results
		results := make([]T, 0, len(routines))
		for _, r := range routines {
			<-r.Wait()

			r, st1 := r.Result()
			if !st1.OK() && st.OK() {
				st = st1
			}

			results = append(results, r)
		}
		return results, st
	})
}

// internal

var (
	_ Routine[any] = (*routine[any])(nil)
	_ CancelWaiter = (*routine[any])(nil)
)

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

// Cancel requests the routine to cancel and returns a wait channel.
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
