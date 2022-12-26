package run

import (
	"fmt"
	"os/exec"
)

// NewExitError returns a new exit error for code.
func NewExitError(code int) *ExitError {
	return &ExitError{
		ExitError: &exec.ExitError{},
		code:      code,
	}
}

// ExitError is a light wrapper for [exec.ExitError] that allows for easily accessing
// and stubbing the exit code.
//
// The reason we don't use [exec.ExitError] directly is that it stores the
// exit code inside `os.ProcessState`, which in turn uses `syscall.WaitStatus`,
// and WaitStatus is implemented differently per-system.
type ExitError struct {
	*exec.ExitError
	code int
}

func (e *ExitError) Unwrap() error {
	return e.ExitError
}

// Code returns the exit code.
func (e *ExitError) Code() int {
	return e.code
}

// Error returns the error message.
func (e *ExitError) Error() string {
	return fmt.Sprintf("exit status %v", e.Code())
}
