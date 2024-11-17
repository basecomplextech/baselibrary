// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"sync"

	"github.com/basecomplextech/baselibrary/status"
)

var _ sync.Locker = (Lock)(nil)

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
type Lock chan struct{}

// NewLock returns a new unlocked lock.
func NewLock() Lock {
	l := make(Lock, 1)
	l <- struct{}{}
	return l
}

// Lock locks the lock.
func (l Lock) Lock() {
	<-l
}

// LockContext awaits and locks the lock, or awaits the context cancellation.
func (l Lock) LockContext(ctx Context) status.Status {
	// Try to lock the lock without waiting for the context.
	// Context wait lazily allocates the internal channel.
	select {
	case <-l:
		return status.OK
	default:
	}

	select {
	case <-l:
		return status.OK
	case <-ctx.Wait():
		return ctx.Status()
	}
}

// Unlock unlocks the lock, or panics if the lock is already unlocked.
func (l Lock) Unlock() {
	select {
	case l <- struct{}{}:
	default:
		panic("unlock of unlocked lock")
	}
}

// UnlockIfLocked unlocks the lock if it is locked, otherwise does nothing.
func (l Lock) UnlockIfLocked() {
	select {
	case l <- struct{}{}:
	default:
	}
}
