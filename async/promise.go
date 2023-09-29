package async

import (
	"sync"

	"github.com/basecomplextech/baselibrary/status"
)

// Promise is a future which can be completed.
type Promise[T any] interface {
	Future[T]

	// Complete completes the promise with a status and a result.
	Complete(result T, st status.Status) bool

	// Reject rejects the promise with a status.
	Reject(st status.Status) bool
}

// NewPromise returns a pending promise.
func NewPromise[T any]() Promise[T] {
	return newPromise[T]()
}

// internal

var _ Promise[any] = (*promise[any])(nil)

type promise[T any] struct {
	mu   sync.Mutex
	wait chan struct{}

	st     status.Status
	done   bool
	result T
}

func newPromise[T any]() *promise[T] {
	return &promise[T]{
		wait: make(chan struct{}),
	}
}

// Result returns a value and a status.
func (p *promise[T]) Result() (T, status.Status) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.result, p.st
}

// Status returns a status.
func (p *promise[T]) Status() status.Status {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.st
}

// Wait returns a channel which is closed when the result is available.
func (p *promise[T]) Wait() <-chan struct{} {
	return p.wait
}

// Complete completes the promise with a status and a result.
func (p *promise[T]) Complete(result T, st status.Status) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.done {
		return false
	}

	p.result = result
	p.st = st
	p.done = true

	close(p.wait)
	return true
}

// Reject rejects the promise with a status.
func (p *promise[T]) Reject(st status.Status) bool {
	var zero T
	return p.Complete(zero, st)
}
