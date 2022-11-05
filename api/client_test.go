package api

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_Constructors(t *testing.T) {
	client := NewClient()
	assert.Equal(t, time.Duration(0), client.Timeout)

	transport := &http.Transport{}
	client = NewClientWith(&ClientOptions{
		AuthToken: "TOKEN",
		BaseURL:   "http://example.com",
		Timeout:   10,
		Transport: transport,
	})
	assert.Equal(t, time.Duration(10), client.Timeout)
}

func TestClient_RegisterStub(t *testing.T) {
	// Stubbing not yet enabled
	client := NewClient()
	assert.Panics(t, func() {
		client.RegisterStub(MatchAny, StringResponse(""))
	})

	// Stubbing enabled
	client = NewClient().WithStubbing()
	defer client.VerifyStubs(t)

	client.RegisterStub(MatchAny, StringResponse("howdy"))
	resp, err := client.Transport.RoundTrip(&http.Request{})
	assert.NoError(t, err)
	assert.Equal(t, "howdy", httpResponseBody(resp))
}

func TestClient_StubbingMethods(t *testing.T) {
	client := NewClient()
	assert.Equal(t, false, client.IsStubbed())
	assert.IsType(t, &http.Transport{}, client.Transport)

	client = NewClient().WithStubbing()
	transport := client.Transport

	assert.Equal(t, true, client.IsStubbed())
	assert.IsType(t, &StubbedTransport{}, transport)

	c := client.WithStubbing()
	assert.Equal(t, client, c, "should return the receiver and be chainable")
	assert.Equal(t, transport, c.Transport, "should not overwrite existing transport")
}
