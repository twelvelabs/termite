package run

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewExitError(t *testing.T) {
	err := NewExitError(0)
	assert.Error(t, err)

	// Should be a wrapper around std lib exit error
	var exitErr *exec.ExitError
	assert.True(t, errors.As(err, &exitErr))
}

func TestExitError_Code(t *testing.T) {
	err := NewExitError(12)
	assert.Equal(t, 12, err.Code())
}

func TestExitError_Error(t *testing.T) {
	err := NewExitError(12)
	assert.Equal(t, "exit status 12", err.Error())
}
