package run

import (
	"fmt"
	"regexp"
)

// Responder is a function that returns stubbed command output.
type Responder func(cmd *Cmd) ([]byte, error)

// ErrorResponse creates a responder that returns err.
func ErrorResponse(err error) Responder {
	return func(cmd *Cmd) ([]byte, error) {
		return nil, err
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
	return func(cmd *Cmd) ([]byte, error) {
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
		return []byte(matches[index]), nil
	}
}

// StringResponse creates a responder that returns s.
func StringResponse(s string) Responder {
	return func(cmd *Cmd) ([]byte, error) {
		return []byte(s), nil
	}
}
