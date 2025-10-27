package utils

import (
	"fmt"

	"github.com/SCL/internal/color"
	"github.com/antlr4-go/antlr/v4"
)

// ErrorListener for better error reporting
type ErrorListener struct {
	*antlr.DefaultErrorListener
	HasErrors bool
}

func (l *ErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	l.HasErrors = true
	PrintError(fmt.Sprintf("Syntax error at line %d:%d", line, column))
	fmt.Printf("  %s%s%s\n", color.ColorDim, msg, color.ColorReset)
}
