package logging

import (
	"bytes"
	"fmt"
	"io"
	"strings"
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

	timeLen    int32
	loggerLen  int32
	messageLen int32
}

func newTextFormatter(color bool) *textFormatter {
	return &textFormatter{
		color: color,
		theme: DefaultColorTheme(),

		loggerLen:  10,
		messageLen: 20,
	}
}

// Format formats the record as "time logger level message fields" separated by tabs.
func (f *textFormatter) Format(w io.Writer, rec *Record) error {
	tw := terminal.NewWriterColor(w, f.color)
	f.writeTime(tw, rec.Time)
	f.writeLogger(tw, rec.Logger)
	f.writeLevel(tw, rec.Level)
	f.writeMessage(tw, rec)
	f.writeFields(tw, rec.Fields)
	f.writeStack(tw, rec.Stack)
	tw.WriteString("\n")
	return nil
}

// time

func (f *textFormatter) writeTime(w *terminal.Writer, t time.Time) {
	s := t.Format("2006-01-02T15:04:05.999Z07:00")
	w.Color(f.theme.Time)
	w.WriteString(s)
	w.ResetColor()

	f.maybePadding(w, len(s), 0, &f.timeLen)
	w.WriteString("\t")
}

// logger

func (f *textFormatter) writeLogger(w *terminal.Writer, logger string) {
	w.Color(f.theme.Logger)
	w.WriteString(logger)
	w.ResetColor()

	f.maybePadding(w, len(logger), maxLoggerLength, &f.loggerLen)
	w.WriteString(" ")
	w.WriteString("\t")
}

// level

func (f *textFormatter) writeLevel(w *terminal.Writer, level Level) {
	s := level.String()
	color := f.theme.Level(level)
	w.Color(color)
	w.WriteString(s)
	w.ResetColor()

	f.writePadding(w, len(s), 6)
	w.WriteString("\t")
}

// message

func (f *textFormatter) writeMessage(w *terminal.Writer, rec *Record) {
	color := f.theme.Level(rec.Level)
	w.Color(color)
	w.WriteString(rec.Message)
	w.ResetColor()

	f.maybePadding(w, len(rec.Message), maxMessageLength, &f.messageLen)
	w.WriteString("\t")
}

// fields

func (f *textFormatter) writeFields(w *terminal.Writer, fields []Field) {
	var valueBytes []byte // scratch buffer

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
			valueBytes = fmt.Appendf(valueBytes, "%v", v)
		}

		w.Color(f.theme.FieldKey)
		w.WriteString(key)
		w.ResetColor()

		w.Color(f.theme.FieldEqualSign)
		w.WriteString("=")
		w.ResetColor()

		w.Color(f.theme.FieldValue)
		if len(valueBytes) > 0 {
			f.writeFieldValueBytes(w, valueBytes)
			valueBytes = valueBytes[:0]
		} else {
			f.writeFieldValue(w, value)
		}
		w.ResetColor()
	}
}

func (f *textFormatter) writeFieldValue(w *terminal.Writer, s string) {
	i := strings.IndexAny(s, " \t")
	if i > 0 {
		w.WriteString("'")
	}

	w.WriteString(s)

	if i > 0 {
		w.WriteString("'")
	}
}

func (f *textFormatter) writeFieldValueBytes(w *terminal.Writer, b []byte) {
	i := bytes.IndexAny(b, " \t")
	if i > 0 {
		w.WriteString("'")
	}

	w.Write(b)

	if i > 0 {
		w.WriteString("'")
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

// private

func (f *textFormatter) maybePadding(w *terminal.Writer, n int, max int32, ptr *int32) error {
	prev := atomic.LoadInt32(ptr)

	// maybe write padding
	if n <= int(prev) {
		return f.writePadding(w, n, int(prev))
	}

	// check already max
	if max > 0 && prev == max {
		return nil
	}

	// update previous
	next := int32(n)
	if max > 0 && next > max {
		next = max
	}
	atomic.CompareAndSwapInt32(ptr, prev, next)
	return nil
}

func (f *textFormatter) writePadding(w *terminal.Writer, n int, total int) error {
	for i := n; i < total; i++ {
		if _, err := w.WriteString(" "); err != nil {
			return err
		}
	}
	return nil
}
