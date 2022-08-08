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

const maxLoggerLength = 50

type textFormatter struct {
	loggerLength int32
}

func newTextFormatter() *textFormatter {
	return &textFormatter{
		loggerLength: 10,
	}
}

// Format formats the record as "INFO [main] message fields" separated by tabs.
func (f *textFormatter) Format(rec Record) string {
	b := strings.Builder{}
	b.WriteByte('\t')

	// level
	level := rec.Level.String()
	b.WriteString(level)
	f.writePadding(&b, level, 6)
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

	// message
	b.WriteString(rec.Message)

	// fields
	for i, field := range rec.Fields {
		b.WriteByte(' ')

		first := i == 0
		if first {
			b.WriteByte('\t')
		}

		key := field.Key
		value, ok := field.Value.(string)
		if !ok {
			value = fmt.Sprintf("%v", field.Value)
		}

		b.WriteString(key)
		b.WriteByte('=')
		b.WriteString(value)
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
