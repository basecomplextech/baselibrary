// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package clock

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/basecomplextech/baselibrary/proto/pclock"
)

// HLClock is a thread-safe hybrid logical clock.
//
// See "Logical Physical Clocks and Consistent Snapshots in Globally Distributed Databases"
// https://cse.buffalo.edu/tech-reports/2014-04.pdf
type HLClock interface {
	// Read returns the current time, does not update the last time, non-blocking.
	Read() pclock.HLTimestamp

	// Next returns a time for the next event, updates the last time.
	Next() pclock.HLTimestamp

	// Update updates the last time if another is greater, increments the sequence number.
	Update(t pclock.HLTimestamp) pclock.HLTimestamp
}

// NewHLClock returns a new hybrid logical clock.
func NewHLClock() HLClock {
	return newHLClock()
}

// internal

var _ HLClock = (*hlClock)(nil)

type hlClock struct {
	mu    sync.RWMutex
	wall  int64 // can be accessed atomically by readers
	logic uint32
}

func newHLClock() *hlClock {
	return &hlClock{}
}

// Read returns the current time, does not update the last time, non-blocking.
func (c *hlClock) Read() pclock.HLTimestamp {
	// Return now if greater than last
	now := time.Now().UnixNano()
	last := c.loadWall()
	if now > last {
		return pclock.HLTimestamp{Wall: now}
	}

	// Return last time
	c.mu.RLock()
	defer c.mu.RUnlock()

	return pclock.HLTimestamp{Wall: c.wall, Logic: c.logic}
}

// Next returns a time for the next event, updates the last time.
func (c *hlClock) Next() pclock.HLTimestamp {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Update last time, or increment logic
	next := time.Now().UnixNano()
	if next > c.wall {
		c.storeWall(next)
		c.logic = 0
	} else {
		c.logic++
	}

	// Return last time
	return pclock.HLTimestamp{Wall: c.wall, Logic: c.logic}
}

// Update updates the last time if another is greater, increments the sequence number.
func (c *hlClock) Update(t pclock.HLTimestamp) pclock.HLTimestamp {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch {
	case t.Wall == c.wall:
		c.logic = max(t.Logic, c.logic) + 1

	case t.Wall > c.wall:
		c.storeWall(t.Wall)
		c.logic = t.Logic + 1
	}

	return pclock.HLTimestamp{Wall: c.wall, Logic: c.logic}
}

// private

func (c *hlClock) loadWall() int64 {
	return atomic.LoadInt64(&c.wall)
}

func (c *hlClock) storeWall(next int64) {
	atomic.StoreInt64(&c.wall, next)
}
