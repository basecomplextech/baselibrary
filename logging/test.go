package logging

import (
	"os"

	"github.com/epochtimeout/baselibrary/build"
	"github.com/epochtimeout/baselibrary/tests"
)

// TEST_LOG specifies an env variable that can be used to override the log level.
const TEST_LOG = "TEST_LOG"

// TestLevel is the default log level for tests, you can change it in your modules for tests.
var TestLevel = LevelDebug

// Test returns a new test logging service.
func Test(t tests.T) Logging {
	level := TestLevelEnv()
	writer := newStdoutWriter(level)
	return newLogging(level, writer)
}

// TestLogger returns a new test logger.
func TestLogger(t tests.T) Logger {
	return Test(t).Main()
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

// TestModule is a build module that provides test logging and logger.
func TestModule(t tests.T) build.Module {
	return func(x *build.X) {
		build.Add(x, func() Logging {
			return Test(t)
		})
		build.Add(x, func() Logger {
			l := build.Get[Logging](x)
			return l.Main()
		})
	}
}
