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
	// Cancel cancels the context.
	Cancel()

	// Wait returns a channel which is closed when the context is cancelled.
	Wait() <-chan struct{}

	// Status returns a cancellation status or OK.
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

// ContextCallback receives context cancellation notifications.
type ContextCallback interface {
	// OnCancelled is called when the context is cancelled.
	OnCancelled(status.Status)
}

// New

// NewContext returns a new context.
func NewContext() Context {
	return newContext(nil /* no parent */)
}

// NoContext returns a non-cancellable background context.
func NoContext() Context {
	return noCtx
}

// DoneContext returns a cancelled context.
func DoneContext() Context {
	return doneCtx
}

// More

// NewContextTimeout returns a new context with a timeout.
func NewContextTimeout(timeout time.Duration) Context {
	return newContextTimeout(nil /* no parent */, timeout)
}

// NewContextDeadline returns a new context with a deadline.
func NewContextDeadline(deadline time.Time) Context {
	timeout := time.Until(deadline)
	return newContextTimeout(nil /* no parent */, timeout)
}

// Child

// ChildContextTimeout returns a child context with a timeout.
func ChildContextTimeout(parent Context, timeout time.Duration) Context {
	return newContextTimeout(parent, timeout)
}

// ChildContextDeadline returns a child context with a deadline.
func ChildContextDeadline(parent Context, deadline time.Time) Context {
	timeout := time.Until(deadline)
	return newContextTimeout(parent, timeout)
}

// internal

var _ Context = (*context)(nil)

type context struct {
	mu    sync.Mutex
	state *contextState
}

type contextState struct {
	parent Context // maybe nil

	done    bool
	doneMu  sync.Mutex // enforce single canceller
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
	c := newContextTimeout1(parent, timeout)

	// Maybe add callback outside of lock
	if parent != nil {
		parent.AddCallback(c)
	}
	return c
}

func newContextTimeout1(parent Context, timeout time.Duration) *context {
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
	c.mu.Lock()
	defer c.mu.Unlock()

	s.timer = time.AfterFunc(timeout, c.timeout)
	return c
}

// Cancel cancels the context.
func (c *context) Cancel() {
	c.cancel(status.None)
}

// Wait returns a channel which is closed when the context is cancelled.
func (c *context) Wait() <-chan struct{} {
	s, ok := c.lock()
	defer c.mu.Unlock()

	if !ok || s.done {
		return closedChan
	}
	return s.channel
}

// Status returns a cancellation status or OK.
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
		cb.OnCancelled(status.Cancelled)
		return
	}
	defer c.mu.Unlock()

	// Maybe done
	if s.done {
		cb.OnCancelled(s.cause)
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

// OnCancelled is called when a parent context is done.
func (c *context) OnCancelled(st status.Status) {
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

	// Enforce single canceller til the end, but do not hold the `s.mu` lock
	// while calling callbacks. Mostly doneMu exists to prevent race conditions
	// between cancel() and Free().
	s.doneMu.Lock()
	defer s.doneMu.Unlock()

	// Maybe already done
	if s.done {
		c.mu.Unlock()
		return
	}

	// Mark as done, close channel
	{
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
			cb.OnCancelled(s.cause)
		}
	}

	// Remove from parent
	if s.parent != nil {
		s.parent.RemoveCallback(c)
	}
}

func (c *context) timeout() {
	c.cancel(status.None)
}

// done context

var doneCtx Context = &doneContext{}

type doneContext struct{}

func (*doneContext) Cancel()                           {}
func (*doneContext) Wait() <-chan struct{}             { return closedChan }
func (*doneContext) Status() status.Status             { return status.OK }
func (*doneContext) AddCallback(cb ContextCallback)    { cb.OnCancelled(status.Cancelled) }
func (*doneContext) RemoveCallback(cb ContextCallback) {}
func (*doneContext) Free()                             {}

// no context

var noCtx Context = &noContext{}

type noContext struct{}

func (*noContext) Cancel()                        {}
func (*noContext) Wait() <-chan struct{}          { return nil }
func (*noContext) Status() status.Status          { return status.OK }
func (*noContext) AddCallback(ContextCallback)    {}
func (*noContext) RemoveCallback(ContextCallback) {}
func (*noContext) Free()                          {}

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
