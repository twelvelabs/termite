package api

import (
	"net/http"
	"time"
)

// ClientOptions allow for configuring new clients.
type ClientOptions struct {
	// AuthToken is the authorization token included on requests made to BaseURL.
	AuthToken string

	// BaseURL is the base URL for relative API requests.
	BaseURL string

	// Headers are the HTTP headers that will be sent with every API request.
	Headers map[string]string

	// Timeout specifies a time limit for each API request.
	// Default is no timeout.
	Timeout time.Duration

	// Transport specifies the mechanism by which individual API requests are made.
	// Default is http.DefaultTransport.
	Transport http.RoundTripper
}
