// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"sync"
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

// Unlock unlocks the lock.
func (l Lock) Unlock() {
	select {
	case l <- struct{}{}:
	default:
		panic("unlock of unlocked lock")
	}
}
