package logging

import (
	"fmt"
)

type RecordBuilder interface {
	// Send sends the record.
	Send()

	// Level

	// Level sets the record level.
	Level(level Level) RecordBuilder

	// Message

	// Message sets the record message.
	Message(msg string) RecordBuilder

	// Messagef formats and sets the record message.
	Messagef(msg string, args ...any) RecordBuilder

	// Fields

	// Field adds a field to the record.
	Field(key string, value any) RecordBuilder

	// Fieldf formats a field value and adds a field to the record.
	Fieldf(key string, format string, a ...any) RecordBuilder

	// Fields adds fields to the record.
	Fields(keyValuePairs ...any) RecordBuilder

	// Stack adds a stack trace.
	Stack(stack []byte) RecordBuilder

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
	l   *logging
	rec Record
}

func newRecordBuilder(l *logging, logger string) *recordBuilder {
	return &recordBuilder{
		l:   l,
		rec: NewRecord(logger),
	}
}

// Send sends the record.
func (b *recordBuilder) Send() {
	b.l.send(b.rec)
}

// Level

// Level sets the record level.
func (b *recordBuilder) Level(level Level) RecordBuilder {
	b.rec.Level = level
	return b
}

// Message

// Message sets the record message.
func (b *recordBuilder) Message(msg string) RecordBuilder {
	b.rec.Message = msg
	return b
}

// Messagef formats and sets the record message.
func (b *recordBuilder) Messagef(format string, args ...any) RecordBuilder {
	b.rec.Message = fmt.Sprintf(format, args...)
	return b
}

// Fields

// Field adds a field to the record.
func (b *recordBuilder) Field(key string, value any) RecordBuilder {
	field := NewField(key, value)
	b.rec.Fields = append(b.rec.Fields, field)
	return b
}

// Fieldf formats a field value and adds a field to the record.
func (b *recordBuilder) Fieldf(key string, format string, a ...any) RecordBuilder {
	value := fmt.Sprintf(format, a...)
	return b.Field(key, value)
}

// Fields adds fields to the record.
func (b *recordBuilder) Fields(keyValuePairs ...any) RecordBuilder {
	fields := NewFields(keyValuePairs...)
	b.rec.Fields = append(b.rec.Fields, fields...)
	return b
}

// Stack adds a stack trace.
func (b *recordBuilder) Stack(stack []byte) RecordBuilder {
	b.rec.Stack = stack
	return b
}

// Utility

// Trace sets the level to trace and adds the message.
func (b *recordBuilder) Trace(msg string) RecordBuilder {
	return b.Level(LevelTrace).Message(msg)
}

// Tracef sets the level to trace and formats the message.
func (b *recordBuilder) Tracef(format string, args ...any) RecordBuilder {
	return b.Level(LevelTrace).Messagef(format, args...)
}

// Debug sets the level to debug and adds the message.
func (b *recordBuilder) Debug(msg string) RecordBuilder {
	return b.Level(LevelDebug).Message(msg)
}

// Debugf sets the level to debug and formats the message.
func (b *recordBuilder) Debugf(format string, args ...any) RecordBuilder {
	return b.Level(LevelDebug).Messagef(format, args...)
}

// Info sets the level to info and adds the message.
func (b *recordBuilder) Info(msg string) RecordBuilder {
	return b.Level(LevelInfo).Message(msg)
}

// Infof sets the level to info and formats the message.
func (b *recordBuilder) Infof(format string, args ...any) RecordBuilder {
	return b.Level(LevelInfo).Messagef(format, args...)
}

// Notice sets the level to notice and adds the message.
func (b *recordBuilder) Notice(msg string) RecordBuilder {
	return b.Level(LevelNotice).Message(msg)
}

// Noticef sets the level to notice and formats the message.
func (b *recordBuilder) Noticef(format string, args ...any) RecordBuilder {
	return b.Level(LevelNotice).Messagef(format, args...)
}

// Warn sets the level to warn and adds the message.
func (b *recordBuilder) Warn(msg string) RecordBuilder {
	return b.Level(LevelWarn).Message(msg)
}

// Warnf sets the level to warn and formats the message.
func (b *recordBuilder) Warnf(format string, args ...any) RecordBuilder {
	return b.Level(LevelWarn).Messagef(format, args...)
}

// Error sets the level to error and adds the message.
func (b *recordBuilder) Error(msg string) RecordBuilder {
	return b.Level(LevelError).Message(msg)
}

// Errorf sets the level to error and formats the message.
func (b *recordBuilder) Errorf(format string, args ...any) RecordBuilder {
	return b.Level(LevelError).Messagef(format, args...)
}

// Fatal sets the level to fatal and adds the message.
func (b *recordBuilder) Fatal(msg string) RecordBuilder {
	return b.Level(LevelFatal).Message(msg)
}

// Fatalf sets the level to fatal and formats the message.
func (b *recordBuilder) Fatalf(format string, args ...any) RecordBuilder {
	return b.Level(LevelFatal).Messagef(format, args...)
}
