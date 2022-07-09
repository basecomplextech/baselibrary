package logging

import (
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

var _ Writer = &writer{}

type writer struct {
	level  Level
	logger *log.Logger
}

func openFileWriter(config *FileConfig) *writer {
	if !config.Enabled {
		return nil
	}

	level, err := ParseLevel(config.Level)
	if err != nil {
		log.Fatal(err)
	}

	if err := checkLogFile(config.Path); err != nil {
		log.Fatal(err)
	}

	rotated := &lumberjack.Logger{
		Filename:   config.Path,
		MaxSize:    config.MaxSize,
		MaxAge:     config.MaxAge,
		MaxBackups: config.MaxBackups,
		LocalTime:  true,
	}

	return newWriter(level, rotated)
}

func openConsoleWriter(config *ConsoleConfig) *writer {
	if !config.Enabled {
		return nil
	}

	level, err := ParseLevel(config.Level)
	if err != nil {
		log.Fatal(err)
	}

	return newWriter(level, os.Stderr)
}

func newWriter(level Level, w io.Writer) *writer {
	logger := log.New(w, "", log.LstdFlags)

	return &writer{
		level:  level,
		logger: logger,
	}
}

func (w *writer) Level() Level {
	return w.level
}

func (w *writer) Write(record Record) {
	if w.level > record.Level {
		return
	}

	msg := formatRecord(record)
	w.logger.Output(4, msg)
}

func formatRecord(r Record) string {
	var msg string

	if r.Format == "" {
		msg = fmt.Sprint(r.Args...)
	} else {
		msg = fmt.Sprintf(r.Format, r.Args...)
	}

	msg = fmt.Sprintf("\t%-5v\t%-10v\t%v", r.Level.Name(), r.Logger, msg)
	if r.Stack != nil {
		msg += "\n" + string(r.Stack)
	}
	return msg
}

func checkLogFile(name string) error {
	// Open or create a file.
	file, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0600)
	if os.IsNotExist(err) {
		file, err = os.Create(name)
	}
	if err != nil {
		return err
	}

	file.Close()
	return nil
}
