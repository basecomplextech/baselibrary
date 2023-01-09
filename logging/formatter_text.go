package logging

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"github.com/complex1tech/baselibrary/terminal"
)

const (
	maxLoggerLength  = 50
	maxMessageLength = 50
)

type textFormatter struct {
	color bool
	theme ColorTheme

	loggerLength  int32
	messageLength int32
}

func newTextFormatter(color bool) *textFormatter {
	return &textFormatter{
		color: color,
		theme: DefaultColorTheme(),

		loggerLength:  10,
		messageLength: 20,
	}
}

// Format formats the record as "time logger level message fields" separated by tabs.
func (f *textFormatter) Format(w io.Writer, rec Record) error {
	tw := terminal.NewWriterColor(w, f.color)
	f.writeTime(tw, rec.Time)
	f.writeLogger(tw, rec.Logger)
	f.writeLevel(tw, rec.Level)
	f.writeMessage(tw, rec)
	f.writeFields(tw, rec.Fields)
	f.writeStack(tw, rec.Stack)
	return nil
}

// time

func (f *textFormatter) writeTime(w *terminal.Writer, time time.Time) {
	w.Color(f.theme.Time)
	w.WriteString(time.Format("2006-01-02 15:04:05.000000"))
	w.ResetColor()
	w.WriteString("\t")
}

// logger

func (f *textFormatter) writeLogger(w *terminal.Writer, logger string) {
	w.Color(f.theme.Logger)
	w.WriteString(logger)
	w.ResetColor()

	f.writeLoggerPadding(w, logger)
	w.WriteString(" ")
	w.WriteString("\t")
}

func (f *textFormatter) writeLoggerPadding(w *terminal.Writer, logger string) {
	n := f.loadLoggerLength()
	if len(logger) <= n {
		f.writePadding(w, logger, n)
		return
	}

	n = len(logger)
	f.storeLoggerLength(n)
}

func (f *textFormatter) loadLoggerLength() int {
	length := atomic.LoadInt32(&f.loggerLength)
	return int(length)
}

func (f *textFormatter) storeLoggerLength(length int) {
	prev := atomic.LoadInt32(&f.loggerLength)
	if prev == maxLoggerLength {
		return
	}
	if length > maxLoggerLength {
		length = maxLoggerLength
	}
	atomic.CompareAndSwapInt32(&f.loggerLength, prev, int32(length))
}

// level

func (f *textFormatter) writeLevel(w *terminal.Writer, level Level) {
	s := level.String()
	color := f.theme.Level(level)
	w.WriteString(color)
	w.WriteString(s)
	w.ResetColor()

	f.writePadding(w, s, 6)
	w.WriteString("\t")
}

// message

func (f *textFormatter) writeMessage(w *terminal.Writer, rec Record) {
	color := f.theme.Level(rec.Level)
	w.Color(color)
	w.WriteString(rec.Message)
	w.ResetColor()

	f.writeMessagePadding(w, rec.Message)
	w.WriteString("\t")
}

func (f *textFormatter) writeMessagePadding(w *terminal.Writer, message string) {
	pad := f.loadMessageLength()
	if len(message) > int(pad) {
		pad = len(message)
		f.storeMessageLength(pad)
	}
	f.writePadding(w, message, pad)
}

func (f *textFormatter) loadMessageLength() int {
	length := atomic.LoadInt32(&f.messageLength)
	return int(length)
}

func (f *textFormatter) storeMessageLength(length int) {
	prev := atomic.LoadInt32(&f.messageLength)
	if prev == maxMessageLength {
		return
	}
	if length > maxMessageLength {
		length = maxMessageLength
	}
	atomic.CompareAndSwapInt32(&f.messageLength, prev, int32(length))
}

// fields

func (f *textFormatter) writeFields(w *terminal.Writer, fields []Field) {
	for _, field := range fields {
		w.WriteString(" ")

		key := field.Key
		var value string

		switch v := field.Value.(type) {
		case string:
			value = v
		case fmt.Stringer:
			value = v.String()
		default:
			value = fmt.Sprintf("%v", v)
		}

		w.Color(f.theme.FieldKey)
		w.WriteString(key)
		w.ResetColor()

		w.Color(f.theme.FieldEqualSign)
		w.WriteString("=")
		w.ResetColor()

		w.Color(f.theme.FieldValue)
		w.WriteString(value)
		w.ResetColor()
	}
}

// stack

func (f *textFormatter) writeStack(w *terminal.Writer, stack []byte) {
	if len(stack) == 0 {
		return
	}

	w.WriteString("\n")
	w.Write(stack)
}

// padding

func (f *textFormatter) writePadding(w *terminal.Writer, s string, n int) error {
	for i := len(s); i < n; i++ {
		if _, err := w.WriteString(" "); err != nil {
			return err
		}
	}
	return nil
}
