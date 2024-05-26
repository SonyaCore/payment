package validators

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

// DescriptionValidator are a custom validation function for description field
func DescriptionValidator(fl validator.FieldLevel) bool {
	// Regular expression to match only alphanumeric characters, spaces, and common punctuation
	regex := regexp.MustCompile(`^[a-zA-Z0-9\s,.!?'-]+$`)
	description := fl.Field().String()
	return regex.MatchString(description)
}
