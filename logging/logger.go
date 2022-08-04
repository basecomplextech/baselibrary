package logging

type Logger interface {
	// Name returns the logger name.
	Name() string

	// Log returns a record builder with info level.
	Log() RecordBuilder

	// Logger returns a child logger.
	Logger(name string) Logger

	// Records

	// Trace logs a trace record.
	Trace(msg string, keyValues ...any)

	// Debug logs a debug record.
	Debug(msg string, keyValues ...any)

	// Info logs an info record.
	Info(msg string, keyValues ...any)

	// Warn logs a warning record.
	Warn(msg string, keyValues ...any)

	// Error logs an error record.
	Error(msg string, keyValues ...any)

	// Fatal logs a fatal record.
	Fatal(msg string, keyValues ...any)

	// Level checks

	// Enabled returns true if a level is enabled.
	Enabled(level Level) bool

	// TraceEnabled returns true if trace level is enabled.
	TraceEnabled() bool

	// DebugEnabled returns true if debug level is enabled.
	DebugEnabled() bool

	// InfoEnabled returns true if info level is enabled.
	InfoEnabled() bool

	// WarnEnabled returns true if warn level is enabled.
	WarnEnabled() bool

	// ErrorEnabled returns true if error level is enabled.
	ErrorEnabled() bool

	// FatalEnabled returns true if fatal level is enabled.
	FatalEnabled() bool
}

// internal

var _ Logger = (*logger)(nil)

type logger struct {
	l    *logging
	name string
	main bool
}

func newLogger(l *logging, name string) *logger {
	return &logger{
		l:    l,
		name: name,
	}
}

// Name returns the logger name.
func (l *logger) Name() string {
	return l.name
}

// Log returns a record builder with info level.
func (l *logger) Log() RecordBuilder {
	return newRecordBuilder(l.l, l.name).Level(LevelInfo)
}

// Logger returns a child logger.
func (l *logger) Logger(name string) Logger {
	if !l.main {
		name = l.name + "." + name
	}
	return l.l.logger(name)
}

// Records

// Trace logs a trace record.
func (l *logger) Trace(msg string, keyValues ...any) {
	rec := Record{
		Level:   LevelTrace,
		Logger:  l.name,
		Message: msg,
		Fields:  NewFields(keyValues...),
	}
	l.l.send(rec)
}

// Debug logs a debug record.
func (l *logger) Debug(msg string, keyValues ...any) {
	rec := Record{
		Level:   LevelDebug,
		Logger:  l.name,
		Message: msg,
		Fields:  NewFields(keyValues...),
	}
	l.l.send(rec)
}

// Info logs an info record.
func (l *logger) Info(msg string, keyValues ...any) {
	rec := Record{
		Level:   LevelInfo,
		Logger:  l.name,
		Message: msg,
		Fields:  NewFields(keyValues...),
	}
	l.l.send(rec)
}

// Warn logs a warning record.
func (l *logger) Warn(msg string, keyValues ...any) {
	rec := Record{
		Level:   LevelWarn,
		Logger:  l.name,
		Message: msg,
		Fields:  NewFields(keyValues...),
	}
	l.l.send(rec)
}

// Error logs an error record.
func (l *logger) Error(msg string, keyValues ...any) {
	rec := Record{
		Level:   LevelError,
		Logger:  l.name,
		Message: msg,
		Fields:  NewFields(keyValues...),
	}
	l.l.send(rec)
}

// Fatal logs a fatal record.
func (l *logger) Fatal(msg string, keyValues ...any) {
	rec := Record{
		Level:   LevelFatal,
		Logger:  l.name,
		Message: msg,
		Fields:  NewFields(keyValues...),
	}
	l.l.send(rec)
}

// Level checks

// Enabled returns true if a level is enabled.
func (l *logger) Enabled(level Level) bool {
	return l.l.enabled(l.name, level)
}

// TraceEnabled returns true if trace level is enabled.
func (l *logger) TraceEnabled() bool {
	return l.l.enabled(l.name, LevelTrace)
}

// DebugEnabled returns true if debug level is enabled.
func (l *logger) DebugEnabled() bool {
	return l.l.enabled(l.name, LevelDebug)
}

// InfoEnabled returns true if info level is enabled.
func (l *logger) InfoEnabled() bool {
	return l.l.enabled(l.name, LevelInfo)
}

// WarnEnabled returns true if warn level is enabled.
func (l *logger) WarnEnabled() bool {
	return l.l.enabled(l.name, LevelWarn)
}

// ErrorEnabled returns true if error level is enabled.
func (l *logger) ErrorEnabled() bool {
	return l.l.enabled(l.name, LevelError)
}

// FatalEnabled returns true if fatal level is enabled.
func (l *logger) FatalEnabled() bool {
	return l.l.enabled(l.name, LevelFatal)
}
