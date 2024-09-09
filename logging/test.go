// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package logging

import (
	"os"

	"github.com/basecomplextech/baselibrary/tests"
)

// TEST_LOG specifies an env variable that can be used to override the log level.
const TEST_LOG = "TEST_LOG"

// TestLevel is the default log level for tests, you can change it in your modules for tests.
var TestLevel = LevelError

// Test returns a new test logging service.
func Test(t tests.T) Logging {
	level := TestLevelEnv()
	writer := newConsoleWriter(level, true, os.Stdout)
	return newLogging(writer)
}

// TestLogger returns a new test logger.
func TestLogger(t tests.T) Logger {
	level := TestLevelEnv()
	return TestLoggerLevel(t, level)
}

// TestLoggerDebug returns a new test logger with the debug level.
func TestLoggerDebug(t tests.T) Logger {
	return TestLoggerLevel(t, LevelDebug)
}

// TestLoggerInfo returns a new test logger with the info level.
func TestLoggerInfo(t tests.T) Logger {
	return TestLoggerLevel(t, LevelInfo)
}

// TestLoggerLevel returns a new test logger with the specified level.
func TestLoggerLevel(t tests.T, level Level) Logger {
	writer := newConsoleWriter(level, true, os.Stdout)
	return newLogger("test", true /* root */, writer)
}

// TestLevelEnv returns a log level from the env variable TEST_LOG or the default test level.
func TestLevelEnv() Level {
	if v := os.Getenv(TEST_LOG); v != "" {
		lv := LevelFromString(v)
		if lv != LevelUndefined {
			return lv
		}
	}
	return TestLevel
}
