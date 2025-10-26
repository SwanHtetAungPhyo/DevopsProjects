package executor

import (
	"strings"

	"github.com/SCL/internal/parser"
)

// ModeDetector extracts the mode and setting variables from SCL file
type ModeDetector struct {
	*parser.BaseInfraDSLListener
	mode    string
	setting string
}

func NewModeDetector() *ModeDetector {
	return &ModeDetector{
		mode:    "compile",       // default mode
		setting: "configuration", // default setting
	}
}

func (m *ModeDetector) GetMode() string {
	return m.mode
}

func (m *ModeDetector) GetSetting() string {
	return m.setting
}

func (m *ModeDetector) EnterAssignment(ctx *parser.AssignmentContext) {
	varName := ctx.IDENTIFIER().GetText()
	expr := ctx.Expression().GetText()
	cleanValue := m.cleanExpression(expr)

	if varName == "mode" {
		m.mode = cleanValue
	} else if varName == "setting" {
		m.setting = cleanValue
	}
}

func (m *ModeDetector) cleanExpression(expr string) string {
	expr = strings.TrimSpace(expr)
	expr = strings.Trim(expr, "\"")
	return expr
}
