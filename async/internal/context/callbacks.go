// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package context

import (
	"sync"

	"github.com/basecomplextech/baselibrary/opt"
	"github.com/basecomplextech/baselibrary/status"
)

type callbacks struct {
	mu sync.Mutex
	m  opt.Opt[map[Callback]struct{}]
}

func (c *callbacks) init(m opt.Opt[map[Callback]struct{}]) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.m = m
}

func (c *callbacks) add(cb Callback, result *result) {
	// Maybe already done
	done, st := result.get()
	if done {
		cb.OnCancelled(st)
		return
	}

	// Lock callbacks
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check again
	done, st = result.get()
	if done {
		cb.OnCancelled(st)
		return
	}

	// Add callback
	m, ok := c.m.Unwrap()
	if !ok {
		m = make(map[Callback]struct{})
		c.m.Set(m)
	}
	m[cb] = struct{}{}
}

func (c *callbacks) remove(cb Callback) {
	c.mu.Lock()
	defer c.mu.Unlock()

	m, ok := c.m.Unwrap()
	if !ok {
		return
	}

	delete(m, cb)
}

func (c *callbacks) notify(st status.Status) {
	c.mu.Lock()
	defer c.mu.Unlock()

	m, ok := c.m.Unwrap()
	if !ok {
		return
	}

	for cb := range m {
		cb.OnCancelled(st)
	}
	clear(m)
}

func (c *callbacks) reset() opt.Opt[map[Callback]struct{}] {
	c.mu.Lock()
	defer c.mu.Unlock()

	m, ok := c.m.Unwrap()
	if ok {
		clear(m)
	}
	return c.m
}
