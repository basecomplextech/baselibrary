package async

import (
	"sync"
	"sync/atomic"
)

// Flag is a routine-safe boolean flag that can be set, reset, and waited on until set.
//
// Example:
//
//	serving := async.UnsetFlag()
//
//	func serve() {
//		s.serving.Set()
//		defer s.serving.Unset()
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
	set     atomic.Bool
	setChan chan struct{} // closed when set
}

// SetFlag returns a new set flag.
func SetFlag() *Flag {
	f := UnsetFlag()
	f.Set()
	return f
}

// UnsetFlag returns a new unset flag.
func UnsetFlag() *Flag {
	return &Flag{
		setChan: make(chan struct{}),
	}
}

// IsSet returns true if the flag is set.
func (f *Flag) IsSet() bool {
	return f.set.Load()
}

// Set sets the flag, notifies waiters and returns true, or false if already set.
func (f *Flag) Set() bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.set.Load() {
		return false
	}

	f.set.Store(true)
	close(f.setChan)
	return true
}

// Unset unsets the flag and replaces its wait channel with an open one.
func (f *Flag) Unset() {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.set.Load() {
		return
	}

	f.set.Store(false)
	f.setChan = make(chan struct{})
}

// Wait waits for the flag to be set.
func (f *Flag) Wait() <-chan struct{} {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.setChan
}
