package nodes

import (
	"fmt"
	"strings"

	"github.com/SCL/internal/parser"
	"github.com/SCL/internal/utils"
	"github.com/antlr4-go/antlr/v4"
)

// BashCodeGenerator generates bash code from the AST
type BashCodeGenerator struct {
	*parser.BaseInfraDSLListener
	bashCode      strings.Builder
	indent        int
	variables     map[string]string
	imports       map[string]bool
	functionCount int
	lineCount     int
	verbose       bool
}

func NewBashCodeGenerator() *BashCodeGenerator {
	gen := &BashCodeGenerator{
		variables: make(map[string]string),
		imports:   make(map[string]bool),
		verbose:   false,
	}
	gen.bashCode.WriteString("#!/bin/bash\n\n")
	gen.bashCode.WriteString("set -e  # Exit on error\n\n")
	return gen
}

// SetVerbose enables or disables verbose output
func (g *BashCodeGenerator) SetVerbose(v bool) {
	g.verbose = v
}

func (g *BashCodeGenerator) GetFunctionCount() int {
	return g.functionCount
}

func (g *BashCodeGenerator) GetLineCount() int {
	return g.lineCount
}

func (g *BashCodeGenerator) GetImports() map[string]bool {
	return g.imports
}

func (g *BashCodeGenerator) getIndent() string {
	return strings.Repeat("    ", g.indent)
}

func (g *BashCodeGenerator) EnterImportStatement(ctx *parser.ImportStatementContext) {
	moduleName := ctx.IDENTIFIER().GetText()
	g.imports[moduleName] = true
	g.bashCode.WriteString(fmt.Sprintf("# import %s\n", moduleName))

	if g.verbose {
		utils.PrintInfo(fmt.Sprintf("Importing module: %s", moduleName))
	}
}

func (g *BashCodeGenerator) EnterAssignment(ctx *parser.AssignmentContext) {
	varName := ctx.IDENTIFIER().GetText()
	expr := ctx.Expression().GetText()

	// Clean up the expression
	value := g.cleanExpression(expr)

	g.variables[varName] = value
	g.bashCode.WriteString(fmt.Sprintf("%s%s=%s\n",
		g.getIndent(), strings.ToUpper(varName), value))
	g.lineCount++

	if g.verbose {
		utils.PrintInfo(fmt.Sprintf("Assignment: %s = %s", varName, value))
	}
}

func (g *BashCodeGenerator) EnterDeclaration(ctx *parser.DeclarationContext) {
	varName := ctx.IDENTIFIER().GetText()
	typeName := ctx.Type_().GetText()
	expr := ctx.Expression().GetText()

	if g.verbose {
		utils.PrintInfo(fmt.Sprintf("Declaring variable: %s (%s)", varName, typeName))
	}

	// Handle special cases like check() function
	if strings.Contains(expr, "check(") {
		g.generateCheckCommand(varName, expr)
	} else {
		value := g.cleanExpression(expr)
		g.variables[varName] = value
		g.bashCode.WriteString(fmt.Sprintf("%s%s=%s\n",
			g.getIndent(), strings.ToUpper(varName), value))
	}
	g.lineCount++
}

func (g *BashCodeGenerator) generateCheckCommand(varName, expr string) {
	target := g.variables["target"]

	if g.verbose {
		utils.PrintInfo(fmt.Sprintf("Generating check command for: %s on target %s", varName, target))
	}

	g.bashCode.WriteString(g.getIndent())
	g.bashCode.WriteString(fmt.Sprintf("# Check if docker and python3 exist on %s\n", target))
	g.bashCode.WriteString(g.getIndent())
	g.bashCode.WriteString(fmt.Sprintf("ssh %s \"command -v docker && command -v python3\" > /dev/null 2>&1\n", target))
	g.bashCode.WriteString(g.getIndent())
	g.bashCode.WriteString(fmt.Sprintf("%s=$?\n", strings.ToUpper(varName)))
	g.lineCount += 3
}

