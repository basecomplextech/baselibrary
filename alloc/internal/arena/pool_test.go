// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the Business Source License (BSL 1.1)
// that can be found in the LICENSE file.

package arena

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// NewPool

func TestNewPool__should_allocate_pool(t *testing.T) {
	a := Test()
	p := NewPool[int64](a)

	v0, _ := p.Get()
	*v0 = math.MaxInt64
	p.Put(v0)

	v1, _ := p.Get()
	assert.Equal(t, int64(math.MaxInt64), *v1)
}

func TestNewPool__should_return_different_pools_for_different_types_with_same_size(t *testing.T) {
	type Value struct {
		V int64
	}

	a := Test()
	list0 := NewPool[int64](a)
	list1 := NewPool[Value](a)

	assert.NotSame(t, list0, list1)
}

// Get

func TestPool_Get__should_allocate_new_object(t *testing.T) {
	a := Test()
	p := newPool[int64](a)

	v, ok := p.Get()
	*v = math.MaxInt64

	assert.Equal(t, int64(math.MaxInt64), *v)
	assert.False(t, ok)
}

func TestPool_Get__should_return_free_object(t *testing.T) {
	a := Test()
	p := newPool[int64](a)

	v0, ok := p.Get()
	p.Put(v0)
	assert.False(t, ok)

	v1, ok := p.Get()
	assert.Same(t, v0, v1)
	assert.True(t, ok)
}

func TestPool_Get__should_consume_free_item(t *testing.T) {
	a := Test()
	p := newPool[int64](a)

	v0, _ := p.Get()
	p.Put(v0)

	p.Get()
	assert.Nil(t, p.head.Load())
}

func TestPool_Get__should_swap_head_item_with_next(t *testing.T) {
	a := Test()
	p := newPool[int64](a)

	v0, _ := p.Get()
	v1, _ := p.Get()
	assert.Nil(t, p.head.Load())

	p.Put(v0)
	m := p.head.Load()
	assert.Same(t, v0, &m.obj)

	p.Put(v1)
	m = p.head.Load()
	assert.Same(t, v1, &m.obj)

	p.Get()
	m = p.head.Load()
	assert.Same(t, &m.obj, v0)
}

// Put

func TestPool_Put__should_swap_head_item(t *testing.T) {
	a := Test()
	p := newPool[int64](a)

	v0, _ := p.Get()
	v1, _ := p.Get()

	p.Put(v0)
	m := p.head.Load()
	assert.Same(t, v0, &m.obj)

	p.Put(v1)
	m = p.head.Load()
	assert.Same(t, v1, &m.obj)

	m1 := m.next.Load()
	assert.Same(t, v0, &m1.obj)
}
