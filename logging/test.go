package logging

func Test() Logging {
	return New(TestConfig())
}

func TestLogger() Logger {
	return Stderr
}

func TestConfig() *Config {
	return &Config{
		Console: &ConsoleConfig{
			Enabled: true,
			Level:   "debug",
		},
	}
}
