package logging2

import (
	"log"
	"os"
)

func newStdoutWriter(level Level) *writer {
	logger := log.New(os.Stdout, "", lflags)
	formatter := newTextFormatter()
	return newWriter(level, logger, formatter)
}

func newStderrWriter(level Level) *writer {
	logger := log.New(os.Stderr, "", lflags)
	formatter := newTextFormatter()
	return newWriter(level, logger, formatter)
}

func newConsoleWriter(config *ConsoleConfig) *writer {
	if !config.Enabled {
		return nil
	}

	level := config.Level
	logger := log.New(os.Stdout, "", lflags)
	formatter := newTextFormatter()
	return newWriter(level, logger, formatter)
}
