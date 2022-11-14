package run

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatcher(t *testing.T) {
	client := NewClient()
	tests := []struct {
		desc    string
		matcher Matcher
		cmd     *Cmd
		matches bool
	}{
		{
			desc:    "Any",
			matcher: MatchAny,
			cmd:     &Cmd{},
			matches: true,
		},

		{
			desc:    "MatchString: match",
			matcher: MatchString("/bin/echo foo"),
			cmd:     client.Command("/bin/echo", "foo"),
			matches: true,
		},
		{
			desc:    "MatchString: no match",
			matcher: MatchString("/bin/echo foo"),
			cmd:     client.Command("/bin/echo", "bar"),
			matches: false,
		},

		{
			desc:    "MatchRegexp: match",
			matcher: MatchRegexp(`echo\s+(\w+)$`),
			cmd:     client.Command("/bin/echo", "foo"),
			matches: true,
		},
		{
			desc:    "MatchRegexp: no match",
			matcher: MatchRegexp(`echo\s+(\w+)$`),
			cmd:     client.Command("/bin/echo", "foo", "bar", "baz"),
			matches: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			actual := tt.matcher(tt.cmd)
			assert.Equal(t, tt.matches, actual)
		})
	}
}
