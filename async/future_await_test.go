// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"testing"
	"time"

	"github.com/basecomplextech/baselibrary/status"
	"github.com/stretchr/testify/assert"
)

func TestAwaitAll__should_wait_for_all_futures(t *testing.T) {
	r0 := Run(func(ctx Context) (int, status.Status) {
		return 1, status.OK
	})
	r1 := Run(func(ctx Context) (int, status.Status) {
		return 2, status.OK
	})

	ctx := NoContext()
	results, st := AwaitAll(ctx, r0, r1)
	if !st.OK() {
		t.Fatal(st)
	}

	assert.Equal(t, []Result[int]{{1, status.OK}, {2, status.OK}}, results)
}

func TestAwaitAll__should_accept_slice_of_routines(t *testing.T) {
	r0 := Run(func(ctx Context) (int, status.Status) {
		return 1, status.OK
	})
	r1 := Run(func(ctx Context) (int, status.Status) {
		return 2, status.OK
	})

	ctx := NoContext()
	routines := []Routine[int]{r0, r1}
	results, st := AwaitAll(ctx, routines...)
	if !st.OK() {
		t.Fatal(st)
	}

	assert.Equal(t, []Result[int]{{1, status.OK}, {2, status.OK}}, results)
}

// AwaitAny

func TestAwaitAny__should_wait_for_any_future(t *testing.T) {
	r0 := Run(func(ctx Context) (int, status.Status) {
		time.Sleep(time.Millisecond * 50)
		return 1, status.OK
	})
	r1 := Run(func(ctx Context) (int, status.Status) {
		time.Sleep(time.Millisecond * 100)
		return 2, status.OK
	})
	r2 := Run(func(ctx Context) (int, status.Status) {
		return 3, status.OK
	})

	ctx := NoContext()
	i, result, st := AwaitAny(ctx, r0, r1, r2)
	if !st.OK() {
		t.Fatal(st)
	}
	assert.Equal(t, 2, i)
	assert.Equal(t, 3, result)
}

// AwaitError

func TestAwaitError__should_wait_for_any_error_error(t *testing.T) {
	r0 := Run(func(ctx Context) (int, status.Status) {
		return 1, status.OK
	})
	r1 := Run(func(ctx Context) (int, status.Status) {
		time.Sleep(time.Millisecond * 10)
		return 2, status.Error("test error")
	})
	r2 := Run(func(ctx Context) (int, status.Status) {
		time.Sleep(time.Millisecond * 20)
		return 3, status.OK
	})

	ctx := NoContext()
	i, st := AwaitError(ctx, r0, r1, r2)
	assert.Equal(t, 1, i)
	assert.Equal(t, status.Error("test error"), st)
}
