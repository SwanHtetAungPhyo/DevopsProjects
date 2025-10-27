package utils

import (
	"fmt"
	"time"

	"github.com/SCL/internal/color"
)

// PrintBanner Colored print functions
func PrintBanner() {
	banner := `
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║   @  Infrastructure DSL Compiler                          ║
║   @   Declarative → Bash Translation Engine               ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
`
	fmt.Printf("%s%s%s%s\n", color.ColorBold, color.ColorCyan, banner, color.ColorReset)
}

func PrintSuccess(msg string) {
	fmt.Printf("%s%s %s%s %s\n", color.ColorGreen, color.EmojiSuccess, color.ColorBold, msg, color.ColorReset)
}

func PrintError(msg string) {
	fmt.Printf("%s%s %s%s %s\n", color.ColorRed, color.EmojiError, color.ColorBold, msg, color.ColorReset)
}

func PrintWarning(msg string) {
	fmt.Printf("%s%s %s%s %s\n", color.ColorYellow, color.EmojiWarning, color.ColorBold, msg, color.ColorReset)
}

func PrintInfo(msg string) {
	fmt.Printf("%s%s%s %s%s\n", color.ColorDim, color.ColorCyan, color.EmojiInfo, msg, color.ColorReset)
}

func PrintStep(step, msg string) {
	fmt.Printf("%s[%s]%s %s\n", color.ColorBlue, step, color.ColorReset, msg)
}

func PrintStats(imports, functions, lines int, duration time.Duration) {
	fmt.Printf("\n%s%s═══════════════════════════════════════════════════════════%s\n", color.ColorBold, color.ColorCyan, color.ColorReset)
	fmt.Printf("%s%s Compilation Statistics:%s\n", color.ColorBold, color.ColorWhite, color.ColorReset)
	fmt.Printf("%s%s═══════════════════════════════════════════════════════════%s\n\n", color.ColorBold,color.ColorCyan, color.ColorReset)

	fmt.Printf("  %s%-20s%s %s%d%s\n", color.ColorCyan, "Imports:", color.ColorReset, color.ColorYellow, imports, color.ColorReset)
	fmt.Printf("  %s%-20s%s %s%d%s\n", color.ColorCyan, "Functions generated:", color.ColorReset, color.ColorYellow, functions, color.ColorReset)
	fmt.Printf("  %s%-20s%s %s%d%s\n", color.ColorCyan, "Lines of bash:", color.ColorReset, color.ColorYellow, lines,color.ColorReset)
	fmt.Printf("  %s%-20s%s %s%.2fms%s\n", color.ColorCyan, "Compilation time:", color.ColorReset, color.ColorYellow, float64(duration.Microseconds())/1000, color.ColorReset)

	fmt.Printf("\n%s%s═══════════════════════════════════════════════════════════%s\n", color.ColorBold, color.ColorCyan,color.ColorReset)
}
