package async

import (
	"sync/atomic"
	"time"
)

// Timeout starts a timeout which closes a cancel channel after a duration.
// Do not forget to call Stop() to release resources.
//
//	Usage:
//
//	t := async.NewTimeout(5 * time.Second)
//	defer t.Stop()
//
//	result, st := doSomething(t.Cancel)
//	if st.Code == status.CodeCancelled {
//		if t.Done() {
//			return nil, status.Timeoutf("doSomething timeout")
//		}
//		return nil, st
//	}
type Timeout struct {
	Cancel <-chan struct{} // public channel closed on timeout or cancel

	reached atomic.Bool
	stopped atomic.Bool

	timer *time.Timer
	stop  chan struct{} // optional, channel to stop goroutine
}

// NewTimeout returns a timeout with a channel which is closed after a duration.
func NewTimeout(after time.Duration) *Timeout {
	timeout := make(chan struct{})

	t := &Timeout{Cancel: timeout}
	t.timer = time.AfterFunc(after, func() {
		close(timeout)
		t.reached.Store(true)
	})
	return t
}

// NewTimeoutCancel returns a timeout with combines a cancel channel and a duration.
func NewTimeoutCancel(cancel <-chan struct{}, after time.Duration) *Timeout {
	cancel_ := make(chan struct{})
	timeout := make(chan struct{})

	t := &Timeout{
		Cancel: cancel_,
		stop:   make(chan struct{}),
	}
	t.timer = time.AfterFunc(after, func() {
		close(timeout)
	})

	go func() {
		select {
		case <-t.stop:
			return
		case <-cancel:
		case <-timeout:
			t.reached.Store(true)
		}

		close(cancel_)
	}()

	return t
}

// Reached returns true if the timeout has been reached.
func (t *Timeout) Reached() bool {
	return t.reached.Load()
}

// Stop stops the timeout and releases its resources.
// The method does not close the cancel channel or sets the done flag.
func (t *Timeout) Stop() {
	ok := t.stopped.CompareAndSwap(false, true)
	if !ok {
		return
	}

	t.timer.Stop()
	if t.stop != nil {
		close(t.stop)
	}
}
