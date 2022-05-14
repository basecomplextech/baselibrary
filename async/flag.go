package async

import "sync"

type Flag struct {
	mu  sync.Mutex
	ch  chan struct{}
	set bool
}

// NewFlag returns a new unset flag.
func NewFlag() *Flag {
	return &Flag{
		ch: make(chan struct{}),
	}
}

// SetFlag returns a new set flag.
func SetFlag() *Flag {
	f := NewFlag()
	f.Signal()
	return f
}

// IsSet returns true if the flag is set.
func (f *Flag) IsSet() bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.set
}

// Signal sets the flag and closes its wait channel.
func (f *Flag) Signal() {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.set {
		return
	}

	close(f.ch)
	f.set = true
}

// Reset resets the flag and replaces its wait channel with an open one.
func (f *Flag) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.set {
		return
	}

	f.ch = make(chan struct{})
	f.set = false
}

// Wait waits for the flag to be set.
func (f *Flag) Wait() <-chan struct{} {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.ch
}
