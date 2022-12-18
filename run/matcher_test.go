package run

import (
	"bytes"
	"errors"
	"testing"

	"github.com/prashantv/gostub"
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
			desc:    "MatchAny: match",
			matcher: MatchAny,
			cmd:     &Cmd{},
			matches: true,
		},

		{
			desc: "MatchAll: match",
			matcher: MatchAll(
				MatchString("/bin/echo"),
				MatchStdin("howdy"),
			),
			cmd: func() *Cmd {
				cmd := client.Command("/bin/echo")
				cmd.Stdin = bytes.NewBufferString("howdy")
				return cmd
			}(),
			matches: true,
		},
		{
			desc: "MatchAll: no match",
			matcher: MatchAll(
				MatchString("/bin/true"),
				MatchStdin("howdy"),
			),
			cmd: func() *Cmd {
				cmd := client.Command("/bin/echo")
				cmd.Stdin = bytes.NewBufferString("howdy")
				return cmd
			}(),
			matches: false,
		},

		{
			desc:    "MatchStdin: match",
			matcher: MatchStdin("howdy"),
			cmd: func() *Cmd {
				cmd := client.Command("/bin/echo")
				cmd.Stdin = bytes.NewBufferString("howdy")
				return cmd
			}(),
			matches: true,
		},
		{
			desc:    "MatchStdin: no match",
			matcher: MatchStdin("howdy"),
			cmd: func() *Cmd {
				cmd := client.Command("/bin/echo")
				cmd.Stdin = bytes.NewBufferString("nope")
				return cmd
			}(),
			matches: false,
		},
		{
			desc:    "MatchStdin: no stdin",
			matcher: MatchStdin("howdy"),
			cmd:     client.Command("/bin/echo"),
			matches: false,
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

func TestMatcher_ioReadAllError(t *testing.T) {
	stubs := gostub.StubFunc(&ioReadAll, nil, errors.New("boom"))
	defer stubs.Reset()

	matcher := MatchStdin("howdy")

	cmd := NewClient().Command("/bin/echo")
	cmd.Stdin = bytes.NewBufferString("")

	assert.PanicsWithError(t, "boom", func() {
		matcher(cmd)
	})
}
