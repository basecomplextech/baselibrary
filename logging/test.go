package logging

import "github.com/epochtimeout/library/tests"

func Test(t tests.T) Logging {
	config := TestConfig(t)
	return New(config)
}

func TestLogger(t tests.T) Logger {
	return Stderr
}

func TestConfig(t tests.T) *Config {
	return &Config{
		Console: &ConsoleConfig{
			Enabled: true,
			Level:   "debug",
		},
	}
}