func (g *BashCodeGenerator) EnterFunctionDeclaration(ctx *parser.FunctionDeclarationContext) {
	funcName := ctx.IDENTIFIER().GetText()
	g.bashCode.WriteString(fmt.Sprintf("\n%s%s() {\n", g.getIndent(), funcName))
	g.indent++
	g.functionCount++
	g.lineCount++

	utils.PrintSuccess(fmt.Sprintf("Generated function: %s()", funcName))

	if g.verbose {
		utils.PrintInfo(fmt.Sprintf("Function depth: %d", g.indent))
	}
}

func (g *BashCodeGenerator) ExitFunctionDeclaration(ctx *parser.FunctionDeclarationContext) {
	g.indent--
	g.bashCode.WriteString(fmt.Sprintf("%s}\n", g.getIndent()))
	g.lineCount++

	if g.verbose {
		funcName := ctx.IDENTIFIER().GetText()
		utils.PrintInfo(fmt.Sprintf("Closed function: %s()", funcName))
	}
}

// EnterExpressionStatement Handle function call statements (with semicolon)
func (g *BashCodeGenerator) EnterExpressionStatement(ctx *parser.ExpressionStatementContext) {
	qualifiedName := ctx.QualifiedName().GetText()

	if g.verbose {
		utils.PrintInfo(fmt.Sprintf("Expression statement: %s()", qualifiedName))
	}

	g.handleFunctionCall(qualifiedName, ctx.ArgumentList())
}

// EnterFunctionCallExpr Handle function calls in expressions
func (g *BashCodeGenerator) EnterFunctionCallExpr(ctx *parser.FunctionCallExprContext) {
	exprText := ctx.GetChild(0).(antlr.ParseTree).GetText()

	if g.verbose {
		utils.PrintInfo(fmt.Sprintf("Function call expression: %s()", exprText))
	}

	g.handleFunctionCall(exprText, ctx.ArgumentList())
}

// EnterMemberAccessExpr Handle member access (for primary.pkg_snap etc)
func (g *BashCodeGenerator) EnterMemberAccessExpr(ctx *parser.MemberAccessExprContext) {
	if g.verbose {
		memberName := ctx.GetText()
		utils.PrintInfo(fmt.Sprintf("Member access: %s", memberName))
	}
}

