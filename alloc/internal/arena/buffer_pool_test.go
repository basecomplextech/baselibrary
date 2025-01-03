// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package arena

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// NewBufferPool

func TestNewBufferPool__should_allocate_pool(t *testing.T) {
	a := Test()
	p := NewBufferPool(a)

	buf := p.Get()
	p.Put(buf)
}

// Get

func TestBufferPool_Get__should_allocate_new_buffer(t *testing.T) {
	a := Test()
	p := newBufferPool(a)

	buf0 := p.Get()
	buf0.Write([]byte("hello, world"))

	buf1 := p.Get()
	assert.Zero(t, buf1.Len())
}

func TestBufferPool_Get__should_return_free_buffer(t *testing.T) {
	a := Test()
	p := newBufferPool(a)

	buf0 := p.Get()
	buf0.Write([]byte("hello, world"))
	p.Put(buf0)

	buf1 := p.Get()
	assert.Same(t, buf0, buf1)
}

// Put

func TestBufferPool_Put__should_reset_buffer(t *testing.T) {
	a := Test()
	p := newBufferPool(a)

	buf := p.Get()
	buf.Write([]byte("hello, world"))
	p.Put(buf)

	assert.Zero(t, buf.Len())
}
