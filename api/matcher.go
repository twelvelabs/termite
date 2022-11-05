package api

import (
	"net/http"
	"net/url"
	"strings"
)

type Matcher func(req *http.Request) bool

// MatchAny is a matcher that matches any request.
func MatchAny(*http.Request) bool {
	return true
}

func MatchDelete(path string) Matcher {
	return MatchRequest(http.MethodDelete, path)
}

func MatchGet(path string) Matcher {
	return MatchRequest(http.MethodGet, path)
}

func MatchPatch(path string) Matcher {
	return MatchRequest(http.MethodPatch, path)
}

func MatchPost(path string) Matcher {
	return MatchRequest(http.MethodPost, path)
}

func MatchPut(path string) Matcher {
	return MatchRequest(http.MethodPut, path)
}

// MatchRequest creates a matcher that matches on method and path.
func MatchRequest(method string, path string) Matcher {
	return func(req *http.Request) bool {
		if !strings.EqualFold(req.Method, method) {
			return false
		}
		return req.URL.EscapedPath() == path
	}
}

// MatchRequestQuery creates a matcher that matches on method, path, and query params.
func MatchRequestQuery(method string, path string, query url.Values) Matcher {
	return func(req *http.Request) bool {
		if !MatchRequest(method, path)(req) {
			return false
		}
		actualQuery := req.URL.Query()
		for param := range query {
			if !(actualQuery.Get(param) == query.Get(param)) {
				return false
			}
		}
		return true
	}
}
