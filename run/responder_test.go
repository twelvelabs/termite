package run

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	cmd := NewClient().Command("/bin/echo")
	responder := ErrorResponse(errors.New("boom"))
	buf, err := responder(cmd)
	assert.Nil(t, buf)
	assert.ErrorContains(t, err, "boom")
}

func TestRegexpResponse(t *testing.T) {
	cmd := NewClient().Command("/bin/echo", "something")
	responder := RegexpResponse(`echo\s+(\w+)$`, 1)
	buf, err := responder(cmd)
	assert.NoError(t, err)
	assert.Equal(t, "something", string(buf))
}

func TestRegexpResponse_WhenInvalidIndex(t *testing.T) {
	cmd := NewClient().Command("/bin/echo", "something")
	responder := RegexpResponse(`echo\s+(\w+)$`, 2)
	assert.PanicsWithValue(
		t,
		`invalid match: cmd='/bin/echo something', pattern='echo\s+(\w+)$', index='2'`,
		func() {
			_, _ = responder(cmd)
		},
	)
}

func TestStringResponse(t *testing.T) {
	cmd := NewClient().Command("/bin/echo")
	responder := StringResponse("foo")
	buf, err := responder(cmd)
	assert.NoError(t, err)
	assert.Equal(t, "foo", string(buf))
}
