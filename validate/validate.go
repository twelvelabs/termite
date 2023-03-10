// Package validate is a light wrapper around the [validator] library
// that configures it for english translation.
//
// [validator]: https://github.com/go-playground/validator
package validate

import (
	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	validators "github.com/go-playground/validator/v10/non-standard/validators"
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

	RegisterValidation("kebabcase", "{0} must be kebabcase", isKebabCase)
	RegisterValidation("not-blank", "{0} must not be blank", validators.NotBlank)
	RegisterValidation("posix-mode", "{0} must be a valid posix file mode", isPosixMode)
}

// RegisterValidation adds a custom validation function and translation message
// for the given tag name.
func RegisterValidation(tag string, translation string, fn validator.Func) {
	err := validate.RegisterValidation(tag, fn)
	if err != nil {
		panic(err)
	}

	err = validate.RegisterTranslation(
		tag,
		translator,
		func(ut ut.Translator) error {
			return ut.Add(tag, translation, true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(tag, fe.Field())
			return t
		},
	)
	if err != nil {
		panic(err)
	}
}

// Struct validates a struct using tag style validation.
func Struct(data any) error {
	return translateErr("", validate.Struct(data))
}

// Var validates a single value using rules.
func Var(value any, rules string) error {
	return translateErr("", validate.Var(value, rules))
}

// KeyVal validates the key/value pair using rules.
func KeyVal(key string, value any, rules string) (err error) {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			err = panicToErr(key, rules, panicVal)
		}
	}()
	err = translateErr(key, validate.Var(value, rules))
	return err
}
