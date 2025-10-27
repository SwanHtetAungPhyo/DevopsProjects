package cmd

import (
	"fmt"
	"os"

	"github.com/SCL/internal/parser"
	"github.com/SCL/internal/utils"
	"github.com/antlr4-go/antlr/v4"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check [input-file]",
	Short: "Check syntax without compiling",
	Long:  "Validate the syntax of an SCL file without generating bash output",
	Args:  cobra.ExactArgs(1),
	Run:   checkSyntax,
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func checkSyntax(cmd *cobra.Command, args []string) {
	inputFile := args[0]

	utils.PrintStep("1/2", "Reading file...")
	input, err := antlr.NewFileStream(inputFile)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to read file: %v", err))
		os.Exit(1)
	}
	utils.PrintSuccess(fmt.Sprintf("Loaded: %s", inputFile))

	utils.PrintStep("2/2", "Checking syntax...")
	lexer := parser.NewInfraDSLLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewInfraDSLParser(stream)

	p.RemoveErrorListeners()
	errorListener := &utils.ErrorListener{}
	p.AddErrorListener(errorListener)

	p.Program()

	if errorListener.HasErrors {
		utils.PrintError("Syntax check failed")
		os.Exit(1)
	}

	utils.PrintSuccess("✓ Syntax is valid!")
	fmt.Println()
}
