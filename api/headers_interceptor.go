package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	hAccept          = "Accept"
	hAuthorization   = "Authorization"
	hContentType     = "Content-Type"
	vAcceptJSON      = "application/json, text/*;q=0.9, */*;q=0.8"
	vContentTypeJSON = "application/json; charset=utf-8"
)

func newHeadersInterceptor(baseURL string, authToken string, headers map[string]string, rt http.RoundTripper) http.RoundTripper {
	if _, ok := headers[hAuthorization]; !ok && authToken != "" {
		headers[hAuthorization] = fmt.Sprintf("Bearer %s", authToken)
	}
	if len(headers) == 0 {
		return rt
	}
	url, err := url.ParseRequestURI(baseURL)
	if err != nil {
		panic(fmt.Sprintf("invalid BaseURL: '%v'", baseURL))
	}
	return &headersInterceptor{
		host:    url.Hostname(),
		headers: headers,
		wrapped: rt,
	}
}

type headersInterceptor struct {
	headers map[string]string
	host    string
	wrapped http.RoundTripper
}

func (hi *headersInterceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range hi.headers {
		// Only add the auth header when making requests to the host configured in ClientOptions.
		if k == hAuthorization && !isSameOrSubDomain(req.URL.Hostname(), hi.host) {
			continue
		}
		// Don't overwrite headers that have already been explicitly set.
		if req.Header.Get(k) != "" {
			continue
		}
		req.Header.Set(k, v)
	}

	return hi.wrapped.RoundTrip(req)
}

func isSameOrSubDomain(hostname string, domain string) bool {
	return isSameDomain(hostname, domain) || isSubDomain(hostname, domain)
}

func isSameDomain(hostname string, domain string) bool {
	return strings.EqualFold(hostname, domain)
}

func isSubDomain(hostname string, domain string) bool {
	return strings.HasSuffix(strings.ToLower(hostname), "."+strings.ToLower(domain))
}
