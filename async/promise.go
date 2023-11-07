package async

import (
	"sync"

	"github.com/basecomplextech/baselibrary/status"
)

// Promise is a future which can be completed.
type Promise[T any] interface {
	Future[T]

	// Cancelled returns a channel which is closed when the promise is cancelled.
	Cancelled() <-chan struct{}

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
	mu sync.Mutex

	// channels are lazily allocated
	wait   chan struct{}
	cancel chan struct{}

	st        status.Status
	done      bool
	result    T
	cancelled bool
}

func newPromise[T any]() *promise[T] {
	return &promise[T]{}
}

// Cancelled returns a channel which is closed when the promise is cancelled.
func (p *promise[T]) Cancelled() <-chan struct{} {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cancelled {
		return closedChan
	}

	p.initCancel()
	return p.cancel
}

// Cancel requests the future to cancel and returns a wait channel.
func (p *promise[T]) Cancel() <-chan struct{} {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.done || p.cancelled {
		return closedChan
	}

	p.initCancel()
	p.initWait()

	p.cancelled = true
	close(p.cancel)
	return p.wait
}

// Wait returns a channel which is closed when the result is available.
func (p *promise[T]) Wait() <-chan struct{} {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.done {
		return closedChan
	}

	p.initWait()
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

	if p.wait != nil {
		close(p.wait)
	}
	return true
}

// Reject rejects the promise with a status.
func (p *promise[T]) Reject(st status.Status) bool {
	var zero T
	return p.Complete(zero, st)
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

// private

func (p *promise[T]) initCancel() {
	if p.cancel == nil {
		p.cancel = make(chan struct{})
	}
}

func (p *promise[T]) initWait() {
	if p.wait == nil {
		p.wait = make(chan struct{})
	}
}
