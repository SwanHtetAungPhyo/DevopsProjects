package validation

type RequiredFieldValidator struct {
	variables map[string]string
	errors    []ValidationError
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Value   string
	Allowed []string
	Message string
}

func NewRequiredFieldValidator(vars map[string]string) *RequiredFieldValidator {
	return &RequiredFieldValidator{
		variables: vars,
		errors:    []ValidationError{},
	}
}

func (v *RequiredFieldValidator) Validate() error {
	v.checkSetting()
	v.checkSuperUser()
	v.checkOnError()

	if len(v.errors) > 0 {
		return v.formatErrors()
	}
	return nil
}
