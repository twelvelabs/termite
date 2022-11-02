// Package validate is a light wrapper around the [validator] library
// that configures it for english translation.
//
// [validator]: https://github.com/go-playground/validator
package validate

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	uni        *ut.UniversalTranslator
	translator ut.Translator
	validate   *validator.Validate
)

func init() {
	en := en.New()
	uni = ut.New(en, en)
	translator, _ = uni.GetTranslator("en")

	validate = validator.New()
	_ = translations.RegisterDefaultTranslations(validate, translator)
}

// Struct validates a struct using tag style validation.
func Struct(data any) error {
	return translateErr(validate.Struct(data))
}

// Var validates a single variable using tag style validation.
func Var(value any, rules string) error {
	return translateErr(validate.Var(value, rules))
}

func translateErr(err error) error {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		lines := []string{}
		for _, v := range errs.Translate(translator) {
			lines = append(lines, strings.TrimSpace(v))
		}
		sort.Strings(lines)
		err = fmt.Errorf(strings.Join(lines, ", "))
	}
	return err
}
