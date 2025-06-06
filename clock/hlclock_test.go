// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package clock

import (
	"testing"
	"time"

	"github.com/basecomplextech/baselibrary/proto/pclock"
	"github.com/stretchr/testify/assert"
)

func TestHLClock_Read__should_return_current_time(t *testing.T) {
	c := newHLClock()
	now := c.Read()

	assert.NotZero(t, now.Wall)
	assert.Zero(t, now.Logic)
}

func TestHLClock_Read__should_return_last_time_if_now_less(t *testing.T) {
	a := pclock.HLTimestamp{
		Wall:  time.Now().UnixNano() + 1000_000,
		Logic: 123,
	}

	c := newHLClock()
	c.mu.Lock()
	c.wall = a.Wall
	c.logic = a.Logic
	c.mu.Unlock()

	b := c.Read()
	assert.Equal(t, a, b)
}

// Next

func TestHLClock_Next__should_return_next_time(t *testing.T) {
	c := newHLClock()
	a := c.Next()
	b := c.Next()

	assert.True(t, a.Less(b))
}

func TestHLClock_Next__should_increment_logic_when_now_less_than_last_wall(t *testing.T) {
	a := pclock.HLTimestamp{
		Wall:  time.Now().UnixNano() + 1000_000,
		Logic: 123,
	}

	c := newHLClock()
	c.mu.Lock()
	c.wall = a.Wall
	c.logic = a.Logic
	c.mu.Unlock()

	b := c.Next()
	assert.Equal(t, a.Wall, b.Wall)
	assert.Equal(t, a.Logic+1, b.Logic)
}

func TestHLClock_Next__should_update_last_time(t *testing.T) {
	c := newHLClock()
	a := c.Next()
	b := pclock.HLTimestamp{Wall: c.wall, Logic: c.logic}

	assert.Equal(t, a, b)
}

// Update

func TestHLClock_Update__should_update_last_time_increment_logicuence(t *testing.T) {
	a := pclock.HLTimestamp{
		Wall:  time.Now().UnixNano() + 1000_000,
		Logic: 123,
	}

	c := newHLClock()

	b := c.Update(a)
	assert.Equal(t, a.Wall, b.Wall)
	assert.Equal(t, a.Logic+1, b.Logic)
}

func TestHLClock_Update__should_increment_logic_when_equal_wall_times(t *testing.T) {
	a := pclock.HLTimestamp{
		Wall:  time.Now().UnixNano() + 1000_000,
		Logic: 123,
	}

	c := newHLClock()

	_ = c.Update(a)
	b := c.Update(a)
	assert.Equal(t, a.Wall, b.Wall)
	assert.Equal(t, a.Logic+2, b.Logic)
}
