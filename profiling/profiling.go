package profiling

import (
	"errors"
	"os"
	"sync"
)

var (
	mu   sync.Mutex
	prof *profiler
)

// Init initializes the profiling service, uses the default config when config is nil.
func Init(config *Config) error {
	mu.Lock()
	defer mu.Unlock()

	if prof != nil {
		return errors.New("profiling already initialized")
	}

	if config == nil {
		config = DefaultConfig()
	}

	prof = newProfiler(config)
	return nil
}

// Load loads a profiling config from a json file if it exists and initializes the profiling service.
func Load(path string) error {
	config, err := LoadConfig(path)
	switch {
	case os.IsNotExist(err):
		config = DefaultConfig()
	case err != nil:
		return err
	}

	return Init(config)
}

// Start starts enabled profilers.
func Start() error {
	mu.Lock()
	defer mu.Unlock()

	if prof == nil {
		return errors.New("profiling not initialized")
	}

	return prof.Start()
}

// Stop stops running profilers, writes all profiles to disk.
func Stop() error {
	mu.Lock()
	defer mu.Unlock()

	if prof == nil {
		return nil
	}

	defer func() {
		prof = nil
	}()
	return prof.Stop()
}
