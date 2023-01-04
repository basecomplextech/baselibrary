package async

import (
	"github.com/complex1tech/baselibrary/panics"
	"github.com/complex1tech/baselibrary/status"
)

var _ CancelWaiter = (Process[any])(nil)

// Process is a generic concurrent process which returns the result as a future,
// and can be cancelled.
type Process[T any] interface {
	Future[T]

	// Cancel requests the process to cancel and returns a wait channel.
	Cancel() <-chan struct{}
}

// Run runs a function in a new process, recovers on panics.
func Run(fn func(cancel <-chan struct{}) status.Status) Process[struct{}] {
	p := newProcess[struct{}]()

	go func() {
		defer func() {
			e := recover()
			if e == nil {
				return
			}

			err := panics.Recover(e)
			st := status.WrapError(err)
			p.Complete(struct{}{}, st)
		}()

		st := fn(p.cancel)
		p.Complete(struct{}{}, st)
	}()

	return p
}

// RunSelf runs a function in a new process, recovers on panics.
func RunSelf(fn func(cancel <-chan struct{}, p Process[struct{}]) status.Status) Process[struct{}] {
	p := newProcess[struct{}]()

	go func() {
		defer func() {
			e := recover()
			if e == nil {
				return
			}

			err := panics.Recover(e)
			st := status.WrapError(err)
			p.Complete(struct{}{}, st)
		}()

		st := fn(p.cancel, p)
		p.Complete(struct{}{}, st)
	}()

	return p
}

// Execute executes a function in a new process, recovers on panics.
func Execute[T any](fn func(cancel <-chan struct{}) (T, status.Status)) Process[T] {
	p := newProcess[T]()

	go func() {
		defer func() {
			e := recover()
			if e == nil {
				return
			}

			err := panics.Recover(e)
			st := status.WrapError(err)

			var zero T
			p.Complete(zero, st)
		}()

		result, st := fn(p.cancel)
		p.Complete(result, st)
	}()

	return p
}

// internal

type process[T any] struct {
	promise[T]

	cancel    chan struct{}
	cancelled bool
}

func newProcess[T any]() *process[T] {
	return &process[T]{
		promise: promise[T]{
			wait: make(chan struct{}),
		},
		cancel: make(chan struct{}),
	}
}

// Cancel requests the process to cancel and returns a wait channel.
func (p *process[T]) Cancel() <-chan struct{} {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.done || p.cancelled {
		return p.wait
	}

	close(p.cancel)
	p.cancelled = true
	return p.wait
}
