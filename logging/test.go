package logging

import "github.com/epochtimeout/baselibrary/tests"

func Test(t tests.T) Logging {
	config := TestConfig()
	l, err := New(config)
	if err != nil {
		t.Fatal(err)
	}
	return l
}

func TestLogger(t tests.T) Logger {
	return Stdout
}

func TestConfig() *Config {
	return &Config{
		Console: &ConsoleConfig{
			Enabled: true,
			Level:   LevelDebug,
		},
	}
}
