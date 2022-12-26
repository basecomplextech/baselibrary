package logging

// Foreground colors. basic foreground colors 30 - 37
const (
	FgBlack   = "\033[0;30m" //  int = iota + 30
	FgRed     = "\033[0;31m"
	FgGreen   = "\033[0;32m"
	FgYellow  = "\033[0;33m"
	FgBlue    = "\033[0;34m"
	FgMagenta = "\033[0;35m"
	FgCyan    = "\033[0;36m"
	FgWhite   = "\033[0;37m"
	FgDefault = "\033[0;39m"
)

// Extra foreground color 90 - 97(非标准)
const (
	FgDarkGray     = "\033[0;90m" //  int = iota + 90 // 亮黑（灰）
	FgLightRed     = "\033[0;91m"
	FgLightGreen   = "\033[0;92m"
	FgLightYellow  = "\033[0;93m"
	FgLightBlue    = "\033[0;94m"
	FgLightMagenta = "\033[0;95m"
	FgLightCyan    = "\033[0;96m"
	FgLightWhite   = "\033[0;97m"
	// FgGray is alias of FgDarkGray
	FgGray = "\033[0;90m" // int = 90 // 亮黑（灰）
)

var (
	colorPureRed      = "\033[0;31m"
	colorDarkGreen    = "\033[0;32m"
	colorOrange       = "\033[0;33m"
	colorDarkBlue     = "\033[0;34m"
	colorBrightPurple = "\033[0;35m"
	colorDarkCyan     = "\033[0;36m"
	colorDullWhite    = "\033[0;37m"
	colorPureBlack    = "\033[0;30m"
	colorBrightRed    = "\033[0;91m"
	colorLightGreen   = "\033[0;92m"
	colorYellow       = "\033[0;93m"
	colorBrightBlue   = "\033[0;94m"
	colorMagenta      = "\033[0;95m"
	colorLightCyan    = "\033[0;96m"
	colorBrightBlack  = "\033[0;90m"
	colorBrightWhite  = "\033[0;97m"
	colorGrey         = "\x1b[90m"

	colorCyanBack   = "\033[0;46m"
	colorPurpleBack = "\033[0;45m"
	colorWhiteBack  = "\033[0;47m"
	colorBlueBack   = "\033[0;44m"
	colorOrangeBack = "\033[0;43m"
	colorGreenBack  = "\033[0;42m"
	colorPinkBack   = "\033[0;41m"
	colorGreyBack   = "\033[0;40m"

	colorBold      = "\033[1m"
	colorUnderline = "\033[4m"
	colorItalic    = "\033[3m"
	colorDarken    = "\033[2m"
	colorInvisible = "\033[08m"
	colorReverse   = "\033[07m"
	colorReset     = "\033[0m"
)

func levelColor(lv Level) string {
	switch lv {
	case LevelTrace:
	case LevelDebug:
	case LevelInfo:
		return FgBlue
		return FgCyan
	case LevelNotice:
		return FgGreen
	case LevelWarn:
		return FgYellow
	case LevelError:
		return FgRed
	case LevelFatal:
		return FgRed
	}
	return ""
}
