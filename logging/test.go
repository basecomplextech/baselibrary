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
	writer := newConsoleWriter(level, true, os.Stdout)
	return newLogger("main", true /* root */, writer)
}

// TestLoggerLevel returns a new test logger with the specified level.
func TestLoggerLevel(t tests.T, level Level) Logger {
	writer := newConsoleWriter(level, true, os.Stdout)
	return newLogger("main", true /* root */, writer)
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
