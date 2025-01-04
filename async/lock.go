// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import "github.com/basecomplextech/baselibrary/async/internal/lock"

// Lock is a channel-based lock, which be used directly as a channel, or via the Lock/Unlock methods.
// To lock a lock receive an empty struct from it, to unlock a lock send an empty struct to it.
//
// Example:
//
//	lock := async.NewLock()
//	select {
//	case <-lock:
//	case <-cancel:
//		return status.Cancelled
//	}
//	defer lock.Unlock()
type Lock = lock.Lock

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
type WaitLock = lock.WaitLock

// New

// NewLock returns a new unlocked lock.
func NewLock() Lock {
	return lock.NewLock()
}

// NewWaitLock returns a new unlocked lock.
func NewWaitLock() WaitLock {
	return lock.NewWaitLock()
}
