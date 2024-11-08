// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package logging

import (
	"os"

	"github.com/basecomplextech/baselibrary/status"
)

var (
	Null   Logger = newLogger("null", true, newNullWriter())
	Stdout Logger = newLogger("main", true, newConsoleWriter(LevelDebug, true, os.Stdout))
	Stderr Logger = newLogger("main", true, newConsoleWriter(LevelDebug, true, os.Stderr))
)

type Logger interface {
	// Name returns the logger name.
	Name() string

	// Logger returns a child logger.
	Logger(name string) Logger

	// WithFields returns a chained logger with the default fields.
	WithFields(keyValuePairs ...any) Logger

	// Enabled returns true if a level is enabled.
	Enabled(level Level) bool

	// Record

	// Begin returns a record builder with the info level.
	Begin() RecordBuilder

	// Write sets the logger if abset, adds the default fields and writes the record.
	Write(rec *Record) error

	// Trace

	// Trace logs a trace message.
	Trace(msg string, keyValues ...any)

	// TraceStatus logs a trace message with a status and a stack trace.
	TraceStatus(msg string, st status.Status, keyValues ...any)

	// TraceOn returns true if trace level is enabled.
	TraceOn() bool

	// Debug

	// Debug logs a debug message.
	Debug(msg string, keyValues ...any)

	// DebugStatus logs a debug message with a status and a stack trace.
	DebugStatus(msg string, st status.Status, keyValues ...any)

	// DebugOn returns true if debug level is enabled.
	DebugOn() bool

	// Info

	// Info logs an info message.
	Info(msg string, keyValues ...any)

	// InfoStatus logs an info message with a status and a stack trace.
	InfoStatus(msg string, st status.Status, keyValues ...any)

	// InfoOn returns true if info level is enabled.
	InfoOn() bool

	// Notice

	// Notice logs a notice message.
	Notice(msg string, keyValues ...any)

	// NoticeStatus logs a notice message with a status and a stack trace.
	NoticeStatus(msg string, st status.Status, keyValues ...any)

	// NoticeOn return true if notice level is enabled.
	NoticeOn() bool

	// Warn

	// Warn logs a warning message.
	Warn(msg string, keyValues ...any)

	// WarnStatus logs a warning message with a status and a stack trace.
	WarnStatus(msg string, st status.Status, keyValues ...any)

	// WarnOn returns true if warn level is enabled.
	WarnOn() bool

	// Error

	// Error logs an error message.
	Error(msg string, keyValues ...any)

	// ErrorStatus logs an error message with a status and a stack trace.
	ErrorStatus(msg string, st status.Status, keyValues ...any)

	// ErrorOn returns true if error level is enabled.
	ErrorOn() bool

	// Fatal

	// Fatal logs a fatal mesage.
	Fatal(msg string, keyValues ...any)

	// FatalStatus logs a fatal message with a status and a stack trace.
	FatalStatus(msg string, st status.Status, keyValues ...any)

	// FatalOn returns true if fatal level is enabled.
	FatalOn() bool
}

// internal

var _ Logger = (*logger)(nil)

