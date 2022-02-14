package async

import (
	"sync"
)

// Promise is a completable future.
type Promise[T any] interface {
	Future[T]

	// Reject sets a promise error, returns false when already completed.
	Reject(err error) bool

	// Resolve sets a promise result, returns false when already completed.
	Resolve(result T) bool

	// Complete completes the promise, returns false when already completed.
	Complete(result T, err error) bool
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
func Rejected[T any](err error) Promise[T] {
	p := newPromise[T]()
	p.Reject(err)
	return p
}

var _ Promise[any] = (*promise[any])(nil)

type promise[T any] struct {
	mu sync.Mutex

	done   bool
	result T
	err    error

	wait chan struct{}
}

func newPromise[T any]() *promise[T] {
	return &promise[T]{
		wait: make(chan struct{}),
	}
}

// Err returns the future error or nil.
func (p *promise[T]) Err() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.err
}

// Wait returns a channel which is closed on a future completion.
func (p *promise[T]) Wait() <-chan struct{} {
	return p.wait
}

// Result returns a future result and an error.
func (p *promise[T]) Result() (T, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.result, p.err
}

// Reject sets a promise error, returns false when already completed.
func (p *promise[T]) Reject(err error) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.done {
		return false
	}

	p.done = true
	p.err = err

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

	close(p.wait)
	return true
}

// Complete completes the promise, returns false when already completed.
func (p *promise[T]) Complete(result T, err error) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.done {
		return false
	}

	p.done = true
	p.result = result
	p.err = err

	close(p.wait)
	return true
}
