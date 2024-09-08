// Copyright 2022 Ivan Korobkov. All rights reserved.

package terminal

// Color is a human-readable color.
type Color string

const (
	Default Color = "default"

	Black   Color = "black"
	Red     Color = "red"
	Green   Color = "green"
	Yellow  Color = "yellow"
	Blue    Color = "blue"
	Magenta Color = "magenta"
	Cyan    Color = "cyan"
	White   Color = "white"
	Gray    Color = "gray"

	LightRed     Color = "light_red"
	LightGreen   Color = "light_green"
	LightYellow  Color = "light_yellow"
	LightBlue    Color = "light_blue"
	LightMagenta Color = "light_magenta"
	LightCyan    Color = "light_cyan"
	LightWhite   Color = "light_white"
)

// Code returns the terminal color code.
func (c Color) Code() ColorCode {
	return colorCodeMap[c]
}

// ColorCode is a terminal color code.
type ColorCode string

const (
	FgReset ColorCode = "\033[0m"

	FgBlack   ColorCode = "\033[0;30m"
	FgRed     ColorCode = "\033[0;31m"
	FgGreen   ColorCode = "\033[0;32m"
	FgYellow  ColorCode = "\033[0;33m"
	FgBlue    ColorCode = "\033[0;34m"
	FgMagenta ColorCode = "\033[0;35m"
	FgCyan    ColorCode = "\033[0;36m"
	FgWhite   ColorCode = "\033[0;37m"
	FgDefault ColorCode = "\033[0;39m"

	FgGray         ColorCode = "\033[0;90m"
	FgLightRed     ColorCode = "\033[0;91m"
	FgLightGreen   ColorCode = "\033[0;92m"
	FgLightYellow  ColorCode = "\033[0;93m"
	FgLightBlue    ColorCode = "\033[0;94m"
	FgLightMagenta ColorCode = "\033[0;95m"
	FgLightCyan    ColorCode = "\033[0;96m"
	FgLightWhite   ColorCode = "\033[0;97m"
)

// private

var colorCodeMap = map[Color]ColorCode{
	Default: FgDefault,

	Black:   FgBlack,
	Red:     FgRed,
	Green:   FgGreen,
	Yellow:  FgYellow,
	Blue:    FgBlue,
	Magenta: FgMagenta,
	Cyan:    FgCyan,
	White:   FgWhite,
	Gray:    FgGray,

	LightRed:     FgLightRed,
	LightGreen:   FgLightGreen,
	LightYellow:  FgLightYellow,
	LightBlue:    FgLightBlue,
	LightMagenta: FgLightMagenta,
	LightCyan:    FgLightCyan,
	LightWhite:   FgLightWhite,
}