func (g *BashCodeGenerator) handleFunctionCall(qualifiedName string, argList parser.IArgumentListContext) {
	switch {
	case qualifiedName == "test":
		if g.verbose {
			utils.PrintInfo("Generating test access command")
		}
		g.bashCode.WriteString(g.getIndent())
		g.bashCode.WriteString("# Test access\n")
		g.lineCount++

	case qualifiedName == "print" || qualifiedName == "primary.print":
		if argList != nil {
			args := argList.GetText()
			cleanArgs := strings.Trim(args, "\"")

			if g.verbose {
				utils.PrintInfo(fmt.Sprintf("Generating print: %s", cleanArgs))
			}

			g.bashCode.WriteString(fmt.Sprintf("%secho \"%s\"\n", g.getIndent(), cleanArgs))
			g.lineCount++
		}

	case qualifiedName == "copy":
		if argList != nil {
			args := argList.GetText()
			parts := strings.Split(args, ",")
			if len(parts) >= 2 {
				source := strings.Trim(strings.TrimSpace(parts[0]), "\"")
				dest := strings.Trim(strings.TrimSpace(parts[1]), "\"")
				target := g.variables["target"]

				if g.verbose {
					utils.PrintInfo(fmt.Sprintf("Generating copy: %s → %s:%s", source, target, dest))
				}

				// Create destination directory first
				destDir := dest
				if !strings.HasSuffix(dest, "/") {
					lastSlash := strings.LastIndex(dest, "/")
					if lastSlash > 0 {
						destDir = dest[:lastSlash]
					}
				}

				g.bashCode.WriteString(fmt.Sprintf("%sssh %s \"mkdir -p %s\"\n",
					g.getIndent(), target, destDir))
				g.bashCode.WriteString(fmt.Sprintf("%sscp -r %s %s:%s\n",
					g.getIndent(), source, target, dest))
				g.lineCount += 2
			}
		}

	case qualifiedName == "create":
		if argList != nil {
			args := argList.GetText()
			parts := strings.Split(args, ",")
			if len(parts) >= 3 {
				destination := strings.Trim(strings.TrimSpace(parts[0]), "\"")
				filename := strings.Trim(strings.TrimSpace(parts[1]), "\"")
				mode := strings.Trim(strings.TrimSpace(parts[2]), "\"")
				target := g.variables["target"]
				fullPath := fmt.Sprintf("%s/%s", destination, filename)

				if g.verbose {
					utils.PrintInfo(fmt.Sprintf("Generating create: %s with mode %s on %s", fullPath, mode, target))
				}

				g.bashCode.WriteString(fmt.Sprintf("%sssh %s \"mkdir -p %s && touch %s && chmod %s %s\"\n",
					g.getIndent(), target, destination, fullPath, mode, fullPath))
				g.lineCount++
			}
		}

	case qualifiedName == "install":
		if argList != nil {
			args := argList.GetText()
			parts := strings.Split(args, ",")
			if len(parts) >= 2 {
				pkg := strings.Trim(strings.TrimSpace(parts[1]), "\"")
				target := g.variables["target"]

				if g.verbose {
					utils.PrintInfo(fmt.Sprintf("Generating install: %s on %s", pkg, target))
				}

				g.bashCode.WriteString(fmt.Sprintf("%sssh %s \"snap install %s\"\n",
					g.getIndent(), target, pkg))
				g.lineCount++
			}
		}

	case qualifiedName == "check":
		// Handled in declaration

	default:
		if !strings.Contains(qualifiedName, ".") {
			if g.verbose {
				utils.PrintInfo(fmt.Sprintf("Generating function call: %s()", qualifiedName))
			}

			g.bashCode.WriteString(fmt.Sprintf("%s%s\n", g.getIndent(), qualifiedName))
			g.lineCount++
		}
	}
}

func (g *BashCodeGenerator) EnterIfStatement(ctx *parser.IfStatementContext) {
	condition := ctx.Expression().GetText()
	bashCondition := g.convertCondition(condition)

	if g.verbose {
		utils.PrintInfo(fmt.Sprintf("Generating if statement: %s", condition))
	}

	g.bashCode.WriteString(fmt.Sprintf("%sif %s; then\n", g.getIndent(), bashCondition))
	g.indent++
	g.lineCount++
}

func (g *BashCodeGenerator) ExitIfStatement(ctx *parser.IfStatementContext) {
	g.indent--

	if len(ctx.AllBlock()) > 1 {
		if g.verbose {
			utils.PrintInfo("Generating else block")
		}

		g.bashCode.WriteString(fmt.Sprintf("%selse\n", g.getIndent()))
		g.indent++
		g.lineCount++
		g.indent--
	}

	g.bashCode.WriteString(fmt.Sprintf("%sfi\n", g.getIndent()))
	g.lineCount++

	if g.verbose {
		utils.PrintInfo("Closed if statement")
	}
}

func (g *BashCodeGenerator) convertCondition(condition string) string {
	varName := strings.ToUpper(condition)
	return fmt.Sprintf("[ $%s -eq 0 ]", varName)
}

func (g *BashCodeGenerator) cleanExpression(expr string) string {
	expr = strings.TrimSpace(expr)

	if strings.HasPrefix(expr, "[") && strings.HasSuffix(expr, "]") {
		content := strings.Trim(expr, "[]")
		parts := strings.Split(content, ",")
		if len(parts) > 0 {
			return strings.Trim(strings.TrimSpace(parts[0]), "\"")
		}
	}

	return strings.Trim(expr, "\"")
}

func (g *BashCodeGenerator) GetBashCode() string {
	if g.verbose {
		utils.PrintInfo("Finalizing bash code generation")
	}

	g.bashCode.WriteString("\n# Execute main function\n")
	g.bashCode.WriteString("main\n")
	g.lineCount += 2
	return g.bashCode.String()
}

func (g *BashCodeGenerator) GetVariables() map[string]string {
	return g.variables
}
