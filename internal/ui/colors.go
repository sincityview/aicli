package ui

const (
	ColorGreen = "\033[32m"
	ColorRed   = "\033[31m"
	ColorCyan  = "\033[36m"
	ColorBold  = "\033[1m"
	ColorReset = "\033[0m"
)

func Success(s string) string { return ColorGreen + s + ColorReset }
func Error(s string) string   { return ColorRed + s + ColorReset }
func Info(s string) string    { return ColorCyan + s + ColorReset }
func Bold(s string) string    { return ColorBold + s + ColorReset }
