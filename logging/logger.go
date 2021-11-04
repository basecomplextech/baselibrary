package logging

import (
	"fmt"
	"io"
	"runtime/debug"
	"sync"
)

var _ Logger = &logger{}

type logger struct {
	name    string
	parent  *logger
	writers []Writer

	childMu  *sync.Mutex
	children map[string]*logger
}

func openLogger(parent *logger, config *LoggerConfig) *logger {
	var ww []Writer

	if config.Console != nil && config.Console.Enabled {
		w := openConsoleWriter(config.Console)
		ww = append(ww, w)
	}

	if config.File != nil && config.File.Enabled {
		w := openFileWriter(config.File)
		ww = append(ww, w)
	}

	return newLogger(parent, config.Name, ww...)
}

func stdLogger(name string, out io.Writer) *logger {
	w := newWriter(LevelDebug, out)
	return newLogger(nil, name, w)
}

func newLogger(parent *logger, name string, writers ...Writer) *logger {
	return &logger{
		name:    name,
		parent:  parent,
		writers: writers,

		childMu:  &sync.Mutex{},
		children: make(map[string]*logger),
	}
}

// Logger returns a child logger.
func (l *logger) Logger(name string) Logger {
	l.childMu.Lock()
	defer l.childMu.Unlock()

	child, ok := l.children[name]
	if ok {
		return child
	}

	fullname := fmt.Sprintf("%s/%s", l.name, name)
	child = newLogger(l, fullname, l.writers...)
	l.children[name] = child
	return child
}

func (l *logger) Log(record Record) {
	for _, w := range l.writers {
		w.Write(record)
	}

	if l.parent != nil {
		l.parent.Log(record)
	}
}

func (l *logger) Enabled(level Level) bool {
	for _, w := range l.writers {
		if w.Level() <= level {
			return true
		}
	}

	if l.parent != nil {
		return l.parent.Enabled(level)
	}
	return false
}

func (l *logger) Trace(args ...interface{}) {
	r := Record{
		Logger: l.name,
		Level:  LevelTrace,
		Args:   args,
	}
	l.Log(r)
}

func (l *logger) Tracef(format string, args ...interface{}) {
	r := Record{
		Logger: l.name,
		Level:  LevelTrace,
		Format: format,
		Args:   args,
	}
	l.Log(r)
}

func (l *logger) TraceEnabled() bool {
	return l.Enabled(LevelTrace)
}

func (l *logger) Debug(args ...interface{}) {
	r := Record{
		Logger: l.name,
		Level:  LevelDebug,
		Args:   args,
	}
	l.Log(r)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	r := Record{
		Logger: l.name,
		Level:  LevelDebug,
		Format: format,
		Args:   args,
	}
	l.Log(r)
}

func (l *logger) DebugEnabled() bool {
	return l.Enabled(LevelDebug)
}

func (l *logger) Info(args ...interface{}) {
	r := Record{
		Logger: l.name,
		Level:  LevelInfo,
		Args:   args,
	}
	l.Log(r)
}

func (l *logger) Infof(format string, args ...interface{}) {
	r := Record{
		Logger: l.name,
		Level:  LevelInfo,
		Format: format,
		Args:   args,
	}
	l.Log(r)
}

func (l *logger) InfoEnabled() bool {
	return l.Enabled(LevelInfo)
}

func (l *logger) Warn(args ...interface{}) {
	r := Record{
		Logger: l.name,
		Level:  LevelWarn,
		Args:   args,
	}
	l.Log(r)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	r := Record{
		Logger: l.name,
		Level:  LevelWarn,
		Format: format,
		Args:   args,
	}
	l.Log(r)
}

func (l *logger) WarnEnabled() bool {
	return l.Enabled(LevelWarn)
}

func (l *logger) Error(args ...interface{}) {
	r := Record{
		Logger: l.name,
		Level:  LevelError,
		Args:   args,
	}
	l.Log(r)

}
func (l *logger) Errorf(format string, args ...interface{}) {
	r := Record{
		Logger: l.name,
		Level:  LevelError,
		Format: format,
		Args:   args,
	}
	l.Log(r)
}

func (l *logger) ErrorEnabled() bool {
	return l.Enabled(LevelError)
}

func (l *logger) Panic(args ...interface{}) {
	r := Record{
		Logger: l.name,
		Level:  LevelError,
		Args:   args,
		Stack:  debug.Stack(),
	}
	l.Log(r)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	r := Record{
		Logger: l.name,
		Level:  LevelError,
		Format: format,
		Args:   args,
		Stack:  debug.Stack(),
	}
	l.Log(r)
}
