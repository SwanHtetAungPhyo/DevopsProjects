package cmd

import (
	"fmt"
	"runtime"

	"github.com/SCL/internal/color"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Display detailed version information about the SCL compiler",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVersion() {
	fmt.Printf("\n%s%s╔═══════════════════════════════════════════╗%s\n", color.ColorBold, color.ColorCyan, color.ColorReset)
	fmt.Printf("%s%s║  SCL Compiler - Version Information      ║%s\n", color.ColorBold, color.ColorCyan, color.ColorReset)
	fmt.Printf("%s%s╚═══════════════════════════════════════════╝%s\n\n", color.ColorBold, color.ColorCyan, color.ColorReset)

	fmt.Printf("  %sVersion:%s     %s%s%s\n", color.ColorCyan, color.ColorReset, color.ColorYellow, version, color.ColorReset)
	fmt.Printf("  %sGo Version:%s  %s%s%s\n", color.ColorCyan, color.ColorReset, color.ColorYellow, runtime.Version(), color.ColorReset)
	fmt.Printf("  %sOS/Arch:%s     %s%s/%s%s\n", color.ColorCyan, color.ColorReset, color.ColorYellow, runtime.GOOS, runtime.GOARCH, color.ColorReset)
	fmt.Printf("  %sCompiler:%s    %s%s%s\n\n", color.ColorCyan, color.ColorReset, color.ColorYellow, runtime.Compiler, color.ColorReset)
}
