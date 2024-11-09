// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package context

import (
	"sync"

	"github.com/basecomplextech/baselibrary/collect/chans"
	"github.com/basecomplextech/baselibrary/opt"
	"github.com/basecomplextech/baselibrary/status"
)

type result struct {
	mu      sync.Mutex
	done    bool
	cause   status.Status
	channel opt.Opt[chan struct{}] // lazily created
}

func (r *result) get() (bool, status.Status) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.done, r.cause
}

func (r *result) set(st status.Status) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.done {
		return false
	}

	r.done = true
	r.cause = st

	c, ok := r.channel.Unwrap()
	if ok {
		close(c)
	}
	return true
}

func (r *result) wait() <-chan struct{} {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.done {
		return chans.Closed()
	}

	c, ok := r.channel.Unwrap()
	if !ok {
		c = make(chan struct{})
		r.channel.Set(c)
	}
	return c
}
