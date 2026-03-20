package ui

const (
	ColorReset = "\033[0m"
	ColorBold  = "\033[1m"
	ColorRed   = "\033[31m"
	ColorGreen = "\033[32m"
	ColorCyan  = "\033[36m"
)

func Reset() string {
	return ColorReset
}

func Bold(s string) string {
	return ColorBold + s + ColorReset
}

func Red(s string) string {
	return ColorRed + s + ColorReset
}

func Green(s string) string {
	return ColorGreen + s + ColorReset
}

func Cyan(s string) string {
	return ColorCyan + s + ColorReset
}

func RedBold(s string) string {
	return ColorRed + ColorBold + s + ColorReset
}

// Marked — для выделения текущего элемента в списке
func Marked(s string) string {
	return Green("→ ") + s
}
