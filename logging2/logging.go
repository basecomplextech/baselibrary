package logging2

import "sync"

// Main is the main logger name.
const Main = "main"

// Logging is a logging service.
type Logging interface {
	// Logger is the main logger.
	Logger

	// Send logs the record.
	Send(rec Record)
}

// New returns a new logging service.
func New(config *Config) Logging {
	return newLogging(config)
}

// Default returns a new logging service with the default config.
func Default() Logging {
	config := DefaultConfig()
	return newLogging(config)
}

// internal

type logging struct {
	*logger // main logger

	level   Level
	writers []Writer

	mu      sync.RWMutex
	loggers map[string]*logger
}

func newLogging(config *Config) *logging {
	l := &logging{
		loggers: make(map[string]*logger),
	}

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
