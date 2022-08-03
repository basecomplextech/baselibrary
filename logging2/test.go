package logging2

import "github.com/epochtimeout/baselibrary/tests"

func Test(t tests.T) Logging {
	config := TestConfig()
	l, err := New(config)
	if err != nil {
		t.Fatal(err)
	}
	return l
}

func TestLogger() Logger {
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
