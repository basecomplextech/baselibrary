package async

import (
	"github.com/basecomplextech/baselibrary/status"
)

// Routine is an async routine which returns the result as a future, recovers on panics,
// and can be cancelled.
type Routine[T any] interface {
	Future[T]
}

// Methods

// Go runs a function in a new routine, recovers on panics.
func Go(fn func(cancel <-chan struct{}) status.Status) Routine[struct{}] {
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

// Call calls a function in a new routine, and returns its result, recovers on panics.
func Call[T any](fn func(cancel <-chan struct{}) (T, status.Status)) Routine[T] {
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
	return Call(func(cancel <-chan struct{}) ([]T, status.Status) {
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

// Exited returns a routine which has exited with the given result and status.
func Exited[T any](result T, st status.Status) Routine[T] {
	r := newRoutine[T]()
	r.Complete(result, st)
	return r
}

// internal

var _ Routine[any] = (*routine[any])(nil)

type routine[T any] struct {
	promise[T]
}

func newRoutine[T any]() *routine[T] {
	return &routine[T]{
		promise: promise[T]{
			wait:   make(chan struct{}),
			cancel: make(chan struct{}),
		},
	}
}
