package cmd

import (
	"fmt"
	"os"

	"github.com/SCL/internal/color"
	"github.com/SCL/internal/utils"
	"github.com/spf13/cobra"
)

var (
	// Flags
	outputFile string
	verbose    bool
	noColor    bool
	version    = "1.0.0"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "scl [input-file]",
	Short: "Infrastructure DSL Compiler",
	Long: fmt.Sprintf(`%s%s
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║   🚀  Infrastructure DSL Compiler                        	║
║   ⚙   Declarative → Bash Translation Engine             	║
║   Version: %s                                            	║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
%s
A powerful compiler that translates declarative infrastructure
configuration into executable bash scripts.

Examples:
  scl deploy.scl
  scl deploy.scl -o production.sh
  scl deploy.scl --output staging.sh --verbose
`, color.ColorCyan, color.ColorBold, version, color.ColorReset),
	Args: cobra.MinimumNArgs(1),
	Run:  compile,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.PrintError(fmt.Sprintf("%v", err))
		os.Exit(1)
	}
}

func init() {
	// Flags
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "output.sh", "Output bash script file")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colored output")

	// Add version flag
	rootCmd.Version = version
}
