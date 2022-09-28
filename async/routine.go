package async

import (
	"github.com/complex1tech/baselibrary/errors2"
	"github.com/complex1tech/baselibrary/status"
)

// Routine is a future alias which indicates that this is an async computation in a gorooutine.
// Commented out until Go supports generic type aliases.
// type Routine[T] = Future[T]

// Run runs a function in a goroutine, recovers on panics, and returns its status as a future.
func Run(fn func(stop <-chan struct{}) status.Status) Future[struct{}] {
	p := newPromise[struct{}]()

	go func() {
		defer func() {
			e := recover()
			if e == nil {
				return
			}

			err := errors2.Recover(e)
			st := status.WrapError(err)
			p.Complete(struct{}{}, st)
		}()

		st := fn(p.cancel)
		p.Complete(struct{}{}, st)
	}()

	return p
}

// Execute executes a function in a goroutine, recovers on panics, and returns its result as a future.
func Execute[T any](fn func(stop <-chan struct{}) (T, status.Status)) Future[T] {
	p := newPromise[T]()

	go func() {
		defer func() {
			e := recover()
			if e == nil {
				return
			}

			err := errors2.Recover(e)
			st := status.WrapError(err)

			var zero T
			p.Complete(zero, st)
		}()

		result, st := fn(p.cancel)
		p.Complete(result, st)
	}()

	return p
}
