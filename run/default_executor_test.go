package run

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultExecutor_ExitCode(t *testing.T) {
	client := NewClient()
	assert.Equal(t, DefaultExecutor, client.Executor)

	cmd := client.Command("echo", "foo", "bar")
	assert.Equal(t, -1, client.Executor.ExitCode(cmd))

	err := cmd.Run()
	assert.NoError(t, err)
	assert.Equal(t, 0, client.Executor.ExitCode(cmd))
}

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
