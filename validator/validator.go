package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(s any) error {
	return validate.Struct(s)
}

func formatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			switch e.Tag() {
			case "required":
				errors[field] = field + " is required"
			case "min":
				errors[field] = field + " must be at least " + e.Param() + " characters"
			case "max":
				errors[field] = field + " must be at most " + e.Param() + " characters"
			case "email":
				errors[field] = field + " must be a valid email"
			default:
				errors[field] = field + " is invalid (" + e.Tag() + ") (" + e.Param() + ")"
			}
		}
	}

	return errors
}

func FormatValidationErrors(err error) map[string]string {
	return formatValidationErrors(err)
}
