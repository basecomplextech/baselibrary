package async

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoutines_Wait__should_await_results(t *testing.T) {
	rr := Routines[int]{
		Call(func(stop <-chan struct{}) (int, error) {
			return 1, nil
		}),
		Call(func(stop <-chan struct{}) (int, error) {
			return 2, nil
		}),
		Call(func(stop <-chan struct{}) (int, error) {
			return 3, nil
		}),
	}

	all := rr.Wait()
	<-all.Wait()

	results, err := all.Result()
	require.Nil(t, err)
	require.Len(t, results, 3)

	assert.Equal(t, 1, results[0].Value)
	assert.Equal(t, 2, results[1].Value)
	assert.Equal(t, 3, results[2].Value)
}
