package run

import (
	"regexp"
)

// Matcher is function that matches commands.
type Matcher func(cmd *Cmd) bool

// MatchAny is a matcher that matches any command.
func MatchAny(cmd *Cmd) bool {
	return true
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
