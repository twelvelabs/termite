package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProgressIndicator(t *testing.T) {
	ios := NewTestIOStreams()
	indicator := NewProgressIndicator(ios)

	assert.Equal(t, false, indicator.IsEnabled())

	indicator.Start()
	indicator.Stop()
	assert.Equal(t, "", ios.Err.String())

	ios.SetStdoutTTY(true)
	ios.SetStderrTTY(true)

	assert.Equal(t, true, indicator.IsEnabled())

	indicator.Start()
	indicator.StartWithLabel("running")
	indicator.StartWithLabel("")
	indicator.Stop()
	indicator.StartWithLabel("updating")
	indicator.Stop()

	// The spinner library does isTTY checks internally, so we can't get output.
	// Doing the above solely for the coverage stats :money:.
	assert.Equal(t, "", ios.Err.String())
}
