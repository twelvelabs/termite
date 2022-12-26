package run

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	cmd := NewClient().Command("/bin/echo")
	responder := ErrorResponse(errors.New("boom"))
	stdout, stderr, err := responder(cmd)
	assert.Nil(t, stdout)
	assert.Nil(t, stderr)
	assert.ErrorContains(t, err, "boom")
}

func TestRegexpResponse(t *testing.T) {
	cmd := NewClient().Command("/bin/echo", "something")
	responder := RegexpResponse(`echo\s+(\w+)$`, 1)
	stdout, stderr, err := responder(cmd)
	assert.Equal(t, "something", string(stdout))
	assert.Nil(t, stderr)
	assert.NoError(t, err)
}

func TestRegexpResponse_WhenInvalidIndex(t *testing.T) {
	cmd := NewClient().Command("/bin/echo", "something")
	responder := RegexpResponse(`echo\s+(\w+)$`, 2)
	assert.PanicsWithValue(
		t,
		`invalid match: cmd='/bin/echo something', pattern='echo\s+(\w+)$', index='2'`,
		func() {
			_, _, _ = responder(cmd)
		},
	)
}

func TestStringResponse(t *testing.T) {
	cmd := NewClient().Command("/bin/echo")
	responder := StringResponse("foo")
	stdout, stderr, err := responder(cmd)
	assert.Equal(t, "foo", string(stdout))
	assert.Nil(t, stderr)
	assert.NoError(t, err)
}
