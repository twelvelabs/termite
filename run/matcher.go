package run

import (
	"io"
	"regexp"
)

var (
	// for stubbing
	ioReadAll = io.ReadAll
)

// Matcher is function that matches commands.
type Matcher func(cmd *Cmd) bool

// MatchAny is a matcher that matches any command.
func MatchAny(cmd *Cmd) bool {
	return true
}

// MatchAll returns a matcher that returns true if all matcher args match.
func MatchAll(matchers ...Matcher) Matcher {
	return func(cmd *Cmd) bool {
		for _, matcher := range matchers {
			if ok := matcher(cmd); !ok {
				return false
			}
		}
		return true
	}
}

// MatchStdin returns a matcher that matches against command stdin.
func MatchStdin(s string) Matcher {
	return func(cmd *Cmd) bool {
		if cmd.Stdin == nil {
			return false
		}
		buf, err := ioReadAll(cmd.Stdin)
		if err != nil {
			panic(err)
		}
		return s == string(buf)
	}
}

// MatchString returns a matcher that matches against command strings.
func MatchString(s string) Matcher {
	return func(cmd *Cmd) bool {
		return s == cmd.String()
	}
}

// MatchRegexp returns a matcher that matches against regular expressions.
func MatchRegexp(s string) Matcher {
	r := regexp.MustCompile(s)
	return func(cmd *Cmd) bool {
		return r.MatchString(cmd.String())
	}
}
