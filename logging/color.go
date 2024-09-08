// Copyright 2022 Ivan Korobkov. All rights reserved.

package logging

import "github.com/basecomplextech/baselibrary/terminal"

// ColorTheme specifies the terminal colors.
type ColorTheme struct {
	Time terminal.Color
	// Logger         terminal.Color
	FieldKey       terminal.Color
	FieldEqualSign terminal.Color
	FieldValue     terminal.Color
	Levels         map[Level]terminal.Color
}

// DefaultColorTheme returns the default terminal colors.
func DefaultColorTheme() ColorTheme {
	return ColorTheme{
		Time: terminal.Gray,
		// Logger:         terminal.Default,
		FieldKey:       terminal.LightBlue,
		FieldEqualSign: terminal.Gray,
		FieldValue:     terminal.Default,
		Levels: map[Level]terminal.Color{
			LevelTrace:  "",
			LevelDebug:  "",
			LevelInfo:   terminal.Blue,
			LevelNotice: terminal.Green,
			LevelWarn:   terminal.Yellow,
			LevelError:  terminal.Red,
			LevelFatal:  terminal.Red,
		},
	}
}

// Level returns a color for the given level.
func (th ColorTheme) Level(level Level) terminal.Color {
	if th.Levels == nil {
		return ""
	}
	return th.Levels[level]
}
