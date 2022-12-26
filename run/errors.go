package run

import (
	"fmt"
)

// NewExitError returns a new exit error for code.
func NewExitError(code int) *ExitError {
	return &ExitError{code: code}
}

// ExitError is a replacement for [exec.ExitError] used when stubbing that
// allows for accessing the exit code.
//
// The reason we don't use [exec.ExitError] directly is that it stores the
// exit code inside `os.ProcessState`, which in turn uses `syscall.WaitStatus`,
// and WaitStatus is implemented differently per-system.
type ExitError struct {
	Stderr []byte
	code   int
}

// Code returns the exit code.
func (e *ExitError) Code() int {
	return e.code
}

// Error returns the error message.
func (e *ExitError) Error() string {
	return fmt.Sprintf("exit status %v", e.Code())
}
