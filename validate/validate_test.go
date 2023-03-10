package validate

import (
	"testing"

	validator "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type widget struct {
	Name   string `validate:"required"`
	Status string `validate:"oneof=active inactive"`
}

func TestRegisterTranslation(t *testing.T) {
	err := KeyVal("key", "val", "some-rule")
	assert.ErrorContains(t, err, "undefined rule")

	isAlwaysValid := func(fl validator.FieldLevel) bool {
		return true
	}

	RegisterValidation("some-rule", "{0} must be something", isAlwaysValid)
	err = KeyVal("key", "val", "some-rule")
	assert.NoError(t, err)

	// Validation tag must be non-empty.
	assert.Panics(t, func() {
		RegisterValidation("", "{0} must be something", isAlwaysValid)
	})
	// Placeholder brackets must be matching.
	assert.Panics(t, func() {
		RegisterValidation("some-rule", "{ must be something", isAlwaysValid)
	})
}

func TestStruct(t *testing.T) {
	assert.NoError(t, Struct(&widget{
		Name:   "untitled",
		Status: "active",
	}))
	assert.ErrorContains(t, Struct(&widget{
		Name:   "",
		Status: "unknown",
	}), "Name is a required field, Status must be one of [active inactive]")
}

func TestVar(t *testing.T) {
	assert.NoError(t, Var(5, "gt=0"))
	assert.ErrorContains(t, Var(0, "gt=0"), "must be greater than 0")
}

func TestKeyVal(t *testing.T) {
	var tests = []struct {
		Rule  string
		Value any
		Err   string
	}{
		{
			Rule:  "",
			Value: "",
			Err:   "",
		},
		{
			Rule:  "unknown",
			Value: "",
			Err:   "undefined rule [my-val: unknown]",
		},
		{
			Rule:  "lowercase",
			Value: false,
			Err:   "invalid rule for bool [my-val: lowercase]",
		},
		{
			Rule:  "kebabcase",
			Value: "not_kebab_case",
			Err:   "my-val must be kebabcase",
		},
		{
			Rule:  "kebabcase",
			Value: "is-kebab-case",
			Err:   "",
		},
		{
			Rule:  "not-blank",
			Value: " ",
			Err:   "my-val must not be blank",
		},
		{
			Rule:  "not-blank",
			Value: "is-not-blank",
			Err:   "",
		},
		{
			Rule:  "posix-mode",
			Value: "12345",
			Err:   "my-val must be a valid posix file mode",
		},
		{
			Rule:  "posix-mode",
			Value: "0755",
			Err:   "",
		},
		{
			Rule:  "posix-mode",
			Value: "755",
			Err:   "",
		},
		{
			Rule:  "gt=0,lte=10",
			Value: 12,
			Err:   "my-val must be 10 or less",
		},
		{
			Rule:  "gt=0,lte=10",
			Value: 5,
			Err:   "",
		},
	}

	for _, test := range tests {
		t.Run(test.Rule, func(t *testing.T) {
			err := KeyVal("my-val", test.Value, test.Rule)

			if test.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.Err)
			}
		})
	}
}