type logger struct {
	name   string
	root   bool
	fields []Field

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

// Logger returns a child logger.
func (l *logger) Logger(name string) Logger {
	if !l.root {
		name = l.name + "." + name
	}

	return newLogger(name, false, l.w)
}

// WithFields returns a chained logger with the default fields.
func (l *logger) WithFields(keyValuePairs ...any) Logger {
	if len(keyValuePairs) == 0 {
		return l
	}

	l1 := &logger{
		name:   l.name,
		root:   l.root,
		fields: NewFields(keyValuePairs...),

		w: l,
	}
	return l1
}

// Enabled returns true if a level is enabled.
func (l *logger) Enabled(level Level) bool {
	return l.w.Enabled(level)
}

// Record

// Begin returns a record builder with the info level.
func (l *logger) Begin() RecordBuilder {
	r := NewRecord(l.name, LevelInfo)
	return newRecordBuilder(l, r)
}

// Write sets the logger if abset, adds the default fields and writes the record.
func (l *logger) Write(rec *Record) error {
	if rec.Logger == "" {
		rec.Logger = l.name
	}

	if l.fields != nil {
		rec.Fields = append(rec.Fields, l.fields...)
	}

	return l.w.Write(rec)
}

// Trace

// Trace logs a trace message.
func (l *logger) Trace(msg string, keyValues ...any) {
	rec := NewRecord(l.name, LevelTrace).
		WithMessage(msg).
		WithFields(keyValues...)
	l.Write(rec)
}

// TraceStatus logs a trace message with a status and a stack trace.
func (l *logger) TraceStatus(msg string, st status.Status, keyValues ...any) {
	rec := NewRecord(l.name, LevelTrace).
		WithMessage(msg).
		WithFields(keyValues...).
		WithStatus(st)
	l.Write(rec)
}

// TraceOn returns true if trace level is enabled.
func (l *logger) TraceOn() bool {
	return l.w.Enabled(LevelTrace)
}

// Debug

// Debug logs a debug message.
func (l *logger) Debug(msg string, keyValues ...any) {
	rec := NewRecord(l.name, LevelDebug).
		WithMessage(msg).
		WithFields(keyValues...)
	l.Write(rec)
}

// DebugStatus logs a debug message with a status and a stack trace.
func (l *logger) DebugStatus(msg string, st status.Status, keyValues ...any) {
	rec := NewRecord(l.name, LevelDebug).
		WithMessage(msg).
		WithFields(keyValues...).
		WithStatus(st)
	l.Write(rec)
}

// DebugOn returns true if debug level is enabled.
func (l *logger) DebugOn() bool {
	return l.w.Enabled(LevelDebug)
}

// Info

// Info logs an info message.
func (l *logger) Info(msg string, keyValues ...any) {
	rec := NewRecord(l.name, LevelInfo).
		WithMessage(msg).
		WithFields(keyValues...)
	l.Write(rec)
}

// InfoStatus logs an info message with a status and a stack trace.
func (l *logger) InfoStatus(msg string, st status.Status, keyValues ...any) {
	rec := NewRecord(l.name, LevelInfo).
		WithMessage(msg).
		WithFields(keyValues...).
		WithStatus(st)
	l.Write(rec)
}

// InfoOn returns true if info level is enabled.
func (l *logger) InfoOn() bool {
	return l.w.Enabled(LevelInfo)
}

// Notice

// Notice logs a notice message.
func (l *logger) Notice(msg string, keyValues ...any) {
	rec := NewRecord(l.name, LevelNotice).
		WithMessage(msg).
		WithFields(keyValues...)
	l.Write(rec)
}

// NoticeStatus logs a notice message with a status and a stack trace.
func (l *logger) NoticeStatus(msg string, st status.Status, keyValues ...any) {
	rec := NewRecord(l.name, LevelNotice).
		WithMessage(msg).
		WithFields(keyValues...).
		WithStatus(st)
	l.Write(rec)
}

// NoticeOn return true if notice level is enabled.
func (l *logger) NoticeOn() bool {
	return l.w.Enabled(LevelNotice)
}

// Warn

// Warn logs a warning message.
func (l *logger) Warn(msg string, keyValues ...any) {
	rec := NewRecord(l.name, LevelWarn).
		WithMessage(msg).
		WithFields(keyValues...)
	l.Write(rec)
}

// WarnStatus logs a warning message with a status and a stack trace.
func (l *logger) WarnStatus(msg string, st status.Status, keyValues ...any) {
	rec := NewRecord(l.name, LevelWarn).
		WithMessage(msg).
		WithFields(keyValues...).
		WithStatus(st)
	l.Write(rec)
}

// WarnOn returns true if warn level is enabled.
func (l *logger) WarnOn() bool {
	return l.w.Enabled(LevelWarn)
}

// Error

// Error logs an error message.
func (l *logger) Error(msg string, keyValues ...any) {
	rec := NewRecord(l.name, LevelError).
		WithMessage(msg).
		WithFields(keyValues...)
	l.Write(rec)
}

// ErrorStatus logs an error message with a status and a stack trace.
func (l *logger) ErrorStatus(msg string, st status.Status, keyValues ...any) {
	rec := NewRecord(l.name, LevelError).
		WithMessage(msg).
		WithFields(keyValues...).
		WithStatus(st)
	l.Write(rec)
}

// ErrorOn returns true if error level is enabled.
func (l *logger) ErrorOn() bool {
	return l.w.Enabled(LevelError)
}

// Fatal

// Fatal logs a fatal mesage.
func (l *logger) Fatal(msg string, keyValues ...any) {
	rec := NewRecord(l.name, LevelFatal).
		WithMessage(msg).
		WithFields(keyValues...)
	l.Write(rec)
}

// FatalStatus logs a fatal message with a status and a stack trace.
func (l *logger) FatalStatus(msg string, st status.Status, keyValues ...any) {
	rec := NewRecord(l.name, LevelFatal).
		WithMessage(msg).
		WithFields(keyValues...).
		WithStatus(st)
	l.Write(rec)
}

// FatalOn returns true if fatal level is enabled.
func (l *logger) FatalOn() bool {
	return l.w.Enabled(LevelFatal)
}
