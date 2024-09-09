// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package pool

import (
	"github.com/basecomplextech/baselibrary/pools"
)

// Pool allows to reuse goroutines with preallocated big stacks.
type Pool interface {
	// Go runs a function in the pool.
	Go(fn func())
}

// New returns a new goroutine pool.
func New() Pool {
	return newPool()
}

// internal

var _ Pool = (*pool)(nil)

type pool struct {
	pool pools.Pool[*worker]
}

func newPool() *pool {
	p := &pool{}
	p.pool = pools.NewPoolFunc(func() *worker {
		return newWorker(p)
	})
	return p
}

// Go runs a function in the pool.
func (p *pool) Go(fn func()) {
	w := p.pool.New()
	w.Go(fn)
}

// release releases a worker to the pool.
func (p *pool) release(w *worker) {
	p.pool.Put(w)
}
