package async

import (
	"github.com/epochtimeout/baselibrary/errors2"
	"github.com/epochtimeout/baselibrary/status"
)

// Routine is a future alias which indicates that this is an async computation in a gorooutine.
// Commented out until Go supports generic type aliases.
// type Routine[T] = Future[T]

// Run runs a function in a goroutine, recovers on panics, and returns its status as a future.
func Run(fn func(stop <-chan struct{}) status.Status) Future[struct{}] {
	r := newRoutine[struct{}]()

	go func() {
		defer func() {
			e := recover()
			if e == nil {
				return
			}

			err := errors2.Recover(e)
			st := status.WrapError(err)
			r.reject(st)
		}()

		st := fn(r.stop)
		r.complete(struct{}{}, st)
	}()

	return r
}

// Execute executes a function in a goroutine, recovers on panics, and returns its result as a future.
func Execute[T any](fn func(stop <-chan struct{}) (T, status.Status)) Future[T] {
	r := newRoutine[T]()

	go func() {
		defer func() {
			e := recover()
			if e == nil {
				return
			}

			err := errors2.Recover(e)
			st := status.WrapError(err)
			r.reject(st)
		}()

		result, st := fn(r.stop)
		r.complete(result, st)
	}()

	return r
}

// internal

type routine[T any] struct {
	result[T]

	stop  chan struct{}
	stop_ bool
}

func newRoutine[T any]() *routine[T] {
	return &routine[T]{
		result: newResult[T](),
		stop:   make(chan struct{}),
	}
}

// Cancel tries to cancel the routine and returns the wait channel.
func (r *routine[T]) Cancel() <-chan struct{} {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.stop_ {
		return r.wait
	}

	r.stop_ = true
	close(r.stop)
	return r.wait
}

// private

func (r *routine[T]) reject(st status.Status) {
	var result T
	r.complete(result, st)
}
