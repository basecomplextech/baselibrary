package logging2

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config specifies the logging configuration.
type Config struct {
	Level   Level          `json:"level"`
	Console *ConsoleConfig `json:"console"`
	File    *FileConfig    `json:"file"`
}

type ConsoleConfig struct {
	Enabled bool  `json:"enabled"`
	Level   Level `json:"level"`
}

type FileConfig struct {
	Enabled    bool   `yaml:"enabled"`
	Path       string `yaml:"path"`
	Level      Level  `yaml:"level"`
	MaxSize    int    `yaml:"max_size"`    // Maximum log file size in megabytes
	MaxAge     int    `yaml:"max_age"`     // Maximum days to retain old log files
	MaxBackups int    `yaml:"max_backups"` // Maximum number of old log files to retain
}

// DefaultConfig returns a default configuration.
func DefaultConfig() *Config {
	name := os.Args[0]
	name, _, _ = strings.Cut(name, ".")
	path := fmt.Sprintf("%v.log", name)

	return &Config{
		Console: &ConsoleConfig{
			Enabled: true,
			Level:   LevelInfo,
		},
		File: &FileConfig{
			Enabled:    false,
			Path:       path,
			Level:      LevelDebug,
			MaxSize:    256,
			MaxAge:     7,
			MaxBackups: 10,
		},
	}
}

// ReadConfig reads a configuration from a YAML file.
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
