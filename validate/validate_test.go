package validate_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/twelvelabs/termite/validate"
)

type widget struct {
	Name   string `validate:"required"`
	Status string `validate:"oneof=active inactive"`
}

func ExampleStruct() {
	type Widget struct {
		Name string `validate:"required"`
	}

	widget := &Widget{
		Name: "",
	}
	err := validate.Struct(widget)

	fmt.Println(err.Error())
	// Output: Name is a required field
}

func TestStruct(t *testing.T) {
	assert.NoError(t, validate.Struct(&widget{
		Name:   "untitled",
		Status: "active",
	}))
	assert.ErrorContains(t, validate.Struct(&widget{
		Name:   "",
		Status: "unknown",
	}), "Name is a required field, Status must be one of [active inactive]")
}

func ExampleVar() {
	count := 0
	err := validate.Var(count, "gt=0")

	fmt.Println(err.Error())
	// Output: must be greater than 0
}

func TestVar(t *testing.T) {
	assert.NoError(t, validate.Var(5, "gt=0"))
	assert.ErrorContains(t, validate.Var(0, "gt=0"), "must be greater than 0")
}
