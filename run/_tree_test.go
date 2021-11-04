package run

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var testMu sync.Mutex

func a(stop <-chan struct{}) error {
	println("a started")

	a := Run(a_a)
	b := Run(a_b)
	c := Run(a_c)
	g := Group(a, b, c)

	<-stop
	println("a stopping...")

	<-g.Stop()
	println("a Done")
	return nil
}

func a_a(stop <-chan struct{}) error {
	println("  a started")

	a := Run(a_a_a)
	b := Run(a_a_b)
	g := Group(a, b)

	<-stop
	println("  a stopping...")

	<-g.Stop()
	println("  a Done")
	return nil
}

func a_a_a(stop <-chan struct{}) error {
	println("    a started")
	<-stop
	println("    a Done")
	return nil
}

func a_a_b(stop <-chan struct{}) error {
	println("    b started")

	c := Run(a_a_b_c)
	g := Group(c)

	<-stop
	println("    b stopping...")

	<-g.Stop()
	println("    b Done")
	return nil
}

func a_a_b_c(stop <-chan struct{}) error {
	println("      c started")
	<-stop
	println("      c Done")
	return nil
}

func a_b(stop <-chan struct{}) error {
	println("  b started")
	<-stop
	println("  b Done")
	return nil
}

func a_c(stop <-chan struct{}) error {
	println("  c started")
	<-stop
	println("  c Done")
	return nil
}

func println(f string, v ...interface{}) {
	testMu.Lock()
	defer testMu.Unlock()

	fmt.Printf(f+"\n", v...)
}

func TestTree(t *testing.T) {
	r := Run(a)

	time.Sleep(time.Second)
	println("")

	r.Stop()

	select {
	case <-r.Done():
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}

	println("")
	// t.Fatal()
}
