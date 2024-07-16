package profiling

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	DefaultBlockProfileRate = 10_000_000 // ns
	DefalutMutexProfileRate = 100        // 1/100 events
)

// Config specifies a profiling configuration.
type Config struct {
	Enabled bool         `json:"enabled"`
	CPU     CPUConfig    `json:"cpu"`
	Memory  MemoryConfig `json:"memory"`
	Block   BlockConfig  `json:"block"`
	Mutex   MutexConfig  `json:"mutex"`
}

type CPUConfig struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`

	// Rate sets the CPU profiling rate to hz samples per second.
	//
	// 0 means the default rate (100 hz), see runtime.SetCPUProfileRate.
	Rate int `json:"rate"`
}

type MemoryConfig struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`

	// Rate controls the fraction of memory allocations that are recorded
	// And reported in the memory profile. The profiler aims to sample
	// An average of one allocation per MemProfileRate bytes allocated.
	//
	// To include every allocated block in the profile, set MemProfileRate to 1.
	//
	// 0 means the default rate (512 kb), see runtime.SetMemProfileRate.
	Rate int `json:"rate"`
}

type BlockConfig struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`

	// SetBlockProfileRate controls the fraction of goroutine blocking events
	// That are reported in the blocking profile. The profiler aims to sample
	// An average of one blocking event per rate nanoseconds spent blocked.
	//
	// To include every blocking event in the profile, pass rate = 1.
	//
	// 0 means the default rate (10_000_000 ns), see runtime.SetBlockProfileRate.
	Rate int `json:"rate"`
}

type MutexConfig struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`

	// Rate controls the fraction of mutex contention events that are reported
	// In the mutex profile. On average 1/rate events are reported.
	//
	// 0 means the default rate (1/100), see runtime.SetMutexProfileFraction.
	Rate int `json:"rate"`
}

// DefaultConfig returns the default profiling config.
func DefaultConfig() *Config {
	return &Config{
		Enabled: false,

		CPU: CPUConfig{
			Enabled: true,
			Path:    "cpu.pprof",
		},
		Memory: MemoryConfig{
			Enabled: true,
			Path:    "memory.pprof",
		},
		Block: BlockConfig{
			Path: "block.pprof",
		},
		Mutex: MutexConfig{
			Path: "mutex.pprof",
		},
	}
}

// ReadConfig reads a profiling config from a JSON or YAML file.
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
	// Init default
	config := DefaultConfig()

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse json
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}

func readConfigYaml(path string) (*Config, error) {
	// Init default
	config := DefaultConfig()

	// Open file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Parse yaml
	if err := yaml.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}
