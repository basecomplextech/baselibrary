// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package context

import (
	"github.com/basecomplextech/baselibrary/opt"
	"github.com/basecomplextech/baselibrary/pools"
	"github.com/basecomplextech/baselibrary/status"
)

type state struct {
	parent    opt.Opt[Context]
	timer     timer
	result    result
	callbacks callbacks
}

func newState(parent Context) *state {
	s := statePool.New()
	s.parent = opt.Maybe(parent)
	return s
}

func (s *state) cancel(x *context, st status.Status) {
	// Set result
	ok := s.result.set(st)
	if !ok {
		return
	}

	// Stop timer
	s.timer.stop()

	// Remove from parent
	if p, ok := s.parent.Unwrap(); ok {
		p.RemoveCallback(x)
	}

	// Notify callbacks
	s.callbacks.notify(st)
}

func (s *state) reset() {
	m := s.callbacks.reset()

	*s = state{}
	s.callbacks.init(m)
}

// pool

var statePool = pools.NewPoolFunc(
	func() *state {
		return &state{}
	},
)

func releaseState(s *state) {
	s.reset()
	statePool.Put(s)
}
