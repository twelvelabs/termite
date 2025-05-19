package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var (
	ioReadAll = io.ReadAll
)

func NewRESTClient(opts *ClientOptions) *RESTClient {
	url, err := url.ParseRequestURI(opts.BaseURL)
	if err != nil {
		panic(fmt.Sprintf("invalid BaseURL: '%v'", opts.BaseURL))
	}
	if !url.IsAbs() {
		panic(fmt.Sprintf("invalid BaseURL: '%v'", opts.BaseURL))
	}
	opts.BaseURL = url.String()

	if opts.Headers == nil {
		opts.Headers = map[string]string{}
	}

	// Automatically set Accept and Content-Type headers if empty.
	if _, ok := opts.Headers[hAccept]; !ok {
		opts.Headers[hAccept] = vAcceptJSON
	}
	if _, ok := opts.Headers[hContentType]; !ok {
		opts.Headers[hContentType] = vContentTypeJSON
	}

	return &RESTClient{
		Client:  NewClientWith(opts),
		BaseURL: opts.BaseURL,
	}
}

// RESTClient is a client wrapper for HTTP requests that return JSON.
type RESTClient struct {
	*Client
	BaseURL string
}

type RESTClientError struct {
	HTTPResponse *http.Response
	Message      string
}

func (e *RESTClientError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.HTTPResponse.StatusCode, e.Message)
}

// RegisterStub registers a new stub for the given matcher/responder pair.
func (c *RESTClient) RegisterStub(matcher Matcher, responder Responder) *RESTClient {
	c.Client.RegisterStub(matcher, responder)
	return c
}

// WithStubbing configures stubbing and returns the receiver.
func (c *RESTClient) WithStubbing() *RESTClient {
	c.Client.WithStubbing()
	return c
}

// DoWithContext performs a HTTP request from the given args
// and parses the returned text as JSON into the response struct.
//
//   - url may be a path (relative to BaseURL), or an absolute URL.
//   - body may be either a JSON encodable struct, or an io.Reader.
//   - response may be a JSON encodable struct (or nil if unused).
//
// Note: Absolute URLs will only include the auth token if they share the same root domain.
func (c *RESTClient) DoWithContext(ctx context.Context, method string, url string, payload any, response any) error {
	// Coerce the payload into an io.Reader...
	body, ok := payload.(io.Reader)
	if !ok {
		body := &bytes.Buffer{}
		if err := json.NewEncoder(body).Encode(payload); err != nil {
			return err
		}
	}
	// Create an http.Request w/ it...
	req, err := http.NewRequestWithContext(ctx, method, c.abs(url), body)
	if err != nil {
		return err
	}
	// Then release the hounds™
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !success {
		return &RESTClientError{resp, "received unsuccessful response"}
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if response == nil {
		return nil
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	b, err := ioReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &response)
	if err != nil {
		return err
	}

	return nil
}

// Do performs a HTTP request from the given args and parses the response as JSON.
func (c *RESTClient) Do(method string, path string, body any, response any) error {
	return c.DoWithContext(context.Background(), method, path, body, response)
}

// Delete performs a HTTP Delete and parses the response JSON into resp.
func (c *RESTClient) Delete(path string, resp any) error {
	return c.Do(http.MethodDelete, path, nil, resp)
}

// Get performs a HTTP Get and parses the response JSON into resp.
func (c *RESTClient) Get(path string, resp any) error {
	return c.Do(http.MethodGet, path, nil, resp)
}

// Patch performs a HTTP Patch and parses the response JSON into resp.
func (c *RESTClient) Patch(path string, body any, resp any) error {
	return c.Do(http.MethodPatch, path, body, resp)
}

// Post performs a HTTP Post and parses the response JSON into resp.
func (c *RESTClient) Post(path string, body any, resp any) error {
	return c.Do(http.MethodPost, path, body, resp)
}

// Put performs a HTTP Put and parses the response JSON into resp.
func (c *RESTClient) Put(path string, body any, resp any) error {
	return c.Do(http.MethodPut, path, body, resp)
}

func (c *RESTClient) abs(path string) string {
	if strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "http://") {
		return path
	}
	return strings.TrimSuffix(c.BaseURL, "/") + "/" + strings.TrimPrefix(path, "/")
}
