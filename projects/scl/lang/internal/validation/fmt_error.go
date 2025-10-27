package validation

import (
	"fmt"
	"strings"

	"github.com/SCL/internal/color"
)

func (v *RequiredFieldValidator) formatErrors() error {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\n%s%s%s Configuration Validation Failed%s\n\n",
		color.ColorBold, color.ColorRed, color.EmojiError, color.ColorReset))

	sb.WriteString("Required fields are missing or invalid:\n\n")

	for _, err := range v.errors {
		sb.WriteString(fmt.Sprintf("  %s•%s %s%s:%s %s\n",
			color.ColorRed, color.ColorReset,
			color.ColorYellow, err.Field, color.ColorReset,
			err.Message))

		if err.Value != "" {
			sb.WriteString(fmt.Sprintf("    %sFound:%s %s'%s'%s\n",
				color.ColorDim, color.ColorReset,
				color.ColorRed, err.Value, color.ColorReset))
		}

		if len(err.Allowed) > 0 {
			sb.WriteString(fmt.Sprintf("    %sAllowed:%s %s%s%s\n",
				color.ColorDim, color.ColorReset,
				color.ColorGreen, formatAllowed(err.Allowed), color.ColorReset))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("%sExample:%s\n", color.ColorCyan, color.ColorReset))
	sb.WriteString(fmt.Sprintf("%s  setting := configuration;\n", color.ColorDim))
	sb.WriteString(fmt.Sprintf("  super_user = true;\n"))
	sb.WriteString(fmt.Sprintf("  on_error = rollback;%s\n\n", color.ColorReset))

	return fmt.Errorf("%s", sb.String())
}
