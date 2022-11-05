package api

import (
	"errors"
	"net/http"
	"testing"

	"github.com/prashantv/gostub" // spell: disable-line
	"github.com/stretchr/testify/assert"
)

func TestFileResponse(t *testing.T) {
	responder := FileResponse("testdata/baz.json")
	resp, err := responder(&http.Request{})
	assert.NoError(t, err)
	assert.Equal(t, `{"baz": true}`, httpResponseBody(resp))
}

func TestFileResponseWithError(t *testing.T) {
	responder := FileResponse("testdata/missing.json")
	resp, err := responder(&http.Request{})
	assert.ErrorContains(t, err, "no such file or directory")
	assert.Nil(t, resp)
}

func TestJSONResponse(t *testing.T) {
	type serializable struct {
		Name   string `json:"displayName"`
		Active bool   `json:"isActive"`
	}

	responder := JSONResponse(&serializable{
		Name:   "untitled",
		Active: true,
	})
	resp, err := responder(&http.Request{})
	assert.NoError(t, err)
	assert.Equal(t, `{"displayName":"untitled","isActive":true}`, httpResponseBody(resp))
}

func TestJSONResponseWithError(t *testing.T) {
	stubs := gostub.StubFunc(&jsonMarshal, nil, errors.New("boom"))
	defer stubs.Reset()

	responder := JSONResponse(nil)
	resp, err := responder(&http.Request{})
	assert.ErrorContains(t, err, "boom")
	assert.Nil(t, resp)
}

func TestStringResponse(t *testing.T) {
	responder := StringResponse(`{"foo": true}`)
	resp, err := responder(&http.Request{})
	assert.NoError(t, err)
	assert.Equal(t, `{"foo": true}`, httpResponseBody(resp))
}

func TestWithHeader(t *testing.T) {
	responder := WithHeader("X-Foo", "bar", StringResponse(""))
	resp, err := responder(&http.Request{})
	assert.NoError(t, err)
	assert.Equal(t, "bar", resp.Header.Get("X-Foo"))
}

func TestWithHeaderHandlesNilHeader(t *testing.T) {
	responder := WithHeader("X-Foo", "bar", func(req *http.Request) (*http.Response, error) {
		return &http.Response{}, nil
	})
	resp, err := responder(&http.Request{})
	assert.NoError(t, err)
	assert.Equal(t, "bar", resp.Header.Get("X-Foo"))
}

func TestWithStatus(t *testing.T) {
	responder := WithStatus(404, StringResponse(""))
	resp, err := responder(&http.Request{})
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}
