package ioutil

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystem(t *testing.T) {
	ios := System()
	assert.NotNil(t, ios.In)
	assert.NotNil(t, ios.Out)
	assert.NotNil(t, ios.Err)
}

func TestIOStreamImplementations(t *testing.T) {
	s := &systemIOStream{File: os.Stdin}
	_ = s.String()
	_ = s.Lines()

	m := &mockIOStream{Buffer: bytes.NewBufferString(""), fd: 1}
	assert.Equal(t, 1, int(m.Fd()))
	assert.Equal(t, "", m.String())
	assert.Equal(t, []string{}, m.Lines())

	m = &mockIOStream{Buffer: bytes.NewBufferString("\n"), fd: 1}
	assert.Equal(t, 1, int(m.Fd()))
	assert.Equal(t, "\n", m.String())
	assert.Equal(t, []string{""}, m.Lines())

	m = &mockIOStream{Buffer: bytes.NewBufferString("foo\nbar\nbaz\n"), fd: 1}
	assert.Equal(t, 1, int(m.Fd()))
	assert.Equal(t, "foo\nbar\nbaz\n", m.String())
	assert.Equal(t, []string{"foo", "bar", "baz"}, m.Lines())
}

func TestIOStreams_Formatter(t *testing.T) {
	ios := Test()
	formatter := ios.Formatter()
	f := formatter.ColorFromString("green")
	assert.Equal(t, "foo", f("foo"))
}

func TestIOStreams_ProgressIndicator(t *testing.T) {
	ios := Test()

	assert.Equal(t, false, ios.ProgressIndicatorEnabled())
	ios.StartProgressIndicator()
	ios.StopProgressIndicator()

	assert.Equal(t, "", ios.Err.String())

	ios.SetProgressIndicatorEnabled(true)
	assert.Equal(t, true, ios.ProgressIndicatorEnabled())
	ios.StartProgressIndicator()
	ios.StartProgressIndicatorWithLabel("running")
	ios.StartProgressIndicatorWithLabel("")
	ios.StopProgressIndicator()
	ios.StartProgressIndicatorWithLabel("updating")
	ios.StopProgressIndicator()

	// The spinner library does isTTY checks internally, so we can't get output.
	// Doing the above solely for the coverage stats :money:.
	assert.Equal(t, "", ios.Err.String())
}

func TestIOStreams_TTYMethods(t *testing.T) {
	ios := System()
	// IsTerminal returns false when running tests
	assert.Equal(t, false, ios.IsStdinTTY())
	assert.Equal(t, false, ios.IsStdoutTTY())
	assert.Equal(t, false, ios.IsStderrTTY())
	assert.Equal(t, false, ios.IsInteractive())
	// But we can override the values
	ios.SetStdinTTY(true)
	ios.SetStdoutTTY(true)
	ios.SetStderrTTY(true)
	assert.Equal(t, true, ios.IsStdinTTY())
	assert.Equal(t, true, ios.IsStdoutTTY())
	assert.Equal(t, true, ios.IsStderrTTY())
	assert.Equal(t, true, ios.IsInteractive())
	// And we can disable interactivity, even when supposedly running in a TTY.
	ios.SetInteractive(false)
	assert.Equal(t, false, ios.IsInteractive())
}
