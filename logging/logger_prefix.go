// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package logging

import (
	"sync/atomic"
	"unsafe"

	"github.com/basecomplextech/baselibrary/status"
)

// PrefixLogger is a logger that adds a prefix to messages, including messages from child loggers.
type PrefixLogger interface {
	Logger

	// SetPrefix sets the logger prefix.
	SetPrefix(s string)

	// ClearPrefix clears the logger prefix.
	ClearPrefix()
}

// NewPrefixLogger returns a new prefix logger.
func NewPrefixLogger(logger Logger) PrefixLogger {
	prefix := newPrefix()
	return newPrefixLogger(logger, prefix)
}

type prefixLogger struct {
	logger Logger
	prefix *prefix
}

func newPrefixLogger(logger Logger, prefix *prefix) *prefixLogger {
	return &prefixLogger{
		logger: logger,
		prefix: prefix,
	}
}

// SetPrefix sets the logger prefix.
func (l *prefixLogger) SetPrefix(s string) {
	l.prefix.store(s)
}

// ClearPrefix clears the logger prefix.
func (l *prefixLogger) ClearPrefix() {
	l.prefix.clear()
}

// Build

// WithFields returns a chained logger with the default fields.
func (l *prefixLogger) WithFields(keyValuePairs ...any) Logger {
	child := l.logger.WithFields(keyValuePairs...)
	return newPrefixLogger(child, l.prefix)
}

// Logger

// Name returns the logger name.
func (l *prefixLogger) Name() string {
	return l.logger.Name()
}

// Logger returns a child logger.
func (l *prefixLogger) Logger(name string) Logger {
	child := l.logger.Logger(name)
	return newPrefixLogger(child, l.prefix)
}

// Enabled returns true if a level is enabled.
func (l *prefixLogger) Enabled(level Level) bool {
	return l.logger.Enabled(level)
}

// Record

// Begin returns a record builder with info level.
func (l *prefixLogger) Begin() RecordBuilder {
	return l.logger.Begin()
}

// Write writes a record.
func (l *prefixLogger) Write(rec *Record) error {
	if p := l.prefix.load(); p != "" {
		rec.Message = p + rec.Message
	}

	return l.logger.Write(rec)
}

// Trace

// Trace logs a trace record.
func (l *prefixLogger) Trace(msg string, keyValues ...any) {
	if p := l.prefix.load(); p != "" {
		msg = p + msg
	}

	l.logger.Trace(msg, keyValues...)
}

// TraceStatus logs a trace message with a status and a stack trace.
func (l *prefixLogger) TraceStatus(msg string, st status.Status, keyValues ...any) {
	if p := l.prefix.load(); p != "" {
		msg = p + msg
	}

	l.logger.TraceStatus(msg, st, keyValues...)
}

// TraceOn returns true if trace level is enabled.
func (l *prefixLogger) TraceOn() bool {
	return l.logger.TraceOn()
}

// Debug

// Debug logs a debug record.
func (l *prefixLogger) Debug(msg string, keyValues ...any) {
	if p := l.prefix.load(); p != "" {
		msg = p + msg
	}

	l.logger.Debug(msg, keyValues...)
}

// DebugStatus logs a debug message with a status and a stack trace.
func (l *prefixLogger) DebugStatus(msg string, st status.Status, keyValues ...any) {
	if p := l.prefix.load(); p != "" {
		msg = p + msg
	}

	l.logger.DebugStatus(msg, st, keyValues...)
}

// DebugOn returns true if debug level is enabled.
func (l *prefixLogger) DebugOn() bool {
	return l.logger.DebugOn()
}

// Info

// Info logs an info record.
func (l *prefixLogger) Info(msg string, keyValues ...any) {
	if p := l.prefix.load(); p != "" {
		msg = p + msg
	}

	l.logger.Info(msg, keyValues...)
}

