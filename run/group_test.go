package run

import (
	"testing"
	"time"
)

func TestGroup__should_run_stop_routine_group(t *testing.T) {
	g := NewGroup()

	s0 := make(chan struct{})
	s1 := make(chan struct{})

	r0 := Run(func(stop <-chan struct{}) error {
		<-stop
		close(s0)
		return nil
	})

	r1 := Run(func(stop <-chan struct{}) error {
		<-stop
		close(s1)
		return nil
	})

	g.Add(r0, r1)
	g.Stop()

	select {
	case <-g.Done():
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}

	select {
	case <-s0:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}

	select {
	case <-s1:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestGroup_Remove__should_remove_routine_from_group(t *testing.T) {
	g := NewGroup()

	s0 := make(chan struct{})
	s1 := make(chan struct{})

	r0 := Run(func(stop <-chan struct{}) error {
		<-stop
		close(s0)
		return nil
	})

	r1 := Run(func(stop <-chan struct{}) error {
		<-stop
		close(s1)
		return nil
	})
	defer r1.Stop()

	g.Add(r0, r1)
	g.Remove(r1)
	g.Stop()

	select {
	case <-s0:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}

	select {
	default:
	case <-s1:
		t.Fatal("unexpected routine stop")
	}
}
