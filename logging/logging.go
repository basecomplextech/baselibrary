package logging

import "sync"

// Main is the main logger name.
const Main = "main"

var (
	Null   Logger = newLogging(LevelDebug, nil).Main()
	Stdout Logger = newLogging(LevelDebug, newStdoutWriter(LevelTrace)).Main()
	Stderr Logger = newLogging(LevelDebug, newStderrWriter(LevelTrace)).Main()
)

// Logging is a logging service.
type Logging interface {
	// Main returns the main logger.
	Main() Logger

	// Send logs the record.
	Send(rec Record)

	// Logger returns a logger with the given name or creates a new one.
	Logger(name string) Logger
}

// New returns a new logging service.
func New(config *Config) (Logging, error) {
	return openLogging(config)
}

// Default returns a new logging service with the default config.
func Default() Logging {
	config := DefaultConfig()
	l, err := openLogging(config)
	if err != nil {
		panic(err) // unreachable
	}
	return l
}

// internal

type logging struct {
	main    *logger // main logger
	level   Level
	writers []Writer

	mu      sync.RWMutex
	loggers map[string]*logger
}

func openLogging(config *Config) (*logging, error) {
	level := config.Level
	var writers []Writer

	if config.Console != nil && config.Console.Enabled {
		console, err := newConsoleWriter(config.Console)
		if err != nil {
			return nil, err
		}
		writers = append(writers, console)
	}

	if config.File != nil && config.File.Enabled {
		file, err := newFileWriter(config.File)
		if err != nil {
			return nil, err
		}
		writers = append(writers, file)
	}

	return newLogging(level, writers...), nil
}

func newLogging(level Level, writers ...Writer) *logging {
	l := &logging{
		writers: make([]Writer, len(writers)),
		loggers: make(map[string]*logger),
	}
	copy(l.writers, writers)

	main := newLogger(l, Main)
	main.main = true

	l.main = main
	l.loggers[Main] = main
	return l
}

// Main returns the main logger.
func (l *logging) Main() Logger {
	return l.main
}

// Send logs the record.
func (l *logging) Send(rec Record) {
	l.send(rec)
}

// Logger returns a logger with the given name or creates a new one.
func (l *logging) Logger(name string) Logger {
	return l.logger(name)
}

// internal

func (l *logging) enabled(logger string, level Level) bool {
	return level >= l.level
}

func (l *logging) logger(name string) *logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	logger, ok := l.loggers[name]
	if ok {
		return logger
	}

	logger = newLogger(l, name)
	l.loggers[name] = logger
	return logger
}

func (l *logging) send(rec Record) {
	if rec.Level < l.level {
		return
	}

	for _, w := range l.writers {
		w.Write(rec)
	}
}
