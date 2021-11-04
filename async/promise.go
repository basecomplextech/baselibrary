package async

import (
	"sync"
)

// Promise is a completable future.
type Promise interface {
	Future

	// Stop returns a channel which is closed on a cancellation request.
	Stop() <-chan struct{}

	// Utility

	// OK tries to complete the promise with a result.
	OK(result interface{}) bool

	// Fail tries to complete the promise with an error.
	Fail(err error) bool

	// Exit tries to complete the promise as cancelled.
	Exit() bool

	// Complete fails the future if err is not nil, otherwise, completes the future.
	Complete(result interface{}, err error) bool
}

// NewPromise returns a pending promise.
func NewPromise() Promise {
	return newPromise()
}

// OK returns a completed promise.
func OK(result interface{}) Promise {
	p := newPromise()
	p.OK(result)
	return p
}

// Failed returns a failed promise.
func Failed(err error) Promise {
	p := newPromise()
	p.Fail(err)
	return p
}

// Cancelled returns a cancelled promise.
func Cancelled() Promise {
	p := newPromise()
	p.Cancel()
	p.Exit()
	return p
}

var _ Promise = (*promise)(nil)

type promise struct {
	mu sync.Mutex

	status Status
	result interface{}
	err    error

	done chan struct{}
	stop chan struct{}

	cancelled bool
}

func newPromise() *promise {
	return &promise{
		status: StatusPending,

		done: make(chan struct{}),
		stop: make(chan struct{}),
	}
}

// Future

// Err returns the future error or nil.
func (p *promise) Err() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.err
}

// Done awaits the completion.
func (p *promise) Done() <-chan struct{} {
	return p.done
}

// Result returns the current status, result and error.
func (p *promise) Result() (Status, interface{}, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.status, p.result, p.err
}

// Cancel tries to cancel the future.
func (p *promise) Cancel() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	switch {
	case p.status != StatusPending:
		return false
	case p.cancelled:
		return true
	}

	p.cancelled = true
	close(p.stop)
	return true
}

// Promise

// Stop returns a channel which is closed on a cancellation request.
func (p *promise) Stop() <-chan struct{} {
	return p.stop
}

// Utility

// OK tries to complete the promise with a result.
func (p *promise) OK(result interface{}) bool {
	return p.complete(StatusOK, result, nil)
}

// Fail tries to complete the promise with an error.
func (p *promise) Fail(err error) bool {
	return p.complete(StatusError, nil, err)
}

// Exit tries to complete the promise as cancelled.
func (p *promise) Exit() bool {
	return p.complete(StatusExit, nil, nil)
}

// Complete fails the future if err is not nil, otherwise, completes the future.
func (p *promise) Complete(result interface{}, err error) bool {
	status := StatusOK
	if err != nil {
		status = StatusError
	}
	return p.complete(status, result, err)
}

// private

func (p *promise) complete(status Status, result interface{}, err error) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	switch {
	case status == StatusPending:
		panic("invalid pending status")
	case p.status != StatusPending:
		return false
	}

	p.status = status
	p.result = result
	p.err = err

	close(p.done)
	return true
}
