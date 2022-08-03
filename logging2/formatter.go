package logging2

import (
	"fmt"
	"strings"
)

// Formatter formats records.
type Formatter interface {
	// Format formats the record.
	Format(rec Record) string
}

// text

type textFormatter struct{}

func newTextFormatter() *textFormatter {
	return &textFormatter{}
}

// Format formats the record as "INFO [main] message fields" separated by tabs.
func (f *textFormatter) Format(rec Record) string {
	b := strings.Builder{}
	b.WriteByte('\t')

	// level
	level := rec.Level.String()
	b.WriteString(level)
	f.writePadding(&b, level, 5)
	b.WriteByte('\t')

	// logger
	b.WriteByte('[')
	b.WriteString(rec.Logger)
	b.WriteByte(']')
	f.writePadding(&b, rec.Logger, 10)
	b.WriteByte('\t')

	// message
	b.WriteString(rec.Message)

	// fields
	for i, field := range rec.Fields {
		first := i == 0
		last := i == len(rec.Fields)-1

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

		if !last {
			b.WriteByte(' ')
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
