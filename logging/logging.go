package logging

import (
	"os"
	"slices"
)

// Logging is a logging service.
type Logging interface {
	// Main returns the main logger.
	Main() Logger

	// Logger returns a logger with the given name or creates a new one.
	Logger(name string) Logger

	// Enabled returns true if logging is enabled for the given level.
	Enabled(level Level) bool

	// Write writes a record.
	Write(rec *Record) error
}

// Init initializes and returns a new logging service, uses the default config when config is nil.
func Init(config *Config) (Logging, error) {
	if config == nil {
		config = DefaultConfig()
	}
	return initLogging(config)
}

// Load loads a logging config from a JSON/YAML file if it exists and returns a new logging service.
func Load(path string) (Logging, error) {
	config, err := LoadConfig(path)
	switch {
	case os.IsNotExist(err):
		config = DefaultConfig()
	case err != nil:
		return nil, err
	}
	return initLogging(config)
}

// Default returns a new logging service with the default config.
func Default() Logging {
	config := DefaultConfig()

	l, err := initLogging(config)
	if err != nil {
		panic(err) // unreachable
	}
	return l
}

// internal

type logging struct {
	main    *logger
	writers []Writer
}

func initLogging(config *Config) (*logging, error) {
	var writers []Writer

	if config.Console != nil && config.Console.Enabled {
		console, err := initConsoleWriter(config.Console)
		if err != nil {
			return nil, err
		}
		writers = append(writers, console)
	}

	if config.File != nil && config.File.Enabled {
		file, err := initFileWriter(config.File)
		if err != nil {
			return nil, err
		}
		writers = append(writers, file)
	}

	return newLogging(writers...), nil
}

func newLogging(writers ...Writer) *logging {
	l := &logging{
		writers: slices.Clone(writers),
	}
	l.main = newLogger("main", true, l)
	return l
}

// Main returns the main logger.
func (l *logging) Main() Logger {
	return l.main
}

// Logger returns a logger with the given name or creates a new one.
func (l *logging) Logger(name string) Logger {
	return l.main.Logger(name)
}

// Enabled returns true if logging is enabled for the given level.
func (l *logging) Enabled(level Level) bool {
	for _, w := range l.writers {
		if w.Enabled(level) {
			return true
		}
	}
	return false
}

// Write writes a record.
func (l *logging) Write(rec *Record) error {
	for _, w := range l.writers {
		if err := w.Write(rec); err != nil {
			return err
		}
	}
	return nil
}
