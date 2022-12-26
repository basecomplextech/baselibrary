package logging

import (
	"log"
	"os"
)

func newStdoutWriter(level Level) *writer {
	logger := log.New(os.Stdout, "", 0)
	formatter := newTextFormatterColor(false, true)
	return newWriter(level, logger, formatter)
}

func newStderrWriter(level Level) *writer {
	logger := log.New(os.Stderr, "", 0)
	formatter := newTextFormatterColor(false, true)
	return newWriter(level, logger, formatter)
}

func newConsoleWriter(config *ConsoleConfig) (*writer, error) {
	if !config.Enabled {
		return nil, nil
	}
	level := config.Level
	color := config.Color

	logger := log.New(os.Stdout, "", 0)
	formatter := newTextFormatterColor(color, true)
	return newWriter(level, logger, formatter), nil
}
