package run

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun__should_run_and_stop(t *testing.T) {
	th := Run(func(stop <-chan struct{}) error {
		return nil
	})
	th.Stop()

	select {
	case <-th.Done():
		err := th.Err()
		if err != nil {
			t.Fatal(err)
		}

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestRun__should_run_and_return_error(t *testing.T) {
	err := errors.New("test")
	th := Run(func(stop <-chan struct{}) error {
		return err
	})
	th.Stop()

	select {
	case <-th.Done():
		err1 := th.Err()
		assert.Equal(t, err, err1)

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestRun__should_run_and_recover(t *testing.T) {
	th := Run(func(stop <-chan struct{}) error {
		panic("test")
	})
	th.Stop()

	select {
	case <-th.Done():
		err := th.Err()
		require.IsType(t, &Panic{}, err)
		assert.EqualValues(t, "test", err.(*Panic).E)

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}
