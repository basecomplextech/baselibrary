// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package logging

import (
	"testing"
)

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

	l1 := l.WithFields("key0", "value0", "key1", "value1")
	l1.Error("With default fields", "key2", "value2")
	l1.Error("With default fields", "key2", "another value2")
}
