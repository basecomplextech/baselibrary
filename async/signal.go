package async

import "sync"

type Signal struct {
	mu   sync.Mutex
	ch   chan struct{}
	done bool
}

// NewSignal returns a new pending signal.
func NewSignal() *Signal {
	return &Signal{
		ch: make(chan struct{}),
	}
}

// DoneSignal returns a completed signal.
func DoneSignal() *Signal {
	f := NewSignal()
	f.Emit()
	return f
}

// Emit closes the underlying signal channel.
func (f *Signal) Emit() {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.done {
		return
	}

	close(f.ch)
	f.done = true
}

// Reset replaces the flag channel with an open channel if closed.
func (f *Signal) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.done {
		return
	}

	f.ch = make(chan struct{})
	f.done = false
}

// Wait waits for the signal.
func (f *Signal) Wait() <-chan struct{} {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.ch
}
