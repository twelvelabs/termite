package api

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func parseQuery(s string) url.Values {
	v, err := url.ParseQuery(s)
	if err != nil {
		panic(err)
	}
	return v
}

func parseURL(s string) *url.URL {
	url, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return url
}

func TestMatcher(t *testing.T) {
	tests := []struct {
		desc    string
		matcher Matcher
		request *http.Request
		matches bool
	}{
		{
			desc:    "Any",
			matcher: MatchAny,
			request: &http.Request{
				Method: http.MethodGet,
				URL:    parseURL("http://example.com/yolo"),
			},
			matches: true,
		},

		{
			desc:    "Delete",
			matcher: MatchDelete("/foo/bar"),
			request: &http.Request{
				Method: http.MethodDelete,
				URL:    parseURL("http://example.com/foo/bar"),
			},
			matches: true,
		},
		{
			desc:    "Get",
			matcher: MatchGet("/foo/bar"),
			request: &http.Request{
				Method: http.MethodGet,
				URL:    parseURL("http://example.com/foo/bar"),
			},
			matches: true,
		},
		{
			desc:    "Patch",
			matcher: MatchPatch("/foo/bar"),
			request: &http.Request{
				Method: http.MethodPatch,
				URL:    parseURL("http://example.com/foo/bar"),
			},
			matches: true,
		},
		{
			desc:    "Post",
			matcher: MatchPost("/foo/bar"),
			request: &http.Request{
				Method: http.MethodPost,
				URL:    parseURL("http://example.com/foo/bar"),
			},
			matches: true,
		},
		{
			desc:    "Put",
			matcher: MatchPut("/foo/bar"),
			request: &http.Request{
				Method: http.MethodPut,
				URL:    parseURL("http://example.com/foo/bar"),
			},
			matches: true,
		},

		{
			desc:    "Request: matching method and path",
			matcher: MatchRequest(http.MethodGet, "/foo/bar"),
			request: &http.Request{
				Method: http.MethodGet,
				URL:    parseURL("http://example.com/foo/bar"),
			},
			matches: true,
		},
		{
			desc:    "Request: not matching method",
			matcher: MatchRequest(http.MethodPost, "/foo/bar"),
			request: &http.Request{
				Method: http.MethodGet,
				URL:    parseURL("http://example.com/foo/bar"),
			},
			matches: false,
		},
		{
			desc:    "Request: not matching path",
			matcher: MatchRequest(http.MethodGet, "/baz"),
			request: &http.Request{
				Method: http.MethodGet,
				URL:    parseURL("http://example.com/foo/bar"),
			},
			matches: false,
		},

		{
			desc:    "RequestQuery: matching method path and query",
			matcher: MatchRequestQuery(http.MethodGet, "/foo/bar", parseQuery("sort=date&order=asc")),
			request: &http.Request{
				Method: http.MethodGet,
				URL:    parseURL("http://example.com/foo/bar?order=asc&x=1&sort=date"),
			},
			matches: true,
		},
		{
			desc:    "RequestQuery: not matching method",
			matcher: MatchRequestQuery(http.MethodPost, "/foo/bar", parseQuery("sort=date&order=asc")),
			request: &http.Request{
				Method: http.MethodGet,
				URL:    parseURL("http://example.com/foo/bar?order=asc&x=1&sort=date"),
			},
			matches: false,
		},

		{
			desc:    "RequestQuery: not matching query",
			matcher: MatchRequestQuery(http.MethodGet, "/foo/bar", parseQuery("sort=date&order=asc")),
			request: &http.Request{
				Method: http.MethodGet,
				URL:    parseURL("http://example.com/foo/bar?order=asc&x=1"),
			},
			matches: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			actual := tt.matcher(tt.request)
			assert.Equal(t, tt.matches, actual)
		})
	}
}
