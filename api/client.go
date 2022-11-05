package api

import (
	"net/http"
)

// NewClient returns a new client configured with default options.
func NewClient() *Client {
	return NewClientWith(&ClientOptions{})
}

// NewClientWith returns a new Client configured with opts.
func NewClientWith(opts *ClientOptions) *Client {
	if opts.Headers == nil {
		opts.Headers = map[string]string{}
	}

	transport := http.DefaultTransport
	if opts.Transport != nil {
		transport = opts.Transport
	}
	transport = newHeadersInterceptor(opts.BaseURL, opts.AuthToken, opts.Headers, transport)

	client := &http.Client{
		Transport: transport,
		Timeout:   opts.Timeout,
	}

	return &Client{
		Client: client,
	}
}

// Client is a wrapper around [net/http.Client] that supports stubbing.
type Client struct {
	*http.Client
}

// IsStubbed returns true if the transport is configured for stubbing.
func (c *Client) IsStubbed() bool {
	_, ok := c.Transport.(*StubbedTransport)
	return ok
}

// RegisterStub registers a new stub for the given matcher/responder pair.
func (c *Client) RegisterStub(matcher Matcher, responder Responder) *Client {
	if !c.IsStubbed() {
		// Considered auto-enabling for friendlier DevExp,
		// but it's better to require them to be explicit.
		panic("must enable stubbing before registering stubs")
	}
	transport := c.Transport.(*StubbedTransport)
	transport.RegisterStub(matcher, responder)
	return c
}

// VerifyStubs fails the test if there are unmatched stubs.
func (c *Client) VerifyStubs(t testable) {
	t.Helper()
	transport := c.Transport.(*StubbedTransport)
	transport.VerifyStubs(t)
}

// WithStubbing configures stubbing and returns the receiver.
func (c *Client) WithStubbing() *Client {
	if !c.IsStubbed() {
		c.Transport = NewStubbedTransport()
	}
	return c
}
