package async

import (
	"sync"
	"time"

	"github.com/basecomplextech/baselibrary/status"
)

// Context is an async cancellation context.
//
// Usage:
//
//	ctx := NewContext()
//	defer ctx.Free()
type Context interface {
	// Cancel indicates that the operation should be cancelled.
	Done() <-chan struct{}

	// Cancel cancels the context.
	Cancel()

	// Status returns a cause or OK.
	Status() status.Status

	// Callbacks

	// AddCallback adds a callback.
	AddCallback(c ContextCallback)

	// RemoveCallback removes a callback.
	RemoveCallback(c ContextCallback)

	// Internal

	// Free cancels and releases the context.
	Free()
}

// ContextCallback allows to receive context done notifications.
type ContextCallback interface {
	// OnContextDone is called when the context is done.
	OnContextDone(status.Status)
}

// New

// NewContext returns a new pending context.
func NewContext() Context {
	return newContext(nil /* no parent */)
}

// NewContextTimeout returns a new context with a timeout.
func NewContextTimeout(timeout time.Duration) Context {
	return newContextTimeout(nil /* no parent */, timeout)
}

// NewContextDeadline returns a new context with a deadline.
func NewContextDeadline(deadline time.Time) Context {
	timeout := time.Until(deadline)
	return newContextTimeout(nil /* no parent */, timeout)
}

// Next

// NextContextTimeout returns a child context with a timeout.
func NextContextTimeout(parent Context, timeout time.Duration) Context {
	return newContextTimeout(parent, timeout)
}

// NextContextDeadline returns a child context with a deadline.
func NextContextDeadline(parent Context, deadline time.Time) Context {
	timeout := time.Until(deadline)
	return newContextTimeout(parent, timeout)
}

// channel

var _ Context = (*context)(nil)

type context struct {
	mu    sync.Mutex
	state *contextState
}

type contextState struct {
	parent Context // maybe nil

	done    bool
	cause   status.Status
	channel chan struct{}

	timer     *time.Timer                  // maybe nil
	callbacks map[ContextCallback]struct{} // maybe nil
}

func newContext(parent Context) *context {
	s := acquireContextState()
	s.parent = parent
	s.cause = status.Cancelled
	s.channel = make(chan struct{})
	c := &context{state: s}

	// Maybe add callback
	if parent != nil {
		parent.AddCallback(c)
	}
	return c
}

func newContextTimeout(parent Context, timeout time.Duration) *context {
	s := acquireContextState()
	s.parent = parent
	s.cause = status.Timeout
	s.channel = make(chan struct{})
	c := &context{state: s}

	// Maybe already done
	if timeout <= 0 {
		c.cancel(status.None)
		return c
	}

	// Start timer
	s.timer = time.AfterFunc(timeout, c.timeout)

	// Maybe Add callback
	if parent != nil {
		parent.AddCallback(c)
	}
	return c
}

// Cancel indicates that the operation should be cancelled.
func (c *context) Done() <-chan struct{} {
	s, ok := c.lock()
	defer c.mu.Unlock()

	if !ok || s.done {
		return closedChan
	}
	return s.channel
}

// Cancel cancels the context.
func (c *context) Cancel() {
	c.cancel(status.None)
}

// Status returns a cause or OK.
func (c *context) Status() status.Status {
	s, ok := c.lock()
	if !ok {
		return status.Cancelled
	}
	defer c.mu.Unlock()

	return s.cause
}

// Callbacks

// AddCallback adds a callback.
func (c *context) AddCallback(cb ContextCallback) {
	s, ok := c.lock()
	if !ok {
		cb.OnContextDone(status.Cancelled)
		return
	}
	defer c.mu.Unlock()

	// Maybe done
	if s.done {
		cb.OnContextDone(s.cause)
		return
	}

	// Add callback
	if s.callbacks == nil {
		s.callbacks = make(map[ContextCallback]struct{})
	}
	s.callbacks[cb] = struct{}{}
}

// RemoveCallback removes a callback.
func (c *context) RemoveCallback(cb ContextCallback) {
	s, ok := c.lock()
	if !ok {
		return
	}
	defer c.mu.Unlock()

	if s.callbacks != nil {
		delete(s.callbacks, cb)
	}
}

// OnContextDone is called when a parent context is done.
func (c *context) OnContextDone(st status.Status) {
	c.cancel(st)
}

// Internal

// Free cancels and releases the context.
func (c *context) Free() {
	c.cancel(status.None)

	s, ok := c.lock()
	if !ok {
		return
	}
	c.state = nil
	c.mu.Unlock()

	releaseContextState(s)
}

// internal

func (c *context) lock() (*contextState, bool) {
	c.mu.Lock()

	if c.state == nil {
		c.mu.Unlock()
		return nil, false
	}

	return c.state, true
}

func (c *context) cancel(st status.Status) {
	s, ok := c.lock()
	if !ok {
		return
	}

	// Locked
	{
		// Mark as done, close channel
		s.done = true
		close(s.channel)

		// Maybe set cause
		if st.Code != status.CodeNone {
			s.cause = st
		}

		// Maybe stop timer
		if s.timer != nil {
			s.timer.Stop()
		}
	}
	c.mu.Unlock()

	// Notify callbacks
	if len(s.callbacks) > 0 {
		for cb := range s.callbacks {
			cb.OnContextDone(s.cause)
		}
	}
}

func (c *context) timeout() {
	c.cancel(status.None)
}

// state pool

var contextStatePool = &sync.Pool{
	New: func() any {
		return &contextState{}
	},
}

func acquireContextState() *contextState {
	return contextStatePool.Get().(*contextState)
}

func releaseContextState(s *contextState) {
	callbacks := s.callbacks
	*s = contextState{}

	if callbacks != nil {
		clear(callbacks)
		s.callbacks = callbacks
	}

	contextStatePool.Put(s)
}
