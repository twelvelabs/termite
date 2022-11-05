package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHeadersInterceptor(t *testing.T) {
	baseURL := "https://example.com"
	token := "TOKEN"
	transport := NewStubbedTransport()

	// Explodes when invalid/empty base URL
	assert.PanicsWithValue(t, "invalid BaseURL: ''", func() {
		_ = newHeadersInterceptor("", token, map[string]string{}, transport)
	})

	// Returns the transport unwrapped if no headers to add (empty headers, no auth token)
	rt := newHeadersInterceptor(baseURL, "", map[string]string{}, transport)
	assert.Equal(t, transport, rt)

	// Otherwise returns a new interceptor wrapping transport
	rt = newHeadersInterceptor(baseURL, token, map[string]string{}, transport)
	hi := rt.(*headersInterceptor)
	assert.Equal(t, "example.com", hi.host)
	assert.Equal(t, map[string]string{"Authorization": "Bearer TOKEN"}, hi.headers)
	assert.Equal(t, transport, hi.wrapped)
}

func TestHeadersInterceptor_RoundTrip(t *testing.T) {
	baseURL := "https://example.com"
	token := "TOKEN"
	headers := map[string]string{
		"X-Foo": "bar",
	}
	transport := NewStubbedTransport().
		RegisterStub(
			MatchGet("/aaa"),
			WithRequestHeaders(StringResponse("")),
		).
		RegisterStub(
			MatchGet("/bbb"),
			WithRequestHeaders(StringResponse("")),
		).
		RegisterStub(
			MatchGet("/ccc"),
			WithRequestHeaders(StringResponse("")),
		)

	rt := newHeadersInterceptor(baseURL, token, headers, transport)

	// Normal use case:
	// - Auth header should be added when request URL matches baseURL
	// - Other headers should be added when not already set in request
	req, _ := http.NewRequest(http.MethodGet, "https://example.com/aaa", nil)
	resp, err := rt.RoundTrip(req)
	assert.NoError(t, err)
	assert.Equal(t, "Bearer TOKEN", resp.Header.Get("Authorization"))
	assert.Equal(t, "bar", resp.Header.Get("X-Foo"))

	// Request headers already exist:
	req, _ = http.NewRequest(http.MethodGet, "https://example.com/bbb", nil)
	req.Header.Set("Authorization", "ALT_TOKEN")
	req.Header.Set("X-Foo", "123")
	resp, err = rt.RoundTrip(req)
	assert.NoError(t, err)
	assert.Equal(t, "ALT_TOKEN", resp.Header.Get("Authorization"))
	assert.Equal(t, "123", resp.Header.Get("X-Foo"))

	// Cross domain (should not set auth header):
	req, _ = http.NewRequest(http.MethodGet, "https://other.com/ccc", nil)
	resp, err = rt.RoundTrip(req)
	assert.NoError(t, err)
	assert.Equal(t, "", resp.Header.Get("Authorization"))
	assert.Equal(t, "bar", resp.Header.Get("X-Foo"))
}
