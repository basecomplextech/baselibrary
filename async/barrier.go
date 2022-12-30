package async

import "sync"

// Barrier is a thread-safe barrier that can be locked and waited on until unlocked.
//
// Example:
//
//	flushing := async.NewBarrier()
//
//	func flush(cancel <-chan struct{}) {
//		select {
//		case <-flushing.Lock():
//			// acquired barrier
//		default:
//			// await flushing end
//			select {
//			case <-flushing.Wait():
//			case <-cancel:
//				return status.Cancelled
//			}
//		}
//		defer flushing.Unlock()
//
//		// ... do work ...
//	}
type Barrier struct {
	mu sync.Mutex

	lock       chan struct{}
	wait       chan struct{} // closed when unlocked
	waitClosed bool
}

// NewBarrier returns a new unlocked barrier.
func NewBarrier() *Barrier {
	b := &Barrier{
		lock:       make(chan struct{}, 1),
		wait:       make(chan struct{}),
		waitClosed: true,
	}

	b.lock <- struct{}{}
	close(b.wait)
	return b
}

// Lock returns a channel receiving from which locks the barrier.
func (b *Barrier) Lock() <-chan struct{} {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.waitClosed {
		return b.lock
	}

	b.wait = make(chan struct{})
	b.waitClosed = false
	return b.lock
}

// Unlock unlocks the barrier and notifies all waiters.
func (b *Barrier) Unlock() {
	b.mu.Lock()
	defer b.mu.Unlock()

	select {
	case b.lock <- struct{}{}:
	default:
		panic("unlock of unlocked barrier")
	}

	b.waitClosed = true
	close(b.wait)
}

// Wait returns a channel which is closed when the barrier is unlocked.
func (b *Barrier) Wait() <-chan struct{} {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.wait
}
