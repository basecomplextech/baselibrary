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

	// level
	b.WriteString(rec.Level.String())
	b.WriteByte('\t')

	// logger
	b.WriteByte('[')
	b.WriteString(rec.Logger)
	b.WriteByte(']')
	b.WriteByte('\t')

	// message
	b.WriteString(rec.Message)

	// fields
	for i, field := range rec.Fields {
		first := i == 0
		last := i == len(rec.Fields)-1

		switch {
		case first:
			b.WriteByte('\t')
		case !last:
			b.WriteByte(' ')
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

	return b.String()
}
