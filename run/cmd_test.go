package run

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmd_Output(t *testing.T) {
	client := NewClient()
	client.Executor = &ExecutorMock{
		OutputFunc: func(cmd *Cmd) ([]byte, error) {
			return []byte("foo bar\n"), nil
		},
	}

	cmd := client.Command("/bin/echo", "foo", "bar")
	buf, err := cmd.Output()
	assert.NoError(t, err)
	assert.Equal(t, []byte("foo bar\n"), buf)
}

func TestCmd_Run(t *testing.T) {
	client := NewClient()
	client.Executor = &ExecutorMock{
		RunFunc: func(cmd *Cmd) error {
			return nil
		},
	}

	cmd := client.Command("/bin/echo", "foo", "bar")
	err := cmd.Run()
	assert.NoError(t, err)
}

func TestCmd_DebugString(t *testing.T) {
	client := NewClient()

	cmd := client.Command("/bin/cat")
	assert.Equal(t, `/bin/cat`, cmd.DebugString())

	cmd.Stdin = bytes.NewBufferString("hello")
	assert.Equal(t, `/bin/cat [Stdin: "hello"]`, cmd.DebugString())

	cmd.Stdin = bytes.NewBufferString("Lorem ipsum dolor sit amet")
	assert.Equal(t, `/bin/cat [Stdin: "Lorem ipsum dolor si…"]`, cmd.DebugString())
}

func TestCmd_PeekStdin(t *testing.T) {
	client := NewClient()
	cmd := client.Command("/bin/cat")

	// Stdin not set
	data, err := cmd.PeekStdin()
	assert.Equal(t, "", string(data))
	assert.NoError(t, err)

	cmd.Stdin = bytes.NewBufferString("hello")

	// Std in set
	data, err = cmd.PeekStdin()
	assert.Equal(t, "hello", string(data))
	assert.NoError(t, err)

	// Read offset should still be at the beginning of file
	data, err = io.ReadAll(cmd.Stdin)
	assert.Equal(t, "hello", string(data))
	assert.NoError(t, err)

	// Read offset is now at EOF
	data, err = io.ReadAll(cmd.Stdin)
	assert.Equal(t, "", string(data))
	assert.NoError(t, err)
}
