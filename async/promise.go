package async

import (
	"sync"

	"github.com/epochtimeout/basekit/status"
)

// Promise is a completable future.
type Promise[T any] interface {
	Future[T]

	// Reject sets a promise error, returns false when already completed.
	Reject(st status.Status) bool

	// Resolve sets a promise result, returns false when already completed.
	Resolve(result T) bool

	// Complete completes the promise, returns false when already completed.
	Complete(result T, st status.Status) bool
}

// Pending returns a pending promise.
func Pending[T any]() Promise[T] {
	return newPromise[T]()
}

// Resolved returns a resolved promise.
func Resolved[T any](result T) Promise[T] {
	p := newPromise[T]()
	p.Resolve(result)
	return p
}

// Rejected returns a rejected promise.
func Rejected[T any](st status.Status) Promise[T] {
	p := newPromise[T]()
	p.Reject(st)
	return p
}

var _ Promise[any] = (*promise[any])(nil)

type promise[T any] struct {
	mu sync.Mutex

	done   bool
	result T
	status status.Status

	wait chan struct{}
}

func newPromise[T any]() *promise[T] {
	return &promise[T]{
		wait: make(chan struct{}),
	}
}

// Wait returns a channel which is closed on a future completion.
func (p *promise[T]) Wait() <-chan struct{} {
	return p.wait
}

// Result returns a future result and an error.
func (p *promise[T]) Result() (T, status.Status) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.result, p.status
}

// Status returns a future status when completed or an empty status.
func (p *promise[T]) Status() status.Status {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.status
}

// Reject sets a promise error, returns false when already completed.
func (p *promise[T]) Reject(st status.Status) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.done {
		return false
	}

	p.done = true
	p.status = st

	close(p.wait)
	return true
}

// Resolve sets a promise result, returns false when already completed.
func (p *promise[T]) Resolve(result T) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.done {
		return false
	}

	p.done = true
	p.result = result
	p.status = status.OK

	close(p.wait)
	return true
}

// Complete completes the promise, returns false when already completed.
func (p *promise[T]) Complete(result T, st status.Status) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.done {
		return false
	}

	p.done = true
	p.result = result
	p.status = st

	close(p.wait)
	return true
}
