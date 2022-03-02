package logging

type nullLogger struct{}

// Logger returns a child logger.
func (l nullLogger) Logger(name string) Logger {
	return l
}

// Log logs a record.
func (nullLogger) Log(record Record) {}

// Enabled returns true if any write or parent logger logs a level.
func (nullLogger) Enabled(level Level) bool {
	return false
}

// Utility methods

func (nullLogger) Trace(args ...interface{})                 {}
func (nullLogger) Tracef(format string, args ...interface{}) {}
func (nullLogger) TraceEnabled() bool                        { return false }

func (nullLogger) Debug(args ...interface{})                 {}
func (nullLogger) Debugf(format string, args ...interface{}) {}
func (nullLogger) DebugEnabled() bool                        { return false }

func (nullLogger) Info(args ...interface{})                 {}
func (nullLogger) Infof(format string, args ...interface{}) {}
func (nullLogger) InfoEnabled() bool                        { return false }

func (nullLogger) Warn(args ...interface{})                 {}
func (nullLogger) Warnf(format string, args ...interface{}) {}
func (nullLogger) WarnEnabled() bool                        { return false }

func (nullLogger) Error(args ...interface{})                 {}
func (nullLogger) Errorf(format string, args ...interface{}) {}
func (nullLogger) ErrorEnabled() bool                        { return false }

// Panic logs an error and includes a stack trace.
func (nullLogger) Panic(args ...interface{})                 {}
func (nullLogger) Panicf(format string, args ...interface{}) {}
