package async

import "sync"

// Flag is a routine-safe boolean flag that can be set, reset, and waited on until set.
//
// Example:
//
//	serving := async.NewFlag()
//
//	func serve() {
//		s.serving.Signal()
//		defer s.serving.Reset()
//
//		// ... start server ...
//	}
//
//	func handle(cancel <-chan struct{}, req *request) <-chan struct{} {
//		// await server handling requests
//		select {
//		case <-serving.Wait():
//		case <-cancel:
//			return status.Cancelled
//		}
//
//		// ... handle request ...
//	}
type Flag struct {
	mu      sync.Mutex
	set     bool
	setChan chan struct{} // closed when set
}

// NewFlag returns a new unset flag.
func NewFlag() *Flag {
	return &Flag{
		setChan: make(chan struct{}),
	}
}

// SetFlag returns a new set flag.
func SetFlag() *Flag {
	f := NewFlag()
	f.Set()
	return f
}

// IsSet returns true if the flag is set.
func (f *Flag) IsSet() bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.set
}

// Set sets the flag, notifies waiters and returns true, or false if already set.
func (f *Flag) Set() bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.set {
		return false
	}

	close(f.setChan)
	f.set = true
	return true
}

// Reset resets the flag and replaces its wait channel with an open one.
func (f *Flag) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.set {
		return
	}

	f.setChan = make(chan struct{})
	f.set = false
}

// Wait waits for the flag to be set.
func (f *Flag) Wait() <-chan struct{} {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.setChan
}
