package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPromptParams(t *testing.T) {
	params := GetPromptParams(
		WithHelp("This is a required field."),
		WithValidation("required"),
	)
	assert.Equal(t, "This is a required field.", params.Help)
	assert.Equal(t, "required", params.ValidationRules)
}
