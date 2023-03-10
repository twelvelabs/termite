package validate

import (
	"regexp"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

var (
	validKebabRegex     = regexp.MustCompile(`^[a-z\-\d]+$`)
	validPosixModeRegex = regexp.MustCompile(`^[0-7]?[0-7]{3}$`) // https://rubular.com/r/WY16zVRPA90l2K
)

// Custom validator that ensures a POSIX permission number
// (either octal or int form).
func isPosixMode(fl validator.FieldLevel) bool {
	field := fl.Field()
	return validPosixModeRegex.MatchString(strings.TrimSpace(field.String()))
}

// Custom validator that ensures kebab-case.
func isKebabCase(fl validator.FieldLevel) bool {
	field := fl.Field()
	return validKebabRegex.MatchString(strings.TrimSpace(field.String()))
}
