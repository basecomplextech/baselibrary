package logging

import (
	"log"
)

// Writer writes log records.
type Writer interface {
	// Write writes a record.
	Write(rec Record) error
}

// internal

var _ Writer = (*writer)(nil)

const lflags = log.Ldate | log.Ltime | log.Lmicroseconds

type writer struct {
	level     Level
	logger    *log.Logger
	formatter Formatter
}

func newWriter(level Level, logger *log.Logger, formatter Formatter) *writer {
	return &writer{
		level:     level,
		logger:    logger,
		formatter: formatter,
	}
}

// Write writes a record.
func (w *writer) Write(rec Record) error {
	if rec.Level < w.level {
		return nil
	}

	msg := w.formatter.Format(rec)
	return w.logger.Output(2, msg)
}
