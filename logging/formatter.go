package logging

import (
	"fmt"
	"strings"
	"sync/atomic"
)

// Formatter formats records.
type Formatter interface {
	// Format formats the record.
	Format(rec Record) string
}

// text

const (
	maxLoggerLength  = 50
	maxMessageLength = 50
)

type textFormatter struct {
	color         bool
	time          bool
	loggerLength  int32
	messageLength int32
}

func newTextFormatter() *textFormatter {
	return &textFormatter{
		loggerLength: 10,
	}
}

func newTextFormatterColor(color bool, time bool) *textFormatter {
	return &textFormatter{
		color:        color,
		time:         time,
		loggerLength: 10,
	}
}

// Format formats the record as "INFO [main] message fields" separated by tabs.
func (f *textFormatter) Format(rec Record) string {
	b := strings.Builder{}

	// time
	if f.time {
		t := rec.Time
		if f.color {
			b.WriteString(colorGrey)
		}
		b.WriteString(t.Format("2006-01-02 15:04:05.000000"))
		if f.color {
			b.WriteString(colorReset)
		}
	}
	b.WriteByte('\t')

	// logger
	b.WriteString(rec.Logger)
	pad := f.loadLoggerLength()
	if len(rec.Logger) > int(pad) {
		pad = len(rec.Logger)
		f.storeLoggerLength(pad)
	}
	f.writePadding(&b, rec.Logger, pad)
	b.WriteByte(' ')
	b.WriteByte('\t')

	// level
	level := rec.Level.String()
	if f.color {
		b.WriteString(levelColor(rec.Level))
	}
	b.WriteString(level)
	if f.color {
		b.WriteString(colorReset)
	}
	f.writePadding(&b, level, 6)
	b.WriteByte('\t')

	// message
	if f.color {
		b.WriteString(levelColor(rec.Level))
		// b.WriteString(FgDefault)
	}
	b.WriteString(rec.Message)
	if f.color {
		b.WriteString(colorReset)
	}

	// message padding
	{
		pad := f.loadMessageLength()
		if len(rec.Message) > int(pad) {
			pad = len(rec.Message)
			f.storeMessageLength(pad)
		}
		f.writePadding(&b, rec.Message, pad)
	}
	b.WriteByte('\t')

	// fields
	for _, field := range rec.Fields {
		b.WriteByte(' ')

		// first := i == 0
		// if first {
		// 	// b.WriteByte('\t')
		// }

		key := field.Key
		value, ok := field.Value.(string)
		if !ok {
			value = fmt.Sprintf("%v", field.Value)
		}

		if f.color {
			b.WriteString(FgLightBlue)
			// b.WriteString(FgBlue)
		}
		b.WriteString(key)
		if f.color {
			b.WriteString(colorReset)
		}

		if f.color {
			b.WriteString(colorGrey)
		}
		b.WriteByte('=')
		if f.color {
			b.WriteString(colorReset)
		}

		if f.color {
			// FgYellow
			// b.WriteString(FgMagenta)
		}
		b.WriteString(value)
		if f.color {
			// b.WriteString(colorReset)
		}
	}

	// stack
	if len(rec.Stack) > 0 {
		stack := string(rec.Stack)
		b.WriteByte('\n')
		b.WriteString(stack)
	}

	return b.String()
}

func (f *textFormatter) writePadding(b *strings.Builder, s string, n int) {
	for i := len(s); i < n; i++ {
		b.WriteByte(' ')
	}
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
