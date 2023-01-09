package logging

import "testing"

func TestLogging(t *testing.T) {
	l := Stdout
	l.Trace("Trace message", "key", "value")
	l.Debug("Debug message", "key", "value")
	l.Info("Info message", "key", "value", "key1", 1234)
	l.Notice("Notice message", "key", "value", "key1", 1234)
	l.Warn("Warn message", "key", "value", "key1", 1234)
	l.Error("Error message", "key", "value", "key1", 1234)
	l.Fatal("Fatal message")

	l.Begin().
		Infof("Hello, %v", "world").
		Fields("key", "value", "key1", 1234).
		Send()
}
