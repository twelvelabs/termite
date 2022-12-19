package run

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

// Cmd is a wrapper around [os/exec.Cmd] that supports stubbing.
// All functionality is delegated to [Client.Executor].
type Cmd struct {
	*exec.Cmd

	client *Client
}

// Output runs the command and returns its standard output.
func (c *Cmd) Output() ([]byte, error) {
	return c.client.Executor.Output(c)
}

// Run starts the specified command and waits for it to complete.
func (c *Cmd) Run() error {
	return c.client.Executor.Run(c)
}

// PeekStdin reads Stdin without moving the read offset to EOF.
func (c *Cmd) PeekStdin() ([]byte, error) {
	if c.Stdin == nil {
		return nil, nil
	}

	// io hijinks to reset the read offset
	copy := &bytes.Buffer{}
	reader := io.TeeReader(c.Stdin, copy)
	c.Stdin = copy

	return io.ReadAll(reader)
}

// DebugString the command string plus a truncated preview of stdin.
func (c *Cmd) DebugString() string {
	stdin, _ := c.PeekStdin()
	if len(stdin) > 20 {
		stdin = append(stdin[:20], []byte("…")...)
	}
	suffix := ""
	if len(stdin) > 0 {
		suffix = fmt.Sprintf(" [Stdin: \"%s\"]", string(stdin))
	}
	return c.String() + suffix
}
