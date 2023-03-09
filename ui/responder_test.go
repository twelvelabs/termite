package ui

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRespondBool(t *testing.T) {
	responder := RespondBool(true)
	response, err := responder(Prompt{})
	assert.NoError(t, err)
	assert.Equal(t, true, response)
}

func TestRespondString(t *testing.T) {
	responder := RespondString("foo")
	response, err := responder(Prompt{})
	assert.NoError(t, err)
	assert.Equal(t, "foo", response)
}

func TestRespondStringSlice(t *testing.T) {
	responder := RespondStringSlice([]string{"foo", "bar", "baz"})
	response, err := responder(Prompt{})
	assert.NoError(t, err)
	assert.Equal(t, []string{"foo", "bar", "baz"}, response)
}

func TestRespondError(t *testing.T) {
	responder := RespondError(errors.New("boom"))
	response, err := responder(Prompt{})
	assert.ErrorContains(t, err, "boom")
	assert.Nil(t, response)
}
