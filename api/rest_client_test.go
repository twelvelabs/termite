package api

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/prashantv/gostub" // spell: disable-line
	"github.com/stretchr/testify/assert"
)

func TestNewRESTClient(t *testing.T) {
	client := NewRESTClient(&ClientOptions{
		BaseURL: "http://example.com/",
	})
	assert.Equal(t, "http://example.com/", client.BaseURL)

	assert.Panics(t, func() {
		_ = NewRESTClient(&ClientOptions{
			BaseURL: "",
		})
	})
	assert.Panics(t, func() {
		_ = NewRESTClient(&ClientOptions{
			BaseURL: "/not/fully/qualified",
		})
	})
}

func TestRESTClient_abs(t *testing.T) {
	client := NewRESTClient(&ClientOptions{
		BaseURL: "http://example.com/",
	})
	assert.Equal(t, "http://example.com/api/v1", client.abs("/api/v1"))
	assert.Equal(t, "https://other.com/foo", client.abs("https://other.com/foo"))
}

func TestRESTClient_DoWithContext(t *testing.T) {
	ctx := context.Background()
	client := greetingClient(MatchAny, StringResponse(`{"msg": "Howdy"}`))
	defer client.VerifyStubs(t)

	g := &greeting{}
	err := client.DoWithContext(ctx, http.MethodGet, "/greet", "", g)
	assert.NoError(t, err)
	assert.Equal(t, "Howdy", g.Msg)
}

func TestRESTClient_DoWithContext_JSONEncodeError(t *testing.T) {
	ctx := context.Background()
	client := NewRESTClient(&ClientOptions{
		BaseURL: "http://example.com/",
	}).WithStubbing()
	defer client.VerifyStubs(t)

	err := client.DoWithContext(ctx, http.MethodPost, "/greet", make(chan int), nil)
	assert.ErrorContains(t, err, "unsupported type")
}

func TestRESTClient_DoWithContext_JSONDecodeError(t *testing.T) {
	ctx := context.Background()
	client := greetingClient(MatchAny, StringResponse(`{"msg":`))
	defer client.VerifyStubs(t)

	g := &greeting{}
	err := client.DoWithContext(ctx, http.MethodGet, "/greet", "", g)
	assert.ErrorContains(t, err, "unexpected end of JSON")
}

func TestRESTClient_DoWithContext_RequestError(t *testing.T) {
	ctx := context.Background()
	client := NewRESTClient(&ClientOptions{
		BaseURL: "http://example.com/",
	}).WithStubbing()
	defer client.VerifyStubs(t)

	err := client.DoWithContext(ctx, "{NOPE}", "/greet", "", nil)
	assert.ErrorContains(t, err, "invalid method")
}

func TestRESTClient_DoWithContext_ResponseError(t *testing.T) {
	ctx := context.Background()
	client := greetingClient(MatchAny, ErrorResponse(errors.New("boom")))
	defer client.VerifyStubs(t)

	err := client.DoWithContext(ctx, http.MethodGet, "/greet", "", nil)
	assert.ErrorContains(t, err, "boom")
}

func TestRESTClient_DoWithContext_SuccessError(t *testing.T) {
	ctx := context.Background()
	client := greetingClient(
		MatchAny,
		WithStatus(500, StringResponse("welp")),
	)
	defer client.VerifyStubs(t)

	err := client.DoWithContext(ctx, http.MethodGet, "/greet", "", nil)
	assert.ErrorContains(t, err, "HTTP 500")
}

func TestRESTClient_DoWithContext_NoContent(t *testing.T) {
	ctx := context.Background()
	client := greetingClient(
		MatchAny,
		WithStatus(http.StatusNoContent, StringResponse("")),
	)
	defer client.VerifyStubs(t)

	err := client.DoWithContext(ctx, http.MethodGet, "/greet", "", nil)
	assert.NoError(t, err)
}

func TestRESTClient_DoWithContext_ReadError(t *testing.T) {
	stubs := gostub.StubFunc(&ioReadAll, nil, errors.New("boom"))
	defer stubs.Reset()

	ctx := context.Background()
	client := greetingClient(
		MatchAny,
		StringResponse(""),
	)
	defer client.VerifyStubs(t)

	g := &greeting{}
	err := client.DoWithContext(ctx, http.MethodGet, "/greet", "", g)
	assert.ErrorContains(t, err, "boom")
}

func TestRESTClient_Delete(t *testing.T) {
	client := greetingClient(MatchDelete("/foo"), StringResponse(""))
	defer client.VerifyStubs(t)

	err := client.Delete("/foo", nil)
	assert.NoError(t, err)
}

func TestRESTClient_Get(t *testing.T) {
	client := greetingClient(MatchGet("/foo"), StringResponse(""))
	defer client.VerifyStubs(t)

	err := client.Get("/foo", nil)
	assert.NoError(t, err)
}

func TestRESTClient_Patch(t *testing.T) {
	client := greetingClient(MatchPatch("/foo"), StringResponse(""))
	defer client.VerifyStubs(t)

	err := client.Patch("/foo", "", nil)
	assert.NoError(t, err)
}

func TestRESTClient_Post(t *testing.T) {
	client := greetingClient(MatchPost("/foo"), StringResponse(""))
	defer client.VerifyStubs(t)

	err := client.Post("/foo", "", nil)
	assert.NoError(t, err)
}

func TestRESTClient_Put(t *testing.T) {
	client := greetingClient(MatchPut("/foo"), StringResponse(""))
	defer client.VerifyStubs(t)

	err := client.Put("/foo", "", nil)
	assert.NoError(t, err)
}

type greeting struct {
	Msg string
}

func greetingClient(m Matcher, r Responder) *RESTClient {
	return NewRESTClient(&ClientOptions{
		BaseURL: "http://example.com/",
	}).WithStubbing().RegisterStub(m, r)
}
