package logging

import (
	"os"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Console *ConsoleConfig  `yaml:"console"`
	File    *FileConfig     `yaml:"file"`
	Loggers []*LoggerConfig `yaml:"loggers"`
}

type LoggerConfig struct {
	Name    string         `yaml:"name"`
	Console *ConsoleConfig `yaml:"console"`
	File    *FileConfig    `yaml:"file"`
}

type ConsoleConfig struct {
	Enabled bool   `yaml:"enabled"`
	Level   string `yaml:"level"`
}

type FileConfig struct {
	Enabled    bool   `yaml:"enabled"`
	Level      string `yaml:"level"`
	Path       string `yaml:"path"`
	MaxSize    int    `yaml:"max_size"`    // Maximum size in megabytes of a log file.
	MaxAge     int    `yaml:"max_age"`     // Maximum number of days to retain old log files.
	MaxBackups int    `yaml:"max_backups"` // Maximum number of old log files to retain.
}

func DefaultConfig() *Config {
	return &Config{
		Console: &ConsoleConfig{
			Enabled: true,
			Level:   "debug",
		},
		File: &FileConfig{
			Enabled:    false,
			Level:      "debug",
			MaxSize:    256,
			MaxAge:     7,
			MaxBackups: 10,
		},
	}
}

func ReadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := DefaultConfig()
	if err := yaml.NewDecoder(file).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
