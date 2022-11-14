package run

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	assert.Equal(t, DefaultExecutor, client.Executor)
}

func TestClient_Command(t *testing.T) {
	client := NewClient()
	client.Executor = &ExecutorMock{}

	cmd := client.Command("/bin/echo", "foo", "bar")
	assert.Equal(t, "/bin/echo", cmd.Path)
	assert.Equal(t, []string{"/bin/echo", "foo", "bar"}, cmd.Args)
}

func TestClient_CommandContext(t *testing.T) {
	client := NewClient()
	client.Executor = &ExecutorMock{}
	ctx := context.Background()

	cmd := client.CommandContext(ctx, "/bin/echo", "foo", "bar")
	assert.Equal(t, "/bin/echo", cmd.Path)
	assert.Equal(t, []string{"/bin/echo", "foo", "bar"}, cmd.Args)
}

func TestClient_RegisterStub(t *testing.T) {
	// Stubbing not yet enabled
	client := NewClient()
	assert.Panics(t, func() {
		client.RegisterStub(MatchAny, StringResponse(""))
	})
	assert.Panics(t, func() {
		client.VerifyStubs(t)
	})

	// Stubbing enabled
	client = NewClient().WithStubbing()
	client.RegisterStub(MatchAny, StringResponse(""))
	defer client.VerifyStubs(t)

	err := client.Command("/bin/echo").Run()
	assert.NoError(t, err)
}

func TestClient_StubbingMethods(t *testing.T) {
	client := NewClient()
	assert.Equal(t, false, client.IsStubbed())
	assert.IsType(t, &defaultExecutor{}, client.Executor)

	client = NewClient().WithStubbing()

	assert.Equal(t, true, client.IsStubbed())
	assert.IsType(t, &StubExecutor{}, client.Executor)
}
