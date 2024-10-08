// Copyright 2021 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"testing"

	"github.com/basecomplextech/baselibrary/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPromise_Resolve__should_complete_future(t *testing.T) {
	p := NewPromise[string]()

	ok := p.Complete("hello", status.OK)
	require.True(t, ok)

	result, st := p.Result()
	assert.Equal(t, "hello", result)
	assert.Equal(t, status.CodeOK, st.Code)
	assert.Nil(t, st.Error)
}

func TestPromise_Reject__should_fail_future(t *testing.T) {
	p := NewPromise[string]()
	st := status.Test("test")

	ok := p.Complete("", st)
	require.True(t, ok)

	result, st1 := p.Result()
	assert.Equal(t, "", result)
	assert.Equal(t, st, st1)
}

func TestPromise_Complete__should_resolve_promise(t *testing.T) {
	p := NewPromise[string]()

	ok := p.Complete("hello", status.OK)
	require.True(t, ok)

	result, st := p.Result()
	assert.Equal(t, "hello", result)
	assert.True(t, st.OK())
}

func TestPromise_Complete__should_reject_promise(t *testing.T) {
	p := NewPromise[string]()
	st := status.Test("test")

	ok := p.Complete("failed", st)
	require.True(t, ok)

	result, st1 := p.Result()
	assert.Equal(t, "failed", result)
	assert.Equal(t, st, st1)
}
