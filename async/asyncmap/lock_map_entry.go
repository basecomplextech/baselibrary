// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import "github.com/basecomplextech/baselibrary/opt"

type lockMapEntry[K comparable] struct {
	item opt.Opt[*lockMapItem[K]]
	more opt.Opt[map[K]*lockMapItem[K]]
}

func (e *lockMapEntry[K]) get(key K) (*lockMapItem[K], bool) {
	if m, ok := e.item.Unwrap(); ok {
		if m.key == key {
			return m, true
		}
	}
	if more, ok := e.more.Unwrap(); ok {
		m, ok := more[key]
		if ok {
			return m, true
		}
	}
	return nil, false
}

func (e *lockMapEntry[K]) set(m *lockMapItem[K]) {
	// Maybe set single item
	if !e.item.Valid {
		e.item.Set(m)
		return
	}

	// Othewise add to more items
	more, ok := e.more.Unwrap()
	if !ok {
		more = make(map[K]*lockMapItem[K])
		e.more.Set(more)
	}
	more[m.key] = m
}

func (e *lockMapEntry[K]) delete(m *lockMapItem[K]) {
	if m1, ok := e.item.Unwrap(); ok {
		if m1 == m {
			e.item.Clear()
			return
		}
	}
	if more, ok := e.more.Unwrap(); ok {
		delete(more, m.key)
	}
}

func (e *lockMapEntry[K]) contains(key K) bool {
	if m, ok := e.item.Unwrap(); ok {
		if m.key == key {
			return true
		}
	}
	if more, ok := e.more.Unwrap(); ok {
		_, ok := more[key]
		return ok
	}
	return false
}

func (e *lockMapEntry[K]) range_(fn func(key K) bool) {
	if m, ok := e.item.Unwrap(); ok {
		if !fn(m.key) {
			return
		}
	}
	if more, ok := e.more.Unwrap(); ok {
		for key := range more {
			if !fn(key) {
				break
			}
		}
	}
}
