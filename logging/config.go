// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config specifies the logging configuration.
type Config struct {
	Console *ConsoleConfig `json:"console"`
	File    *FileConfig    `json:"file"`
}

type ConsoleConfig struct {
	Enabled bool       `json:"enabled"`
	Level   Level      `json:"level"`
	Color   bool       `json:"color"`
	Theme   ColorTheme `json:"theme"`
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
			Color:   true,
			Theme:   DefaultColorTheme(),
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

// ReadConfig reads a configuration from a JSON or YAML file.
func ReadConfig(path string) (*Config, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		return readConfigYaml(path)
	}
	return readConfigJson(path)
}

// private

func readConfigJson(path string) (*Config, error) {
	config := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}

func readConfigYaml(path string) (*Config, error) {
	config := DefaultConfig()

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}
