package logging

import "github.com/complex1tech/baselibrary/terminal"

type ColorTheme struct {
	Time           string
	Logger         string
	FieldKey       string
	FieldEqualSign string
	FieldValue     string
	Levels         map[Level]string
}

func DefaultColorTheme() ColorTheme {
	return ColorTheme{
		Time:           terminal.FgGray,
		Logger:         terminal.FgDefault,
		FieldKey:       terminal.FgLightBlue,
		FieldEqualSign: terminal.FgGray,
		FieldValue:     terminal.FgDefault,
		Levels: map[Level]string{
			LevelTrace:  "",
			LevelDebug:  "",
			LevelInfo:   terminal.FgBlue,
			LevelNotice: terminal.FgGreen,
			LevelWarn:   terminal.FgYellow,
			LevelError:  terminal.FgRed,
			LevelFatal:  terminal.FgRed,
		},
	}
}

func (th ColorTheme) Level(level Level) string {
	if th.Levels == nil {
		return ""
	}
	return th.Levels[level]
}
