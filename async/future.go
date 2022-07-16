package async

import (
	"sync"

	"github.com/epochtimeout/baselibrary/status"
)

// Future is an async request with a future result.
type Future[T any] interface {
	Result[T]

	// Cancel tries to cancel the future.
	Cancel() bool
}

// Promise is an async response with a future result.
type Promise[T any] interface {
	Future[T]

	// Cancelled returns a channel which is closed when the promise is cancelled.
	Cancelled() <-chan struct{}

	// Reject rejects the promise with a status.
	Reject(st status.Status) bool

	// Resolve resolve the promise with OK and a given result.
	Resolve(result T) bool

	// Complete completes the promise with a status and a result.
	Complete(result T, st status.Status) bool
}

// Pending returns a pending promise.
func Pending[T any]() Promise[T] {
	return newPromise[T]()
}

// Resolved returns a resolved future.
func Resolved[T any](result T) Future[T] {
	p := newPromise[T]()
	p.Resolve(result)
	return p
}

// Rejected returns a rejected future.
func Rejected[T any](st status.Status) Future[T] {
	p := newPromise[T]()
	p.Reject(st)
	return p
}

// Completed returns a completed future.
func Completed[T any](result T, st status.Status) Future[T] {
	p := newPromise[T]()
	p.Complete(result, st)
	return p
}

// CancelAll cancels all futures.
func CancelAll[T any](futures ...Future[T]) {
	for _, f := range futures {
		f.Cancel()
	}
}

// internal

var _ Promise[any] = (*promise[any])(nil)

type promise[T any] struct {
	mu sync.Mutex

	done   bool
	result T
	status status.Status

	wait      chan struct{}
	cancel    chan struct{}
	cancelled bool
}

func newPromise[T any]() *promise[T] {
	return &promise[T]{
		wait:   make(chan struct{}),
		cancel: make(chan struct{}),
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

// Cancel tries to cancel the future.
func (p *promise[T]) Cancel() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.done {
		return false
	}

	close(p.cancel)
	p.cancelled = true
	return true
}

// Cancelled returns a channel which is closed on a promise cancellation.
func (p *promise[T]) Cancelled() <-chan struct{} {
	return p.cancel
}

// Reject rejects the promise with a status.
func (p *promise[T]) Reject(st status.Status) bool {
	var result T
	return p.complete(result, st)
}

// Resolve resolve the promise with OK and a given result.
func (p *promise[T]) Resolve(result T) bool {
	return p.complete(result, status.OK)
}

// Complete completes the promise with a status and a result.
func (p *promise[T]) Complete(result T, st status.Status) bool {
	return p.complete(result, st)
}

// private

func (p *promise[T]) complete(result T, st status.Status) bool {
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