// InfoStatus logs an info message with a status and a stack trace.
func (l *prefixLogger) InfoStatus(msg string, st status.Status, keyValues ...any) {
	if p := l.prefix.load(); p != "" {
		msg = p + msg
	}

	l.logger.InfoStatus(msg, st, keyValues...)
}

// InfoOn returns true if info level is enabled.
func (l *prefixLogger) InfoOn() bool {
	return l.logger.InfoOn()
}

// Notice

// Notice logs a notice record.
func (l *prefixLogger) Notice(msg string, keyValues ...any) {
	if p := l.prefix.load(); p != "" {
		msg = p + msg
	}

	l.logger.Notice(msg, keyValues...)
}

// NoticeStatus logs a notice message with a status and a stack trace.
func (l *prefixLogger) NoticeStatus(msg string, st status.Status, keyValues ...any) {
	if p := l.prefix.load(); p != "" {
		msg = p + msg
	}

	l.logger.NoticeStatus(msg, st, keyValues...)
}

// NoticeOn return true if notice level is enabled.
func (l *prefixLogger) NoticeOn() bool {
	return l.logger.NoticeOn()
}

// Warn

// Warn logs a warning record.
func (l *prefixLogger) Warn(msg string, keyValues ...any) {
	if p := l.prefix.load(); p != "" {
		msg = p + msg
	}

	l.logger.Warn(msg, keyValues...)
}

// WarnStatus logs a warning message with a status and a stack trace.
func (l *prefixLogger) WarnStatus(msg string, st status.Status, keyValues ...any) {
	if p := l.prefix.load(); p != "" {
		msg = p + msg
	}

	l.logger.WarnStatus(msg, st, keyValues...)
}

// WarnOn returns true if warn level is enabled.
func (l *prefixLogger) WarnOn() bool {
	return l.logger.WarnOn()
}

// Error

// Error logs an error record.
func (l *prefixLogger) Error(msg string, keyValues ...any) {
	if p := l.prefix.load(); p != "" {
		msg = p + msg
	}

	l.logger.Error(msg, keyValues...)
}

// ErrorStatus logs an error message with a status and a stack trace.
func (l *prefixLogger) ErrorStatus(msg string, st status.Status, keyValues ...any) {
	if p := l.prefix.load(); p != "" {
		msg = p + msg
	}

	l.logger.ErrorStatus(msg, st, keyValues...)
}

// ErrorOn returns true if error level is enabled.
func (l *prefixLogger) ErrorOn() bool {
	return l.logger.ErrorOn()
}

// Fatal

// Fatal logs a fatal record.
func (l *prefixLogger) Fatal(msg string, keyValues ...any) {
	if p := l.prefix.load(); p != "" {
		msg = p + msg
	}

	l.logger.Fatal(msg, keyValues...)
}

// FatalStatus logs a fatal message with a status and a stack trace.
func (l *prefixLogger) FatalStatus(msg string, st status.Status, keyValues ...any) {
	if p := l.prefix.load(); p != "" {
		msg = p + msg
	}

	l.logger.FatalStatus(msg, st, keyValues...)
}

// FatalOn returns true if fatal level is enabled.
func (l *prefixLogger) FatalOn() bool {
	return l.logger.FatalOn()
}

// private

type prefix struct {
	s *string
}

func newPrefix() *prefix {
	return &prefix{}
}

func (p *prefix) clear() {
	ptr := (*unsafe.Pointer)(unsafe.Pointer(&p.s))
	atomic.StorePointer(ptr, nil)
}

func (p *prefix) load() string {
	ptr := (*unsafe.Pointer)(unsafe.Pointer(&p.s))
	sptr := atomic.LoadPointer(ptr)
	if sptr == nil {
		return ""
	}

	s := (*string)(sptr)
	return *s
}

func (p *prefix) store(s string) {
	sptr := &s
	ptr := (*unsafe.Pointer)(unsafe.Pointer(&p.s))
	atomic.StorePointer(ptr, unsafe.Pointer(sptr))
}
