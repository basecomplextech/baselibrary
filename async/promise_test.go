package async

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPromise_Resolve__should_complete_future(t *testing.T) {
	p := Pending[string]()

	ok := p.Resolve("hello")
	require.True(t, ok)

	result, err := p.Result()
	assert.Equal(t, "hello", result)
	assert.Nil(t, err)
}

func TestPromise_Reject__should_fail_future(t *testing.T) {
	p := Pending[string]()
	err := errors.New("test")

	ok := p.Reject(err)
	require.True(t, ok)

	result, err1 := p.Result()
	assert.Equal(t, "", result)
	assert.Equal(t, err, err1)
}

func TestPromise_Complete__should_resolve_promise(t *testing.T) {
	p := Pending[string]()

	ok := p.Complete("hello", nil)
	require.True(t, ok)

	result, err := p.Result()
	assert.Equal(t, "hello", result)
	assert.Nil(t, err)
}

func TestPromise_Complete__should_reject_promise(t *testing.T) {
	p := Pending[string]()
	err := errors.New("test")

	ok := p.Complete("failed", err)
	require.True(t, ok)

	result, err1 := p.Result()
	assert.Equal(t, "failed", result)
	assert.Equal(t, err, err1)
}
