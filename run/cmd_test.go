package run

import (
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
