package ui

import (
	"github.com/twelvelabs/termite/ioutil"
)

// NewUserInterface returns a new UserInterface.
func NewUserInterface(ios *ioutil.IOStreams) *UserInterface {
	return &UserInterface{
		Prompter: NewSurveyPrompter(ios.In, ios.Out, ios.Err, ios),
		ios:      ios,
		format:   ios.Formatter(),
	}
}

// UserInterface is a high level abstraction for rendering terminal UIs.
// It supports stubbing interactive prompts and rendering formatted text to os.Stdout.
type UserInterface struct {
	Prompter Prompter

	ios    *ioutil.IOStreams
	format *ioutil.Formatter
}

/**
* Delegated `Prompter` methods.
**/

// Confirm prompts for a boolean yes/no value.
func (ui *UserInterface) Confirm(msg string, value bool, help string) (bool, error) {
	return ui.Prompter.Confirm(msg, value, help)
}

// Input prompts for single string value.
func (ui *UserInterface) Input(msg string, value string, help string) (string, error) {
	return ui.Prompter.Input(msg, value, help)
}

// MultiSelect prompts for a slice of string values w/ a fixed set of options.
func (ui *UserInterface) MultiSelect(msg string, options []string, values []string, help string) ([]string, error) {
	return ui.Prompter.MultiSelect(msg, options, values, help)
}

// Select prompts for single string value w/ a fixed set of options.
func (ui *UserInterface) Select(msg string, options []string, value string, help string) (string, error) {
	return ui.Prompter.Select(msg, options, value, help)
}

/**
* Stub management methods.
**/

// IsStubbed returns true if the executor is configured for stubbing.
func (ui *UserInterface) IsStubbed() bool {
	_, ok := ui.Prompter.(*StubPrompter)
	return ok
}

// RegisterStub registers a new stub for the given matcher/responder pair.
func (ui *UserInterface) RegisterStub(matcher Matcher, responder Responder) *UserInterface {
	if !ui.IsStubbed() {
		panic("must enable stubbing before registering stubs")
	}
	executor := ui.Prompter.(*StubPrompter)
	executor.RegisterStub(matcher, responder)
	return ui
}

// VerifyStubs fails the test if there are unmatched stubs.
func (ui *UserInterface) VerifyStubs(t testable) {
	if !ui.IsStubbed() {
		panic("must enable stubbing before verifying stubs")
	}
	t.Helper()
	executor := ui.Prompter.(*StubPrompter)
	executor.VerifyStubs(t)
}

// WithStubbing configures stubbing and returns the receiver.
func (ui *UserInterface) WithStubbing() *UserInterface {
	if !ui.IsStubbed() {
		ui.Prompter = NewStubPrompter()
	}
	return ui
}
