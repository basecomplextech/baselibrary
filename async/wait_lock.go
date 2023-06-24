package async

import "sync"

// WaitLock is a single writer multiple waiters lock.
//
// It allows one writer to acquire the lock, and other readers to wait until the lock is unlocked.
// WaitLock does not guarantee that the lock is not acquired by another writer after its waiters
// are notified.
//
// WaitLock can be used, for example, to execute a single operation by one writer, while other
// writers wait until the operation is completed.
//
// Example:
//
//	lock := async.NewWaitLock()
//
//	func flush(cancel <-chan struct{}) {
//		select {
//		case <-lock.Lock():
//			// Acquired lock
//		default:
//			// Await flushing end
//			select {
//			case <-lock.Wait():
//			case <-cancel:
//				return status.Cancelled
//			}
//		}
//		defer lock.Unlock()
//
//		// ... Do work ...
//	}
type WaitLock struct {
	mu sync.Mutex

	lock       chan struct{}
	wait       chan struct{} // Open on lock, closed on unlock
	waitClosed bool
}

// NewWaitLock returns a new unlocked lock.
func NewWaitLock() *WaitLock {
	l := &WaitLock{
		lock:       make(chan struct{}, 1),
		wait:       make(chan struct{}),
		waitClosed: true,
	}

	l.lock <- struct{}{}
	close(l.wait)
	return l
}

// Lock returns a channel receiving from which locks the lock.
func (l *WaitLock) Lock() <-chan struct{} {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.waitClosed {
		return l.lock
	}

	l.wait = make(chan struct{})
	l.waitClosed = false
	return l.lock
}

// Unlock unlocks the lock and notifies all waiters.
func (l *WaitLock) Unlock() {
	l.mu.Lock()
	defer l.mu.Unlock()

	select {
	case l.lock <- struct{}{}:
	default:
		panic("unlock of unlocked lock")
	}

	l.waitClosed = true
	close(l.wait)
}

// Wait returns a channel which is closed when the lock is unlocked.
func (l *WaitLock) Wait() <-chan struct{} {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.wait
}
