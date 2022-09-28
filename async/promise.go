package async

import "github.com/complex1tech/baselibrary/status"

// Promise is a future which can be completed.
type Promise[T any] interface {
	Future[T]

	// Cancelled returns a channel which is closed when the promise is cancelled.
	Cancelled() <-chan struct{}

	// Complete completes the promise with a status and a result.
	Complete(result T, st status.Status) bool
}

// NewPromise returns a pending promise.
func NewPromise[T any]() Promise[T] {
	return newPromise[T]()
}

// Reject rejects the promise with a status.
func Reject[T any](p Promise[T], st status.Status) bool {
	var zero T
	return p.Complete(zero, st)
}

// Resolve resolves the promise with a result.
func Resolve[T any](p Promise[T], result T) bool {
	return p.Complete(result, status.OK)
}

// internal

var _ Promise[any] = (*promise[any])(nil)

type promise[T any] struct {
	result[T]

	cancel    chan struct{}
	cancelled bool
}

func newPromise[T any]() *promise[T] {
	return &promise[T]{
		result: newResult[T](),
		cancel: make(chan struct{}),
	}
}

// Cancel tries to cancel the future.
func (p *promise[T]) Cancel() <-chan struct{} {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.done {
		return p.wait
	}

	close(p.cancel)
	p.cancelled = true
	return p.wait
}

// Cancelled returns a channel which is closed on a promise cancellation.
func (p *promise[T]) Cancelled() <-chan struct{} {
	return p.cancel
}

// Complete completes the promise with a status and a result.
func (p *promise[T]) Complete(result T, st status.Status) bool {
	return p.complete(result, st)
}
