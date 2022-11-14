package run

import (
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
