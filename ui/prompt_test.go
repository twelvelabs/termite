package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPromptType_String(t *testing.T) {
	assert.Equal(t, "Unknown", PromptTypeUnknown.String())
	assert.Equal(t, "Confirm", PromptTypeConfirm.String())
	assert.Equal(t, "Input", PromptTypeInput.String())
	assert.Equal(t, "MultiSelect", PromptTypeMultiSelect.String())
	assert.Equal(t, "Select", PromptTypeSelect.String())
}
