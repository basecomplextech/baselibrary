// Copyright 2023 Ivan Korobkov. All rights reserved.

package async

import "sync"

// WaitLock is a lock which allows others to wait until it is unlocked.
//
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
type WaitLock interface {
	// Lock returns a channel receiving from which locks the lock.
	Lock() <-chan struct{}

	// Unlock unlocks the lock and notifies all waiters.
	Unlock()

	// Wait returns a channel which is closed when the lock is unlocked.
	Wait() <-chan struct{}
}

// NewWaitLock returns a new unlocked lock.
func NewWaitLock() WaitLock {
	l := &waitLock{
		lock:       make(chan struct{}, 1),
		wait:       make(chan struct{}),
		waitClosed: true,
	}

	l.lock <- struct{}{}
	close(l.wait)
	return l
}

// internal

type waitLock struct {
	mu sync.Mutex

	lock       chan struct{}
	wait       chan struct{} // Open on lock, closed on unlock
	waitClosed bool
}

// Lock returns a channel receiving from which locks the lock.
func (l *waitLock) Lock() <-chan struct{} {
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
func (l *waitLock) Unlock() {
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
func (l *waitLock) Wait() <-chan struct{} {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.wait
}
