package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/SCL/internal/color"
	"github.com/SCL/internal/nodes"
	"github.com/SCL/internal/parser"
	"github.com/SCL/internal/utils"
	"github.com/SCL/internal/validation"
	"github.com/antlr4-go/antlr/v4"
	"github.com/spf13/cobra"
)

func compile(cmd *cobra.Command, args []string) {
	// Disable colors if requested
	if noColor {
		color.DisableColors()
	}

	utils.PrintBanner()

	inputFile := args[0]
	startTime := time.Now()

	// Step 1: Reading file
	utils.PrintStep("1/4", "Reading input file...")
	if verbose {
		utils.PrintInfo(fmt.Sprintf("Input: %s", inputFile))
	}

	input, err := antlr.NewFileStream(inputFile)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to read file: %v", err))
		os.Exit(1)
	}
	utils.PrintSuccess(fmt.Sprintf("Loaded: %s", inputFile))

	// Step 2: Lexical analysis
	utils.PrintStep("2/4", "Performing lexical analysis...")
	lexer := parser.NewInfraDSLLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)

	if verbose {
		utils.PrintInfo(fmt.Sprintf("Tokens generated: %d", stream.Size()))
	}
	utils.PrintSuccess("Tokenization complete")

	// Step 3: Parsing
	utils.PrintStep("3/4", "Parsing syntax tree...")
	p := parser.NewInfraDSLParser(stream)

	// Add error listener
	p.RemoveErrorListeners()
	errorListener := &utils.ErrorListener{}
	p.AddErrorListener(errorListener)

	// Parse the program
	tree := p.Program()

	// Check for parsing errors
	if errorListener.HasErrors {
		utils.PrintError("Parsing failed with errors")
		os.Exit(1)
	}
	utils.PrintSuccess("AST generated successfully")

	// Step 4: Detect mode and process accordingly
	utils.PrintStep("4/4", "Processing...")

	// First pass: extract variables to detect mode
	modeDetector := nodes.NewModeDetector()
	antlr.ParseTreeWalkerDefault.Walk(modeDetector, tree)

	mode := modeDetector.GetMode()
	if verbose {
		utils.PrintInfo(fmt.Sprintf("Detected mode: %s", mode))
	}

	if mode == "interpret" {
		// Execute directly (Ansible-like configuration management)
		utils.PrintInfo("Mode: interpret - Executing configuration management like Ansible...")
		executeDirectly(tree, startTime)
	} else {
		// Generate bash code (default: compile mode)
		generateBashCode(tree, startTime)
	}
}

func executeDirectly(tree antlr.ParseTree, startTime time.Time) {
	exec := nodes.NewDirectExecutor()
	exec.SetVerbose(verbose)

	antlr.ParseTreeWalkerDefault.Walk(exec, tree)

	target := exec.GetTarget()
	if target != "" {
		name, _ := os.Hostname()
		utils.PrintInfo(fmt.Sprintf("Running on local machine ( %s ) ....", name))
	} else {
		utils.PrintInfo(fmt.Sprintf("Checking SSH connectivity to target: %s", target))
		if err := exec.CheckSSHConnectivity(target); err != nil {
			utils.PrintError(fmt.Sprintf("SSH connectivity check failed: %v", err))
			os.Exit(1)
		}
		utils.PrintSuccess(fmt.Sprintf("✓ SSH connection to %s is working", target))
	}
	v := validation.NewRequiredFieldValidator(exec.GetVariables())
	if err := v.Validate(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	if err := exec.Execute(); err != nil {
		utils.PrintError(fmt.Sprintf("Execution failed: %v", err))
		os.Exit(1)
	}

	duration := time.Since(startTime)

	utils.PrintSuccess("Execution completed successfully")

	utils.PrintStats(
		len(exec.GetImports()),
		exec.GetFunctionCount(),
		exec.GetCommandCount(),
		duration,
	)

	fmt.Printf("\n%s%s%s Deployment complete!%s\n\n",
		color.ColorBold, color.ColorGreen, color.EmojiRocket, color.ColorReset)
}

func generateBashCode(tree antlr.ParseTree, startTime time.Time) {
	utils.PrintInfo("Mode: compile - Generating bash code...")
	if verbose {
		utils.PrintInfo(fmt.Sprintf("Output: %s", outputFile))
	}

	generator := nodes.NewBashCodeGenerator()
	generator.SetVerbose(verbose)
	antlr.ParseTreeWalkerDefault.Walk(generator, tree)

	v := validation.NewRequiredFieldValidator(generator.GetVariables())
	if err := v.Validate(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	bashCode := generator.GetBashCode()

	// Write to output file
	err := os.WriteFile(outputFile, []byte(bashCode), 0755)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to write output: %v", err))
		os.Exit(1)
	}

	duration := time.Since(startTime)

	// Print final results
	utils.PrintSuccess(fmt.Sprintf("Compiled successfully → %s", outputFile))

	// Print statistics
	utils.PrintStats(
		len(generator.GetImports()),
		generator.GetFunctionCount(),
		generator.GetLineCount(),
		duration,
	)

	// Final message
	fmt.Printf("\n%s%s%s Ready to deploy!%s Run with: %sbash %s%s\n\n",
		color.ColorBold, color.ColorGreen, color.EmojiRocket, color.ColorReset,
		color.ColorYellow, outputFile, color.ColorReset)
}
