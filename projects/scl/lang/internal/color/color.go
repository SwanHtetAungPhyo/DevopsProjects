package color

var (
	// ANSI Color codes
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
	ColorBold    = "\033[1m"
	ColorDim     = "\033[2m"

	// Emoji constants
	EmojiSuccess = "✓"
	EmojiError   = "✗"
	EmojiWarning = "⚠"
	EmojiInfo    = "ℹ"
	EmojiRocket  = "🚀"
	EmojiGear    = "⚙"
	EmojiFile    = "📄"
	EmojiCheck   = "✔"
)

// DisableColors disables all color output
func DisableColors() {
	ColorReset = ""
	ColorRed = ""
	ColorGreen = ""
	ColorYellow = ""
	ColorBlue = ""
	ColorMagenta = ""
	ColorCyan = ""
	ColorWhite = ""
	ColorBold = ""
	ColorDim = ""

	EmojiSuccess = "✓"
	EmojiError = "x"
	EmojiWarning = "!"
	EmojiInfo = "i"
	EmojiRocket = ">"
	EmojiGear = "*"
	EmojiFile = "-"
	EmojiCheck = "+"
}
