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

// Begin returns a record builder with info level.
func (l *prefixLogger) Begin() RecordBuilder {
	return l.logger.Begin()
}

// Logger returns a child logger.
func (l *prefixLogger) Logger(name string) Logger {
	child := l.logger.Logger(name)
	return newPrefixLogger(child, l.prefix)
}

// Write writes a record.
func (l *prefixLogger) Write(rec *Record) error {
	if p := l.prefix.load(); p != "" {
		rec.Message = p + rec.Message
	}

	return l.logger.Write(rec)
}

// Records

// Trace logs a trace record.
func (l *prefixLogger) Trace(msg string, keyValues ...any) {
	if p := l.prefix.load(); p != "" {
		msg = p + msg
	}

	l.logger.Trace(msg, keyValues...)
}

// Debug logs a debug record.
func (l *prefixLogger) Debug(msg string, keyValues ...any) {
	if p := l.prefix.load(); p != "" {
		msg = p + msg
	}

	l.logger.Debug(msg, keyValues...)
}

// Info logs an info record.
func (l *prefixLogger) Info(msg string, keyValues ...any) {
	if p := l.prefix.load(); p != "" {
		msg = p + msg
	}

	l.logger.Info(msg, keyValues...)
}

// Notice logs a notice record.
func (l *prefixLogger) Notice(msg string, keyValues ...any) {
	if p := l.prefix.load(); p != "" {
		msg = p + msg
	}

	l.logger.Notice(msg, keyValues...)
}

// Warn logs a warning record.
func (l *prefixLogger) Warn(msg string, keyValues ...any) {
	if p := l.prefix.load(); p != "" {
		msg = p + msg
	}

	l.logger.Warn(msg, keyValues...)
}

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

// Level checks

// Enabled returns true if a level is enabled.
func (l *prefixLogger) Enabled(level Level) bool {
	return l.logger.Enabled(level)
}

// TraceEnabled returns true if trace level is enabled.
func (l *prefixLogger) TraceEnabled() bool {
	return l.logger.TraceEnabled()
}

// DebugEnabled returns true if debug level is enabled.
func (l *prefixLogger) DebugEnabled() bool {
	return l.logger.DebugEnabled()
}

// InfoEnabled returns true if info level is enabled.
func (l *prefixLogger) InfoEnabled() bool {
	return l.logger.InfoEnabled()
}

// NoticeEnabled return true if notice level is enabled.
func (l *prefixLogger) NoticeEnabled() bool {
	return l.logger.NoticeEnabled()
}

// WarnEnabled returns true if warn level is enabled.
func (l *prefixLogger) WarnEnabled() bool {
	return l.logger.WarnEnabled()
}

// ErrorEnabled returns true if error level is enabled.
func (l *prefixLogger) ErrorEnabled() bool {
	return l.logger.ErrorEnabled()
}

// FatalEnabled returns true if fatal level is enabled.
func (l *prefixLogger) FatalEnabled() bool {
	return l.logger.FatalEnabled()
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
