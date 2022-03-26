package async

import (
	"errors"
	"testing"
	"time"

	"github.com/baseblck/library/try"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun__should_run_and_stop(t *testing.T) {
	r := Run(func(stop <-chan struct{}) error {
		return nil
	})
	r.Stop()

	select {
	case <-r.Wait():
		err := r.Err()
		if err != nil {
			t.Fatal(err)
		}

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestRun__should_run_and_return_error(t *testing.T) {
	err := errors.New("test")
	r := Run(func(stop <-chan struct{}) error {
		return err
	})
	r.Stop()

	select {
	case <-r.Wait():
		err1 := r.Err()
		assert.Equal(t, err, err1)

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestRun__should_run_and_recover(t *testing.T) {
	r := Run(func(stop <-chan struct{}) error {
		panic("test")
	})
	r.Stop()

	select {
	case <-r.Wait():
		err := r.Err()
		require.IsType(t, &try.Panic{}, err)
		assert.EqualValues(t, "test", err.(*try.Panic).E)

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}
