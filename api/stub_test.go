package api

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStubbedTransport_StubbingMethods(t *testing.T) {
	transport := NewStubbedTransport()

	transport.
		RegisterStub(
			MatchGet("/foo"),
			StringResponse(`{"foo": true}`),
		).
		RegisterStub(
			MatchGet("/foo"),
			StringResponse(`{"foo": false}`),
		)

	_, err := transport.RoundTrip(&http.Request{
		Method: http.MethodGet,
		URL:    parseURL("http://example.com/not-registered"),
	})
	assert.ErrorContains(t, err, "no registered stubs matching")

	resp, err := transport.RoundTrip(&http.Request{
		Method: http.MethodGet,
		URL:    parseURL("http://example.com/foo"),
	})
	assert.NoError(t, err)
	assert.Equal(t, `{"foo": true}`, httpResponseBody(resp))

	resp, err = transport.RoundTrip(&http.Request{
		Method: http.MethodGet,
		URL:    parseURL("http://example.com/foo"),
	})
	assert.NoError(t, err)
	assert.Equal(t, `{"foo": false}`, httpResponseBody(resp))

	_, err = transport.RoundTrip(&http.Request{
		Method: http.MethodGet,
		URL:    parseURL("http://example.com/foo"),
	})
	assert.ErrorContains(t, err, "wanted 3 of only 2 stubs matching")

	assert.Equal(t, 2, len(transport.Requests))
}

func TestStubbedTransport_VerifyWhenNoStubs(t *testing.T) {
	mt := &mockTest{}
	transport := NewStubbedTransport()

	transport.VerifyStubs(mt)
	assert.Equal(t, true, mt.HelperCalled)
	assert.Equal(t, false, mt.ErrorfCalled)
}

func TestStubbedTransport_VerifyWhenAllStubsMatched(t *testing.T) {
	mt := &mockTest{}
	transport := NewStubbedTransport()

	transport.
		RegisterStub(
			MatchGet("/foo"),
			StringResponse(`{"foo": true}`),
		)

	_, err := transport.RoundTrip(&http.Request{
		Method: http.MethodGet,
		URL:    parseURL("http://example.com/foo"),
	})
	assert.NoError(t, err)

	transport.VerifyStubs(mt)
	assert.Equal(t, true, mt.HelperCalled)
	assert.Equal(t, false, mt.ErrorfCalled)
}

func TestStubbedTransport_VerifyWhenUnmatchedStubs(t *testing.T) {
	mt := &mockTest{}
	transport := NewStubbedTransport()

	transport.
		RegisterStub(
			MatchGet("/foo"),
			StringResponse(`{"foo": true}`),
		)

	transport.VerifyStubs(mt)
	assert.Equal(t, true, mt.HelperCalled)
	assert.Equal(t, true, mt.ErrorfCalled)
	assert.Equal(t, "found 1 unmatched stub(s)", mt.Msg)
}

type mockTest struct {
	Msg          string
	HelperCalled bool
	ErrorfCalled bool
}

func (mt *mockTest) Helper() {
	mt.HelperCalled = true
}
func (mt *mockTest) Errorf(line string, args ...interface{}) {
	mt.ErrorfCalled = true
	mt.Msg = fmt.Sprintf(line, args...)
}

func httpResponseBody(resp *http.Response) string {
	buf := &bytes.Buffer{}
	_, err := buf.ReadFrom(resp.Body)
	if err != nil {
		panic(err)
	}
	return buf.String()
}
