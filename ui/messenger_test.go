package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/twelvelabs/termite/ioutil"
)

func TestMessenger_OutMethods(t *testing.T) {
	ios := ioutil.Test()

	msg := NewMessenger(ios)
	msg.Info("hello %s\n", "world")
	msg.Success("hello %s\n", "world")
	msg.Warning("hello %s\n", "world")
	msg.Failure("hello %s\n", "world")

	assert.Equal(t, "", ios.Err.String())
	assert.Equal(t, []string{
		"• hello world",
		"✓ hello world",
		"! hello world",
		"✖ hello world",
	}, ios.Out.Lines())
}

func TestMessenger_ErrMethods(t *testing.T) {
	ios := ioutil.Test()

	msg := NewMessenger(ios)
	msg.InfoErr("hello %s\n", "world")
	msg.SuccessErr("hello %s\n", "world")
	msg.WarningErr("hello %s\n", "world")
	msg.FailureErr("hello %s\n", "world")

	assert.Equal(t, "", ios.Out.String())
	assert.Equal(t, []string{
		"• hello world",
		"✓ hello world",
		"! hello world",
		"✖ hello world",
	}, ios.Err.Lines())
}

func TestMessenger_TagOutMethods(t *testing.T) {
	ios := ioutil.Test()

	msg := NewMessenger(ios)
	msg.InfoTag("something", "hello %s\n", "world")
	msg.SuccessTag("something", "hello %s\n", "world")
	msg.WarningTag("something", "hello %s\n", "world")
	msg.FailureTag("something", "hello %s\n", "world")

	assert.Equal(t, "", ios.Err.String())
	assert.Equal(t, []string{
		"• [something] hello world",
		"✓ [something] hello world",
		"! [something] hello world",
		"✖ [something] hello world",
	}, ios.Out.Lines())
}

func TestMessenger_TagErrMethods(t *testing.T) {
	ios := ioutil.Test()

	msg := NewMessenger(ios)
	msg.InfoTagErr("something", "hello %s\n", "world")
	msg.SuccessTagErr("something", "hello %s\n", "world")
	msg.WarningTagErr("something", "hello %s\n", "world")
	msg.FailureTagErr("something", "hello %s\n", "world")

	assert.Equal(t, "", ios.Out.String())
	assert.Equal(t, []string{
		"• [something] hello world",
		"✓ [something] hello world",
		"! [something] hello world",
		"✖ [something] hello world",
	}, ios.Err.Lines())
}
