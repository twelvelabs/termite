package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

var (
	jsonMarshal = json.Marshal
)

type Responder func(req *http.Request) (*http.Response, error)

// ErrorResponse creates a responder that returns err.
func ErrorResponse(err error) Responder {
	return func(req *http.Request) (*http.Response, error) {
		return nil, err
	}
}

// FileResponse creates a responder that returns the content located at filename.
func FileResponse(filename string) Responder {
	return func(req *http.Request) (*http.Response, error) {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		return httpResponse(200, req, f), nil
	}
}

// JSONResponse creates a responder that serializes and returns data.
func JSONResponse(data interface{}) Responder {
	return func(req *http.Request) (*http.Response, error) {
		b, err := jsonMarshal(data)
		if err != nil {
			return nil, err
		}
		return httpResponse(200, req, bytes.NewBuffer(b)), nil
	}
}

// StringResponse creates a responder that returns body.
func StringResponse(body string) Responder {
	return func(req *http.Request) (*http.Response, error) {
		return httpResponse(200, req, bytes.NewBufferString(body)), nil
	}
}

// WithHeader wraps a responder so that it responds with the provided HTTP header.
func WithHeader(header string, value string, responder Responder) Responder {
	return func(req *http.Request) (*http.Response, error) {
		resp, _ := responder(req)
		if resp.Header == nil {
			resp.Header = make(http.Header)
		}
		resp.Header.Set(header, value)
		return resp, nil
	}
}

// WithRequestHeaders wraps a responder so that it responds with all HTTP headers from the request.
func WithRequestHeaders(responder Responder) Responder {
	return func(req *http.Request) (*http.Response, error) {
		resp, _ := responder(req)
		resp.Header = req.Header.Clone()
		return resp, nil
	}
}

// WithStatus wraps a responder so that it responds with the provided HTTP status code.
func WithStatus(status int, responder Responder) Responder {
	return func(req *http.Request) (*http.Response, error) {
		resp, _ := responder(req)
		resp.StatusCode = status
		return resp, nil
	}
}

func httpResponse(status int, req *http.Request, body io.Reader) *http.Response {
	return &http.Response{
		StatusCode: status,
		Request:    req,
		Body:       io.NopCloser(body),
		Header:     http.Header{},
	}
}
