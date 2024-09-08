// Copyright 2022 Ivan Korobkov. All rights reserved.

package logging

import (
	"github.com/basecomplextech/baselibrary/status"
)

type RecordBuilder interface {
	// Build builds and returns the record, but does not send it.
	Build() *Record

	// Send sends the record.
	Send()

	// Attrs

	// Level sets the record level.
	Level(lv Level) RecordBuilder

	// Message sets the record message.
	Message(msg string) RecordBuilder

	// Messagef formats and sets the record message.
	Messagef(msg string, args ...any) RecordBuilder

	// Stack adds a stack trace.
	Stack(stack []byte) RecordBuilder

	// Status adds a status with an optional stack trace.
	Status(st status.Status) RecordBuilder

	// Fields

	// Field adds a field to the record.
	Field(key string, value any) RecordBuilder

	// Fieldf formats a field value and adds a field to the record.
	Fieldf(key string, format string, a ...any) RecordBuilder

	// Fields adds fields to the record.
	Fields(keyValuePairs ...any) RecordBuilder

	// Utility

	// Trace sets the level to trace and adds the message.
	Trace(msg string) RecordBuilder

	// Tracef sets the level to trace and formats the message.
	Tracef(format string, args ...any) RecordBuilder

	// Debug sets the level to debug and adds the message.
	Debug(msg string) RecordBuilder

	// Debugf sets the level to debug and formats the message.
	Debugf(format string, args ...any) RecordBuilder

	// Info sets the level to info and adds the message.
	Info(msg string) RecordBuilder

	// Infof sets the level to info and formats the message.
	Infof(format string, args ...any) RecordBuilder

	// Notice sets the level to notice and adds the message.
	Notice(msg string) RecordBuilder

	// Noticef sets the level to notice and formats the message.
	Noticef(format string, args ...any) RecordBuilder

	// Warn sets the level to warn and adds the message.
	Warn(msg string) RecordBuilder

	// Warnf sets the level to warn and formats the message.
	Warnf(format string, args ...any) RecordBuilder

	// Error sets the level to error and adds the message.
	Error(msg string) RecordBuilder

	// Errorf sets the level to error and formats the message.
	Errorf(format string, args ...any) RecordBuilder

	// Fatal sets the level to fatal and adds the message.
	Fatal(msg string) RecordBuilder

	// Fatalf sets the level to fatal and formats the message.
	Fatalf(format string, args ...any) RecordBuilder
}

// internal

var _ RecordBuilder = (*recordBuilder)(nil)

type recordBuilder struct {
	w Writer
	r *Record
}

func newRecordBuilder(w Writer, r *Record) *recordBuilder {
	return &recordBuilder{
		w: w,
		r: r,
	}
}

// Build builds and returns the record, but does not send it.
func (b *recordBuilder) Build() *Record {
	return b.r
}

// Send sends the record.
func (b *recordBuilder) Send() {
	b.w.Write(b.r)
}

// Attrs

// Level sets the record level.
func (b *recordBuilder) Level(lv Level) RecordBuilder {
	b.r.WithLevel(lv)
	return b
}

// Message sets the record message.
func (b *recordBuilder) Message(msg string) RecordBuilder {
	b.r.WithMessage(msg)
	return b
}

// Messagef formats and sets the record message.
func (b *recordBuilder) Messagef(format string, args ...any) RecordBuilder {
	b.r.WithMessagef(format, args...)
	return b
}

// Stack adds a stack trace.
func (b *recordBuilder) Stack(stack []byte) RecordBuilder {
	b.r.WithStack(stack)
	return b
}

// Status adds a status with an optional stack trace.
func (b *recordBuilder) Status(st status.Status) RecordBuilder {
	b.r.WithStatus(st)
	return b
}

// Fields

// Field adds a field to the record.
func (b *recordBuilder) Field(key string, value any) RecordBuilder {
	b.r.WithField(key, value)
	return b
}

// Fieldf formats a field value and adds a field to the record.
func (b *recordBuilder) Fieldf(key string, format string, a ...any) RecordBuilder {
	b.r.WithFieldf(key, format, a...)
	return b
}

// Fields adds fields to the record.
func (b *recordBuilder) Fields(keyValuePairs ...any) RecordBuilder {
	b.r.WithFields(keyValuePairs...)
	return b
}

// Utility

// Trace sets the level to trace and adds the message.
func (b *recordBuilder) Trace(msg string) RecordBuilder {
	b.r.WithLevel(LevelTrace).
		WithMessage(msg)
	return b
}

// Tracef sets the level to trace and formats the message.
func (b *recordBuilder) Tracef(format string, args ...any) RecordBuilder {
	b.r.WithLevel(LevelTrace).
		WithMessagef(format, args...)
	return b
}

// Debug sets the level to debug and adds the message.
func (b *recordBuilder) Debug(msg string) RecordBuilder {
	b.r.WithLevel(LevelDebug).
		WithMessage(msg)
	return b
}

// Debugf sets the level to debug and formats the message.
func (b *recordBuilder) Debugf(format string, args ...any) RecordBuilder {
	b.r.WithLevel(LevelDebug).
		WithMessagef(format, args...)
	return b
}

// Info sets the level to info and adds the message.
func (b *recordBuilder) Info(msg string) RecordBuilder {
	b.r.WithLevel(LevelInfo).
		WithMessage(msg)
	return b
}

// Infof sets the level to info and formats the message.
func (b *recordBuilder) Infof(format string, args ...any) RecordBuilder {
	b.r.WithLevel(LevelInfo).
		WithMessagef(format, args...)
	return b
}

// Notice sets the level to notice and adds the message.
func (b *recordBuilder) Notice(msg string) RecordBuilder {
	b.r.WithLevel(LevelNotice).
		WithMessage(msg)
	return b
}

// Noticef sets the level to notice and formats the message.
func (b *recordBuilder) Noticef(format string, args ...any) RecordBuilder {
	b.r.WithLevel(LevelNotice).
		WithMessagef(format, args...)
	return b
}

// Warn sets the level to warn and adds the message.
func (b *recordBuilder) Warn(msg string) RecordBuilder {
	b.r.WithLevel(LevelWarn).
		WithMessage(msg)
	return b
}

// Warnf sets the level to warn and formats the message.
func (b *recordBuilder) Warnf(format string, args ...any) RecordBuilder {
	b.r.WithLevel(LevelWarn).
		WithMessagef(format, args...)
	return b
}

// Error sets the level to error and adds the message.
func (b *recordBuilder) Error(msg string) RecordBuilder {
	b.r.WithLevel(LevelError).
		WithMessage(msg)
	return b
}

// Errorf sets the level to error and formats the message.
func (b *recordBuilder) Errorf(format string, args ...any) RecordBuilder {
	b.r.WithLevel(LevelError).
		WithMessagef(format, args...)
	return b
}

// Fatal sets the level to fatal and adds the message.
func (b *recordBuilder) Fatal(msg string) RecordBuilder {
	b.r.WithLevel(LevelFatal).
		WithMessage(msg)
	return b
}

// Fatalf sets the level to fatal and formats the message.
func (b *recordBuilder) Fatalf(format string, args ...any) RecordBuilder {
	b.r.WithLevel(LevelFatal).
		WithMessagef(format, args...)
	return b
}
