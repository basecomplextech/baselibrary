package async

import (
	"sync"
	"time"

	"github.com/basecomplextech/baselibrary/opt"
	"github.com/basecomplextech/baselibrary/pools"
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

	// Done returns true if the context is cancelled.
	Done() bool

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

// CancelledContext returns a cancelled context.
func CancelledContext() Context {
	return doneCtx
}

// Timeout

// TimeoutContext returns a context with a timeout.
func TimeoutContext(timeout time.Duration) Context {
	return newContextTimeout(nil /* no parent */, timeout)
}

// DeadlineContext returns a context with a deadline.
func DeadlineContext(deadline time.Time) Context {
	timeout := time.Until(deadline)
	return newContextTimeout(nil /* no parent */, timeout)
}

// Next

// NextContext returns a child context.
func NextContext(parent Context) Context {
	return newContext(parent)
}

// NextTimeoutContext returns a child context with a timeout.
func NextTimeoutContext(parent Context, timeout time.Duration) Context {
	return newContextTimeout(parent, timeout)
}

// NextDeadlineContext returns a child context with a deadline.
func NextDeadlineContext(parent Context, deadline time.Time) Context {
	timeout := time.Until(deadline)
	return newContextTimeout(parent, timeout)
}

// internal

var _ Context = (*context)(nil)

type context struct {
	cmu   sync.Mutex // cancel mutex, prevents data race between cancel/free
	smu   sync.Mutex // state mutex
	state *contextState
}

type contextState struct {
	parent Context // maybe nil

	done    bool
	cause   status.Status
	channel opt.Opt[chan struct{}] // lazily created

	timer     *time.Timer                  // maybe nil
	callbacks map[ContextCallback]struct{} // maybe nil
}

func newContext(parent Context) *context {
	s := acquireContextState()
	s.parent = parent
	s.cause = status.Cancelled
	c := &context{state: s}

	// Maybe add callback
	if parent != nil {
		parent.AddCallback(c)
	}
	return c
}

func newContextTimeout(parent Context, timeout time.Duration) *context {
	c := newContextTimeout1(parent, timeout)

	// Add callback outside of lock
	if parent != nil {
		parent.AddCallback(c)
	}
	return c
}

func newContextTimeout1(parent Context, timeout time.Duration) *context {
	s := acquireContextState()
	s.parent = parent
	s.cause = status.Timeout
	c := &context{state: s}

	// Maybe already done
	if timeout <= 0 {
		c.cancel(status.None)
		return c
	}

	// Start timer
	// Lock to prevent data race with immediate timeout.
	c.smu.Lock()
	defer c.smu.Unlock()

	s.timer = time.AfterFunc(timeout, c.timeout)
	return c
}

func (s *contextState) reset() {
	callbacks := s.callbacks
	*s = contextState{}

	if callbacks != nil {
		clear(callbacks)
		s.callbacks = callbacks
	}
}

// Cancel cancels the context.
func (c *context) Cancel() {
	c.cancel(status.None)
}

// Done returns true if the context is cancelled.
func (c *context) Done() bool {
	s, ok := c.lockState()
	if !ok {
		return true
	}
	defer c.smu.Unlock()

	return s.done
}

// Wait returns a channel which is closed when the context is cancelled.
func (c *context) Wait() <-chan struct{} {
	s, ok := c.lockState()
	defer c.smu.Unlock()

	if !ok || s.done {
		return closedChan
	}

	if !s.channel.Set {
		s.channel = opt.New(make(chan struct{}))
	}
	return s.channel.Value
}

// Status returns a cancellation status or OK.
func (c *context) Status() status.Status {
	s, ok := c.lockState()
	if !ok {
		return status.Cancelled
	}
	defer c.smu.Unlock()

	return s.cause
}

// Callbacks

// AddCallback adds a callback.
func (c *context) AddCallback(cb ContextCallback) {
	s, ok := c.lockState()
	if !ok {
		cb.OnCancelled(status.Cancelled)
		return
	}

	// Maybe done, notify immediately outside of lock
	if s.done {
		c.smu.Unlock()
		cb.OnCancelled(s.cause)
		return
	}

	// Add callback
	defer c.smu.Unlock()

	if s.callbacks == nil {
		s.callbacks = make(map[ContextCallback]struct{})
	}
	s.callbacks[cb] = struct{}{}
}

// RemoveCallback removes a callback.
func (c *context) RemoveCallback(cb ContextCallback) {
	s, ok := c.lockState()
	if !ok {
		return
	}
	defer c.smu.Unlock()

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
	c.cancel(status.None /* default */)

	c.cmu.Lock()
	defer c.cmu.Unlock()

	s, ok := c.lockState()
	if !ok {
		return
	}
	defer c.smu.Unlock()

	c.state = nil
	releaseContextState(s)
}

// internal

func (c *context) cancel(st status.Status) {
	c.cmu.Lock()
	defer c.cmu.Unlock()

	// Try to cancel
	s, ok := c.doCancel(st)
	if !ok {
		return
	}

	// State is immutable here. Notify callbacks,
	// parent outside of the state lock.

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

// private

func (c *context) lockState() (*contextState, bool) {
	c.smu.Lock()

	if c.state == nil {
		c.smu.Unlock()
		return nil, false
	}

	return c.state, true
}

func (c *context) doCancel(st status.Status) (*contextState, bool) {
	s, ok := c.lockState()
	if !ok {
		return nil, false
	}
	defer c.smu.Unlock()

	if s.done {
		return nil, false
	}

	// Mark as done, close
	s.done = true
	if s.channel.Set {
		close(s.channel.Value)
	}

	// Maybe set cause
	if st.Code != status.CodeNone {
		s.cause = st
	}

	// Maybe stop timer
	if s.timer != nil {
		s.timer.Stop()
	}
	return s, true
}

// done context

var doneCtx Context = &doneContext{}

type doneContext struct{}

func (*doneContext) Cancel()                           {}
func (*doneContext) Done() bool                        { return true }
func (*doneContext) Wait() <-chan struct{}             { return closedChan }
func (*doneContext) Status() status.Status             { return status.OK }
func (*doneContext) AddCallback(cb ContextCallback)    { cb.OnCancelled(status.Cancelled) }
func (*doneContext) RemoveCallback(cb ContextCallback) {}
func (*doneContext) Free()                             {}

// no context

var noCtx Context = &noContext{}

type noContext struct{}

func (*noContext) Cancel()                        {}
func (*noContext) Done() bool                     { return false }
func (*noContext) Wait() <-chan struct{}          { return nil }
func (*noContext) Status() status.Status          { return status.OK }
func (*noContext) AddCallback(ContextCallback)    {}
func (*noContext) RemoveCallback(ContextCallback) {}
func (*noContext) Free()                          {}

// state pool

var contextStatePool = pools.MakePool(
	func() *contextState {
		return &contextState{}
	},
)

func acquireContextState() *contextState {
	return contextStatePool.New()
}

func releaseContextState(s *contextState) {
	s.reset()
	contextStatePool.Put(s)
}
