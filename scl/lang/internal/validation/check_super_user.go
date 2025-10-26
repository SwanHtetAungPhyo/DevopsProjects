package validation

import "fmt"

func (v *RequiredFieldValidator) checkSuperUser() {
	field := "super_user"
	allowed := []string{"true", "false"}

	value, exists := v.variables[field]
	if !exists {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Allowed: allowed,
			Message: fmt.Sprintf("Missing required field '%s'", field),
		})
		return
	}

	if !contains(allowed, value) {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   value,
			Allowed: allowed,
			Message: fmt.Sprintf("Invalid value '%s' for '%s'", value, field),
		})
	}
}
