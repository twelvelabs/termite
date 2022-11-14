package run

import (
	"context"
	"os/exec"
)

// NewClient returns a new Client.
func NewClient() *Client {
	return &Client{
		Executor: DefaultExecutor,
	}
}

// Client is an abstraction around [os/exec] to support stubbing.
type Client struct {
	Executor Executor
}

// Command returns the Cmd struct to execute the named program with
// the given arguments.
func (c *Client) Command(name string, arg ...string) *Cmd {
	cmd := exec.Command(name, arg...)
	return &Cmd{Cmd: cmd, client: c}
}

// CommandContext is like Command but includes a context.
func (c *Client) CommandContext(ctx context.Context, name string, arg ...string) *Cmd {
	cmd := exec.CommandContext(ctx, name, arg...)
	return &Cmd{Cmd: cmd, client: c}
}

// IsStubbed returns true if the executor is configured for stubbing.
func (c *Client) IsStubbed() bool {
	_, ok := c.Executor.(*StubExecutor)
	return ok
}

// RegisterStub registers a new stub for the given matcher/responder pair.
func (c *Client) RegisterStub(matcher Matcher, responder Responder) *Client {
	if !c.IsStubbed() {
		panic("must enable stubbing before registering stubs")
	}
	executor := c.Executor.(*StubExecutor)
	executor.RegisterStub(matcher, responder)
	return c
}

// VerifyStubs fails the test if there are unmatched stubs.
func (c *Client) VerifyStubs(t testable) {
	if !c.IsStubbed() {
		panic("must enable stubbing before verifying stubs")
	}
	t.Helper()
	executor := c.Executor.(*StubExecutor)
	executor.VerifyStubs(t)
}

// WithStubbing configures stubbing and returns the receiver.
func (c *Client) WithStubbing() *Client {
	if !c.IsStubbed() {
		c.Executor = NewStubExecutor()
	}
	return c
}
