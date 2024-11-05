// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

// LockedKey is a locked key which is unlocked when freed.
type LockedKey interface {
	// Free unlocks and freed the key.
	Free()
}

// internal

var _ LockedKey = &lockMapLockedKey[any]{}

type lockMapLockedKey[K comparable] struct {
	item *lockMapItem[K]
}

func newLockMapLockedKey[K comparable](item *lockMapItem[K]) *lockMapLockedKey[K] {
	return &lockMapLockedKey[K]{item: item}
}

func (l *lockMapLockedKey[K]) Free() {
	if l.item == nil {
		panic("free of freed locked key")
	}

	m := l.item
	l.item = nil

	m.unlock()
	m.release()
}
