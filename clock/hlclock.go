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
	mu   sync.RWMutex
	wall int64 // can be accessed atomically by readers
	seq  uint32
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

	// Return the last time
	return c.load()
}

// Next returns a time for the next event, updates the last time.
func (c *hlClock) Next() pclock.HLTimestamp {
	c.mu.Lock()
	defer c.mu.Unlock()

	next := time.Now().UnixNano()
	if next > c.wall {
		c.seq = 0
		c.storeWall(next)
	} else {
		c.seq++
	}

	return pclock.HLTimestamp{Wall: c.wall, Seq: c.seq}
}

// Update updates the last time if another is greater, increments the sequence number.
func (c *hlClock) Update(t pclock.HLTimestamp) pclock.HLTimestamp {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch {
	case t.Wall == c.wall:
		c.seq = max(t.Seq, c.seq) + 1

	case t.Wall > c.wall:
		c.seq = t.Seq + 1
		c.wall = t.Wall
	}

	t1 := pclock.HLTimestamp{Wall: c.wall, Seq: c.seq}
	return t1
}

// private

func (c *hlClock) load() pclock.HLTimestamp {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return pclock.HLTimestamp{Wall: c.wall, Seq: c.seq}
}

func (c *hlClock) loadWall() int64 {
	return atomic.LoadInt64(&c.wall)
}

func (c *hlClock) storeWall(next int64) {
	atomic.StoreInt64(&c.wall, next)
}
