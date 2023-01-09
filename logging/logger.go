package logging

import "time"

type Logger interface {
	// Name returns the logger name.
	Name() string

	// Log returns a record builder with the info level.
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

	// Notice logs a notice record.
	Notice(msg string, keyValues ...any)

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

	// NoticeEnabled return true if notice level is enabled.
	NoticeEnabled() bool

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
	name string
	root bool

	w Writer
}

func newLogger(name string, root bool, writer Writer) *logger {
	return &logger{
		name: name,
		root: root,

		w: writer,
	}
}

// Name returns the logger name.
func (l *logger) Name() string {
	return l.name
}

// Log returns a record builder with the info level.
func (l *logger) Log() RecordBuilder {
	return newRecordBuilder(l.w, l.name).Level(LevelInfo)
}

// Logger returns a child logger.
func (l *logger) Logger(name string) Logger {
	if !l.root {
		name = l.name + "." + name
	}

	return newLogger(name, false, l.w)
}

// Records

// Trace logs a trace record.
func (l *logger) Trace(msg string, keyValues ...any) {
	rec := Record{
		Time:    time.Now(),
		Level:   LevelTrace,
		Logger:  l.name,
		Message: msg,
		Fields:  NewFields(keyValues...),
	}
	l.w.Write(rec)
}

// Debug logs a debug record.
func (l *logger) Debug(msg string, keyValues ...any) {
	rec := Record{
		Time:    time.Now(),
		Level:   LevelDebug,
		Logger:  l.name,
		Message: msg,
		Fields:  NewFields(keyValues...),
	}
	l.w.Write(rec)
}

// Info logs an info record.
func (l *logger) Info(msg string, keyValues ...any) {
	rec := Record{
		Time:    time.Now(),
		Level:   LevelInfo,
		Logger:  l.name,
		Message: msg,
		Fields:  NewFields(keyValues...),
	}
	l.w.Write(rec)
}

// Notice logs a notice record.
func (l *logger) Notice(msg string, keyValues ...any) {
	rec := Record{
		Time:    time.Now(),
		Level:   LevelNotice,
		Logger:  l.name,
		Message: msg,
		Fields:  NewFields(keyValues...),
	}
	l.w.Write(rec)
}

// Warn logs a warning record.
func (l *logger) Warn(msg string, keyValues ...any) {
	rec := Record{
		Time:    time.Now(),
		Level:   LevelWarn,
		Logger:  l.name,
		Message: msg,
		Fields:  NewFields(keyValues...),
	}
	l.w.Write(rec)
}

// Error logs an error record.
func (l *logger) Error(msg string, keyValues ...any) {
	rec := Record{
		Time:    time.Now(),
		Level:   LevelError,
		Logger:  l.name,
		Message: msg,
		Fields:  NewFields(keyValues...),
	}
	l.w.Write(rec)
}

// Fatal logs a fatal record.
func (l *logger) Fatal(msg string, keyValues ...any) {
	rec := Record{
		Time:    time.Now(),
		Level:   LevelFatal,
		Logger:  l.name,
		Message: msg,
		Fields:  NewFields(keyValues...),
	}
	l.w.Write(rec)
}

// Level checks

// Enabled returns true if a level is enabled.
func (l *logger) Enabled(level Level) bool {
	return l.w.Enabled(level)
}

// TraceEnabled returns true if trace level is enabled.
func (l *logger) TraceEnabled() bool {
	return l.w.Enabled(LevelTrace)
}

// DebugEnabled returns true if debug level is enabled.
func (l *logger) DebugEnabled() bool {
	return l.w.Enabled(LevelDebug)
}

// InfoEnabled returns true if info level is enabled.
func (l *logger) InfoEnabled() bool {
	return l.w.Enabled(LevelInfo)
}

// NoticeEnabled return true if notice level is enabled.
func (l *logger) NoticeEnabled() bool {
	return l.w.Enabled(LevelNotice)
}

// WarnEnabled returns true if warn level is enabled.
func (l *logger) WarnEnabled() bool {
	return l.w.Enabled(LevelWarn)
}

// ErrorEnabled returns true if error level is enabled.
func (l *logger) ErrorEnabled() bool {
	return l.w.Enabled(LevelError)
}

// FatalEnabled returns true if fatal level is enabled.
func (l *logger) FatalEnabled() bool {
	return l.w.Enabled(LevelFatal)
}
