package async

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPromise_Cancel__should_close_stop_channel(t *testing.T) {
	p := NewPromise()
	ok := p.Cancel()
	require.True(t, ok)

	ch := p.Stop()
	select {
	case <-ch:
	default:
		t.Fatal("stop channel not closed")
	}
}

func TestPromise_OK__should_complete_future(t *testing.T) {
	p := NewPromise()
	ok := p.OK("hello")
	require.True(t, ok)

	status, result, err := p.Result()
	assert.Equal(t, StatusOK, status)
	assert.Equal(t, "hello", result)
	assert.Nil(t, err)
}

func TestPromise_Fail__should_fail_future(t *testing.T) {
	p := NewPromise()
	err := errors.New("test")
	ok := p.Fail(err)
	require.True(t, ok)

	status, result, err1 := p.Result()
	assert.Equal(t, StatusError, status)
	assert.Nil(t, result)
	assert.Equal(t, err, err1)
}

func TestPromise_Exit__should_cancel_future(t *testing.T) {
	p := NewPromise()
	ok := p.Exit()
	require.True(t, ok)

	status, result, err := p.Result()
	assert.Equal(t, StatusExit, status)
	assert.Nil(t, result)
	assert.Nil(t, err)
}

func TestPromise_Complete__should_complete_future_when_error_nil(t *testing.T) {
	p := NewPromise()
	ok := p.Complete("hello", nil)
	require.True(t, ok)

	status, result, err := p.Result()
	assert.Equal(t, StatusOK, status)
	assert.Equal(t, "hello", result)
	assert.Nil(t, err)
}

func TestPromise_Complete__should_fail_future_when_error_nonnil(t *testing.T) {
	p := NewPromise()
	err := errors.New("test")
	ok := p.Complete("failed", err)
	require.True(t, ok)

	status, result, err1 := p.Result()
	assert.Equal(t, StatusError, status)
	assert.Equal(t, "failed", result)
	assert.Equal(t, err, err1)
}
