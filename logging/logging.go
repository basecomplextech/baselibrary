package logging

import "os"

// Main is the root logger name.
const Main = "main"

var (
	Null   Logger = newLogger(Main, true, newNullWriter())
	Stdout Logger = newLogger(Main, true, newConsoleWriter(LevelDebug, true, os.Stdout))
	Stderr Logger = newLogger(Main, true, newConsoleWriter(LevelDebug, true, os.Stderr))
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
	Write(rec Record) error
}

// Init initializes and returns a new logging service.
func Init(config *Config) (Logging, error) {
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
	main   *logger // main logger
	writer Writer
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
		writer: newMultiWriter(writers...),
	}
	l.main = newLogger(Main, true, l)
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
	return l.writer.Enabled(level)
}

// Write writes a record.
func (l *logging) Write(rec Record) error {
	return l.writer.Write(rec)
}
