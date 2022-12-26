package run

import (
	"fmt"
	"regexp"
)

// Responder is a function that returns stubbed command output.
type Responder func(cmd *Cmd) (stdout []byte, stderr []byte, err error)

// ErrorResponse creates a responder that returns err.
func ErrorResponse(err error) Responder {
	return func(cmd *Cmd) ([]byte, []byte, error) {
		return nil, nil, err
	}
}

// RegexpResponse creates a responder that returns the match index
// for the given regular expression pattern.
// Panics if there is no match for index.
//
// For example:
//
//	responder := RegexpResponse(`echo (\w+)$`, 1)
//	buf, err := responder(client.Command("echo", "howdy"))
//	// buf == []byte("howdy")
//	// err == nil
func RegexpResponse(pattern string, index int) Responder {
	r := regexp.MustCompile(pattern)
	return func(cmd *Cmd) ([]byte, []byte, error) {
		cmdStr := cmd.String()
		matches := r.FindStringSubmatch(cmdStr)
		if index >= len(matches) {
			panic(
				fmt.Sprintf(
					"invalid match: cmd='%s', pattern='%s', index='%d'",
					cmdStr,
					pattern,
					index,
				),
			)
		}
		return []byte(matches[index]), nil, nil
	}
}

// StringResponse creates a responder that returns the given string via stdout.
func StringResponse(s string) Responder {
	return func(cmd *Cmd) ([]byte, []byte, error) {
		return []byte(s), nil, nil
	}
}

// StdoutResponse creates a responder that returns the given bytes via stdout.
// If the provided exit code is non-zero, then an exit error will also be returned.
func StdoutResponse(stdout []byte, code int) Responder {
	return func(cmd *Cmd) ([]byte, []byte, error) {
		return stdout, nil, errorForCode(code)
	}
}

// StderrResponse creates a responder that returns the given bytes via stderr.
// If the provided exit code is non-zero, then an exit error will also be returned.
func StderrResponse(stderr []byte, code int) Responder {
	return func(cmd *Cmd) ([]byte, []byte, error) {
		return nil, stderr, errorForCode(code)
	}
}

// MuxResponse creates a responder that returns values for both stdout and stderr.
// If the provided exit code is non-zero, then an exit error will also be returned.
func MuxResponse(stdout []byte, stderr []byte, code int) Responder {
	return func(cmd *Cmd) ([]byte, []byte, error) {
		return stdout, stderr, errorForCode(code)
	}
}

func errorForCode(code int) error {
	if code == 0 {
		return nil
	}
	return NewExitError(code)
}
