package api

import (
	"fmt"
	"net/http"
	"sync"
)

// Stub is a stubbed HTTP response for a specific match pattern.
type Stub struct {
	Matched   bool
	Matcher   Matcher
	Responder Responder
}

func NewStubbedTransport() *StubbedTransport {
	return &StubbedTransport{
		Requests: []*http.Request{},
		stubs:    []*Stub{},
	}
}

// StubbedTransport is a [net/http.RoundTripper] that serves stubbed responses.
type StubbedTransport struct {
	Requests []*http.Request

	mu    sync.Mutex
	stubs []*Stub
}

// RegisterStub registers a new stub for the given matcher/responder pair.
func (t *StubbedTransport) RegisterStub(matcher Matcher, responder Responder) *StubbedTransport {
	t.stubs = append(t.stubs, &Stub{
		Matcher:   matcher,
		Responder: responder,
	})
	return t
}

// RoundTrip implements the RoundTripper interface.
// Will attempt to match a registered stub or return an error if none found.
func (t *StubbedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.mu.Lock()
	var stub *Stub
	var matches []*Stub

	for _, s := range t.stubs {
		if s.Matcher(req) {
			matches = append(matches, s)
			if s.Matched {
				continue
			}
			if stub == nil {
				s.Matched = true
				stub = s
			}
		}
	}

	if stub == nil {
		t.mu.Unlock()
		n := len(matches)
		if n == 0 {
			return nil, fmt.Errorf("no registered stubs matching: %v", req)
		} else {
			return nil, fmt.Errorf("wanted %d of only %d stubs matching: %v", n+1, n, req)
		}
	}

	t.Requests = append(t.Requests, req)
	t.mu.Unlock()

	return stub.Responder(req)
}

// VerifyStubs fails the test if there are unmatched stubs.
func (t *StubbedTransport) VerifyStubs(test testable) {
	test.Helper()

	n := 0
	for _, s := range t.stubs {
		if !s.Matched {
			n++
		}
	}
	if n > 0 {
		test.Errorf("found %d unmatched stub(s)", n)
	}
}

type testable interface {
	Errorf(string, ...interface{})
	Helper()
}
