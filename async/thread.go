package async

import (
	"github.com/complex1tech/baselibrary/status"
)

// Thread is a generic concurrent thread which returns the result as a future,
// and can be cancelled.
type Thread[T any] interface {
	Future[T]

	// Cancel requests the thread to cancel and returns a wait channel.
	Cancel() <-chan struct{}
}

// Methods

// Run runs a function in a new thread, recovers on panics.
func Run(fn func(cancel <-chan struct{}) status.Status) Thread[struct{}] {
	th := newThread[struct{}]()

	go func() {
		defer func() {
			e := recover()
			if e == nil {
				return
			}

			st := status.Recover(e)
			th.Complete(struct{}{}, st)
		}()

		st := fn(th.cancel)
		th.Complete(struct{}{}, st)
	}()

	return th
}

// RunSelf runs a function in a new thread, recovers on panics.
func RunSelf(fn func(cancel <-chan struct{}, th Thread[struct{}]) status.Status) Thread[struct{}] {
	th := newThread[struct{}]()

	go func() {
		defer func() {
			e := recover()
			if e == nil {
				return
			}

			st := status.Recover(e)
			th.Complete(struct{}{}, st)
		}()

		st := fn(th.cancel, th)
		th.Complete(struct{}{}, st)
	}()

	return th
}

// Execute executes a function in a new thread, recovers on panics.
func Execute[T any](fn func(cancel <-chan struct{}) (T, status.Status)) Thread[T] {
	th := newThread[T]()

	go func() {
		defer func() {
			e := recover()
			if e == nil {
				return
			}

			var zero T
			st := status.Recover(e)
			th.Complete(zero, st)
		}()

		result, st := fn(th.cancel)
		th.Complete(result, st)
	}()

	return th
}

// Join joins all threades into a single thread.
// The thread returns all the results and the first non-OK status.
func Join[T any](threads ...Thread[T]) Thread[[]T] {
	return Execute(func(cancel <-chan struct{}) ([]T, status.Status) {
		// Await all or cancel
		st := status.OK
	loop:
		for _, th := range threads {
			select {
			case <-th.Wait():
			case <-cancel:
				st = status.Cancelled
				break loop
			}
		}

		// Cancel all
		for _, th := range threads {
			th.Cancel()
		}

		// Collect results
		results := make([]T, 0, len(threads))
		for _, th := range threads {
			<-th.Wait()

			r, st1 := th.Result()
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
	_ Thread[any]  = (*thread[any])(nil)
	_ CancelWaiter = (*thread[any])(nil)
)

type thread[T any] struct {
	promise[T]

	cancel    chan struct{}
	cancelled bool
}

func newThread[T any]() *thread[T] {
	return &thread[T]{
		promise: promise[T]{
			wait: make(chan struct{}),
		},
		cancel: make(chan struct{}),
	}
}

// Cancel requests the thread to cancel and returns a wait channel.
func (th *thread[T]) Cancel() <-chan struct{} {
	th.mu.Lock()
	defer th.mu.Unlock()

	if th.done || th.cancelled {
		return th.wait
	}

	close(th.cancel)
	th.cancelled = true
	return th.wait
}
