package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserInterface(t *testing.T) {
	ui := NewUserInterface(NewTestIOStreams())
	assert.IsType(t, &UserInterface{}, ui)
}

func TestUserInterface_Out(t *testing.T) {
	ios := NewTestIOStreams()
	ui := NewUserInterface(ios)

	ui.Out("Hello\n")
	ui.Out("Hello %s\n", "World")

	assert.Equal(t, []string{
		"Hello",
		"Hello World",
	}, ios.Out.Lines())
}

func TestUserInterface_Err(t *testing.T) {
	ios := NewTestIOStreams()
	ui := NewUserInterface(ios)

	ui.Err("Hello\n")
	ui.Err("Hello %s\n", "World")

	assert.Equal(t, []string{
		"Hello",
		"Hello World",
	}, ios.Err.Lines())
}

func TestUserInterface_StubbingMethods(t *testing.T) {
	// Stubbing not yet enabled
	ui := NewUserInterface(NewTestIOStreams())
	assert.Equal(t, false, ui.IsStubbed())
	assert.IsType(t, &SurveyPrompter{}, ui.Prompter)
	assert.Panics(t, func() {
		ui.RegisterStub(MatchInput("Name:"), RespondString("foo"))
	})
	assert.Panics(t, func() {
		ui.VerifyStubs(t)
	})

	// Stubbing enabled
	ui = NewUserInterface(NewTestIOStreams()).WithStubbing()
	assert.Equal(t, true, ui.IsStubbed())
	assert.IsType(t, &StubPrompter{}, ui.Prompter)
}

func TestUserInterface_Confirm(t *testing.T) {
	ui := NewUserInterface(NewTestIOStreams()).WithStubbing()
	ui.RegisterStub(MatchConfirm("Perform?"), RespondBool(true))
	defer ui.VerifyStubs(t)

	response, err := ui.Confirm("Perform?", false, "")
	assert.NoError(t, err)
	assert.Equal(t, true, response)
}

func TestUserInterface_Input(t *testing.T) {
	ui := NewUserInterface(NewTestIOStreams()).WithStubbing()
	ui.RegisterStub(MatchInput("Name:"), RespondString("foo"))
	defer ui.VerifyStubs(t)

	response, err := ui.Input("Name:", "", "")
	assert.NoError(t, err)
	assert.Equal(t, "foo", response)
}

func TestUserInterface_MultiSelect(t *testing.T) {
	ui := NewUserInterface(NewTestIOStreams()).WithStubbing()
	ui.RegisterStub(MatchMultiSelect("Color:"), RespondStringSlice([]string{"red"}))
	defer ui.VerifyStubs(t)

	response, err := ui.MultiSelect("Color:", nil, nil, "")
	assert.NoError(t, err)
	assert.Equal(t, []string{"red"}, response)
}

func TestUserInterface_Select(t *testing.T) {
	ui := NewUserInterface(NewTestIOStreams()).WithStubbing()
	ui.RegisterStub(MatchSelect("Country:"), RespondString("US"))
	defer ui.VerifyStubs(t)

	response, err := ui.Select("Country:", nil, "", "")
	assert.NoError(t, err)
	assert.Equal(t, "US", response)
}
