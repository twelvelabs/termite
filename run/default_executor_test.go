package run

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultExecutor_Output(t *testing.T) {
	client := NewClient()
	assert.Equal(t, DefaultExecutor, client.Executor)

	cmd := client.Command("echo", "foo", "bar")
	buf, err := cmd.Output()
	assert.NoError(t, err)
	assert.Equal(t, []byte("foo bar\n"), buf)
}

func TestDefaultExecutor_Run(t *testing.T) {
	client := NewClient()
	assert.Equal(t, DefaultExecutor, client.Executor)

	cmd := client.Command("echo", "foo", "bar")
	err := cmd.Run()
	assert.NoError(t, err)
}
