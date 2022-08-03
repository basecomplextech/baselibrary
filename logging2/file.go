package logging2

import (
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

func newFileWriter(config *FileConfig) (*writer, error) {
	if !config.Enabled {
		return nil, nil
	}

	path := config.Path
	if err := createFile(path); err != nil {
		return nil, err
	}

	out := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    config.MaxSize,
		MaxAge:     config.MaxAge,
		MaxBackups: config.MaxBackups,
		LocalTime:  true,
	}

	level := config.Level
	logger := log.New(out, "", lflags)
	formatter := newTextFormatter()
	return newWriter(level, logger, formatter), nil
}

func createFile(name string) error {
	// open or create file
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
