package async

import (
	"sync"

	"github.com/epochtimeout/baselibrary/errors2"
	"github.com/epochtimeout/baselibrary/status"
)

// Routine executes a function in a goroutine, recovers on panics, and returns a future result.
type Routine[T any] interface {
	Result[T]

	// Stop requests the routine to stop and returns the wait channel.
	Stop() <-chan struct{}
}

// Run runs a function in a routine and returns its status, recovers on panics.
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
			r.reject(st)
		}()

		st := fn(r.stop)
		r.complete(struct{}{}, st)
	}()

	return r
}

// Execute executes a function in a routine and returns its result and status, recovers on panics.
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
			r.reject(st)
		}()

		result, st := fn(r.stop)
		r.complete(result, st)
	}()

	return r
}

// ToFuture returns a routine as a future.
func ToFuture[T any](r Routine[T]) Future[T] {
	return newRoutineFuture(r)
}

// StopAll stops all routines.
func StopAll[T any](routines ...Routine[T]) {
	for _, r := range routines {
		r.Stop()
	}
}

// StopWait stops a routine and awaits its result.
func StopWait[T any](r Routine[T]) (T, status.Status) {
	<-r.Stop()
	return r.Result()
}

// internal

type routine[T any] struct {
	mu sync.Mutex

	done   bool
	result T
	status status.Status

	wait  chan struct{}
	stop  chan struct{}
	stop_ bool
}

func newRoutine[T any]() *routine[T] {
	return &routine[T]{
		wait: make(chan struct{}),
		stop: make(chan struct{}),
	}
}

// Wait awaits the result.
func (r *routine[T]) Wait() <-chan struct{} {
	return r.wait
}

// Result returns a value and a status or zero.
func (r *routine[T]) Result() (T, status.Status) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.result, r.status
}

// Status returns a status.
func (r *routine[T]) Status() status.Status {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.status
}

// Stop requests the routine to stop and returns its wait channel.
func (r *routine[T]) Stop() <-chan struct{} {
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

func (r *routine[T]) complete(result T, st status.Status) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.done {
		return
	}

	r.done = true
	r.result = result
	r.status = st
	close(r.wait)
}

// routine future

type routineFuture[T any] struct {
	Routine[T]
}

func newRoutineFuture[T any](r Routine[T]) *routineFuture[T] {
	return &routineFuture[T]{r}
}

// Cancel tries to cancel the future.
func (r *routineFuture[T]) Cancel() {
	r.Routine.Stop()
}
