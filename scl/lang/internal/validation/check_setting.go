package validation

import (
	"fmt"
	"strings"

	"github.com/SCL/internal/token"
)

func (v *RequiredFieldValidator) checkSetting() {
	fields := token.SETTING
	allowed := token.GetAllowedSettingKeys()

	value, exists := v.variables[fields]
	if !exists {
		v.errors = append(v.errors, ValidationError{
			Field:   fields,
			Allowed: allowed,
			Message: fmt.Sprintf("Ivalid SCL file format.\n Missing required field '%s'", fields),
		})
		return
	}
	if !contains(allowed, value) {
		v.errors = append(v.errors, ValidationError{
			Field:   fields,
			Value:   value,
			Allowed: allowed,
			Message: fmt.Sprintf("Invalid value '%s' for '%s'", value, fields),
		})
	}

}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func formatAllowed(values []string) string {
	quoted := make([]string, len(values))
	for i, v := range values {
		quoted[i] = fmt.Sprintf("'%s'", v)
	}
	return strings.Join(quoted, " or ")
}
