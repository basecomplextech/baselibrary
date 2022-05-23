package async

import (
	"testing"

	"github.com/epochtimeout/basekit/system/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAll__should_await_all_results(t *testing.T) {
	rr := []Routine[int]{
		Call(func(stop <-chan struct{}) (int, status.Status) {
			return 1, status.OK
		}),
		Call(func(stop <-chan struct{}) (int, status.Status) {
			return 2, status.OK
		}),
		Call(func(stop <-chan struct{}) (int, status.Status) {
			return 3, status.OK
		}),
	}

	results, err := All(nil, rr...)
	if err != nil {
		t.Fatal(err)
	}

	require.Len(t, results, 3)
	assert.Equal(t, 1, results[0].Value)
	assert.Equal(t, 2, results[1].Value)
	assert.Equal(t, 3, results[2].Value)
}
