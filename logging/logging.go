package logging

import (
	"os"
	"sync"
)

var (
	Stdout Logger = stdLogger("stdout", os.Stdout)
	Stderr Logger = stdLogger("stderr", os.Stderr)
	Null   Logger = Stdout
)

type Logging interface {
	Logger(name string) Logger
}

type Logger interface {
	// Logger returns a child logger.
	Logger(name string) Logger

	// Log logs a record.
	Log(record Record)

	// Enabled returns true if any write or parent logger logs a level.
	Enabled(level Level) bool

	// Utility methods

	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
	TraceEnabled() bool

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	DebugEnabled() bool

	Info(args ...interface{})
	Infof(format string, args ...interface{})
	InfoEnabled() bool

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	WarnEnabled() bool

	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	ErrorEnabled() bool

	// Panic logs an error and includes a stack trace.
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
}

type Writer interface {
	Level() Level
	Write(record Record)
}

type Record struct {
	Logger string
	Level  Level
	Format string
	Args   []interface{}
	Stack  []byte // Stacktrace.
}

func New(config *Config) Logging {
	main := openLogger(nil, &LoggerConfig{
		Console: config.Console,
		File:    config.File,
	})

	ll := map[string]*logger{}
	for _, lc := range config.Loggers {
		l := openLogger(main, lc)
		ll[l.name] = l
	}

	return &logging{
		main:    main,
		loggers: ll,
	}
}

type logging struct {
	mu      sync.Mutex
	main    *logger
	loggers map[string]*logger
}

func (l *logging) Logger(name string) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	logger, ok := l.loggers[name]
	if ok {
		return logger
	}

	logger = newLogger(l.main, name)
	l.loggers[name] = logger
	return logger
}
