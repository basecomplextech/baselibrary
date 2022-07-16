package async

import (
	"testing"

	"github.com/epochtimeout/baselibrary/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPromise_Resolve__should_complete_future(t *testing.T) {
	p := Pending[string]()

	ok := p.Resolve("hello")
	require.True(t, ok)

	result, st := p.Result()
	assert.Equal(t, "hello", result)
	assert.Equal(t, status.CodeOK, st.Code)
	assert.Nil(t, st.Error)
}

func TestPromise_Reject__should_fail_future(t *testing.T) {
	p := Pending[string]()
	st := status.Error("test")

	ok := p.Reject(st)
	require.True(t, ok)

	result, st1 := p.Result()
	assert.Equal(t, "", result)
	assert.Equal(t, st, st1)
}

func TestPromise_Complete__should_resolve_promise(t *testing.T) {
	p := Pending[string]()

	ok := p.Complete("hello", status.OK)
	require.True(t, ok)

	result, st := p.Result()
	assert.Equal(t, "hello", result)
	assert.True(t, st.OK())
}

func TestPromise_Complete__should_reject_promise(t *testing.T) {
	p := Pending[string]()
	st := status.Error("test")

	ok := p.Complete("failed", st)
	require.True(t, ok)

	result, st1 := p.Result()
	assert.Equal(t, "failed", result)
	assert.Equal(t, st, st1)
}
