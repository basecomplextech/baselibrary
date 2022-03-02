package async

import "sync"

type Condition struct {
	mu   sync.Mutex
	ch   chan struct{}
	done bool
}

// NewCondition returns a new pending condition.
func NewCondition() *Condition {
	return &Condition{
		ch: make(chan struct{}),
	}
}

// DoneCondition returns a completed condition.
func DoneCondition() *Condition {
	c := NewCondition()
	c.Signal()
	return c
}

// Signal closes the wait channel.
func (c *Condition) Signal() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.done {
		return
	}

	close(c.ch)
	c.done = true
}

// Reset replaces the wait channel with an open one if closed.
func (c *Condition) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.done {
		return
	}

	c.ch = make(chan struct{})
	c.done = false
}

// Wait waits for the condition.
func (c *Condition) Wait() <-chan struct{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.ch
}
