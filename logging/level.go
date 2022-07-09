package logging

import (
	"fmt"
	"strings"
)

type Level int

const (
	LevelDisabled Level = iota
	LevelTrace
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var levelNames = map[Level]string{
	LevelDisabled: "",
	LevelTrace:    "TRACE",
	LevelDebug:    "DEBUG",
	LevelInfo:     "INFO",
	LevelWarn:     "WARN",
	LevelError:    "ERROR",
	LevelFatal:    "FATAL",
}

func (l Level) Name() string {
	return levelNames[l]
}

func ParseLevel(s string) (Level, error) {
	s = strings.ToUpper(s)

	for level, name := range levelNames {
		if name == s {
			return level, nil
		}
	}

	return LevelTrace, fmt.Errorf("Unsupported log level %q", s)
}
