package logging

import "sync"

// Main is the main logger name.
const Main = "main"

var (
	Null   Logger = newLogging(LevelDebug, nil)
	Stdout Logger = newLogging(LevelDebug, []Writer{newStdoutWriter(LevelTrace)})
	Stderr Logger = newLogging(LevelDebug, []Writer{newStderrWriter(LevelTrace)})
)

// Logging is a logging service.
type Logging interface {
	// Logger is the main logger.
	Logger

	// Send logs the record.
	Send(rec Record)
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
	*logger // main logger

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

	return newLogging(level, writers), nil
}

func newLogging(level Level, writers []Writer) *logging {
	l := &logging{
		writers: make([]Writer, len(writers)),
		loggers: make(map[string]*logger),
	}
	copy(l.writers, writers)

	main := newLogger(l, Main)
	main.main = true

	l.logger = main
	l.loggers[Main] = main
	return l
}

// Send logs the record.
func (l *logging) Send(rec Record) {
	l.send(rec)
}

// internal

func (l *logging) child(name string) *logger {
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

func (l *logging) enabled(logger string, level Level) bool {
	return level >= l.level
}

func (l *logging) send(rec Record) {
	if rec.Level < l.level {
		return
	}

	for _, w := range l.writers {
		w.Write(rec)
	}
}
