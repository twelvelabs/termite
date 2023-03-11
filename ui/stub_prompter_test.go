package ui

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStubPrompter_Confirm(t *testing.T) {
	ios := NewTestIOStreams()
	prompter := NewStubPrompter(ios)

	_, err := prompter.Confirm("Proceed?", false)
	assert.ErrorContains(t, err, "no registered stubs matching")

	prompter.RegisterStub(
		MatchConfirm("Proceed?"),
		RespondBool(true),
	)
	prompter.RegisterStub(
		MatchConfirm("Proceed?"),
		RespondDefault(),
	)
	prompter.RegisterStub(
		MatchConfirm("Proceed?"),
		RespondError(errors.New("boom")),
	)

	response, err := prompter.Confirm("Proceed?", false)
	assert.NoError(t, err)
	assert.Equal(t, true, response)

	response, err = prompter.Confirm("Proceed?", true)
	assert.NoError(t, err)
	assert.Equal(t, true, response)

	response, err = prompter.Confirm("Proceed?", false)
	assert.ErrorContains(t, err, "boom")
	assert.Equal(t, false, response)

	_, err = prompter.Confirm("Proceed?", false)
	assert.ErrorContains(t, err, "wanted 4 of only 3 stubs matching")

	assert.Equal(t, []string{
		"? Proceed? Yes",
		"? Proceed? Yes",
	}, ios.Out.Lines())
}

func TestStubPrompter_Input(t *testing.T) {
	ios := NewTestIOStreams()
	prompter := NewStubPrompter(ios)

	_, err := prompter.Input("Name:", "")
	assert.ErrorContains(t, err, "no registered stubs matching")

	prompter.RegisterStub(
		MatchInput("Name:"),
		RespondString("foo"),
	)
	prompter.RegisterStub(
		MatchInput("Name:"),
		RespondDefault(),
	)
	prompter.RegisterStub(
		MatchInput("Name:"),
		RespondError(errors.New("boom")),
	)

	response, err := prompter.Input("Name:", "")
	assert.NoError(t, err)
	assert.Equal(t, "foo", response)

	response, err = prompter.Input("Name:", "something")
	assert.NoError(t, err)
	assert.Equal(t, "something", response)

	response, err = prompter.Input("Name:", "")
	assert.ErrorContains(t, err, "boom")
	assert.Equal(t, "", response)

	_, err = prompter.Input("Name:", "")
	assert.ErrorContains(t, err, "wanted 4 of only 3 stubs matching")

	assert.Equal(t, []string{
		"? Name: foo",
		"? Name: something",
	}, ios.Out.Lines())
}

func TestStubPrompter_MultiSelect(t *testing.T) {
	ios := NewTestIOStreams()
	prompter := NewStubPrompter(ios)

	_, err := prompter.MultiSelect("Colors:", nil, nil)
	assert.ErrorContains(t, err, "no registered stubs matching")

	prompter.RegisterStub(
		MatchMultiSelect("Colors:"),
		RespondStringSlice([]string{"red", "yellow", "blue"}),
	)
	prompter.RegisterStub(
		MatchMultiSelect("Colors:"),
		RespondDefault(),
	)
	prompter.RegisterStub(
		MatchMultiSelect("Colors:"),
		RespondError(errors.New("boom")),
	)

	response, err := prompter.MultiSelect("Colors:", nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, []string{"red", "yellow", "blue"}, response)

	response, err = prompter.MultiSelect("Colors:", nil, []string{"purple", "green"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"purple", "green"}, response)

	response, err = prompter.MultiSelect("Colors:", nil, nil)
	assert.ErrorContains(t, err, "boom")
	assert.Nil(t, response)

	_, err = prompter.MultiSelect("Colors:", nil, nil)
	assert.ErrorContains(t, err, "wanted 4 of only 3 stubs matching")

	assert.Equal(t, []string{
		"? Colors: red, yellow, blue",
		"? Colors: purple, green",
	}, ios.Out.Lines())
}

func TestStubPrompter_Select(t *testing.T) {
	ios := NewTestIOStreams()
	prompter := NewStubPrompter(ios)

	_, err := prompter.Select("Country:", nil, "")
	assert.ErrorContains(t, err, "no registered stubs matching")

	prompter.RegisterStub(
		MatchSelect("Country:"),
		RespondString("US"),
	)
	prompter.RegisterStub(
		MatchSelect("Country:"),
		RespondDefault(),
	)
	prompter.RegisterStub(
		MatchSelect("Country:"),
		RespondError(errors.New("boom")),
	)

	response, err := prompter.Select("Country:", nil, "")
	assert.NoError(t, err)
	assert.Equal(t, "US", response)

	response, err = prompter.Select("Country:", nil, "CA")
	assert.NoError(t, err)
	assert.Equal(t, "CA", response)

	response, err = prompter.Select("Country:", nil, "")
	assert.ErrorContains(t, err, "boom")
	assert.Equal(t, "", response)

	_, err = prompter.Select("Country:", nil, "")
	assert.ErrorContains(t, err, "wanted 4 of only 3 stubs matching")

	assert.Equal(t, []string{
		"? Country: US",
		"? Country: CA",
	}, ios.Out.Lines())
}

func TestStubPrompter_VerifyWhenNoStubs(t *testing.T) {
	mt := &mockTest{}
	prompter := NewStubPrompter(NewTestIOStreams())

	prompter.VerifyStubs(mt)
	assert.Equal(t, true, mt.HelperCalled)
	assert.Equal(t, false, mt.ErrorfCalled)
}

func TestStubPrompter_VerifyWhenAllStubsMatched(t *testing.T) {
	mt := &mockTest{}
	prompter := NewStubPrompter(NewTestIOStreams())
	prompter.RegisterStub(
		MatchConfirm("Proceed?"),
		RespondBool(true),
	)

	_, err := prompter.Confirm("Proceed?", false)
	assert.NoError(t, err)

	prompter.VerifyStubs(mt)
	assert.Equal(t, true, mt.HelperCalled)
	assert.Equal(t, false, mt.ErrorfCalled)
}

func TestStubPrompter_VerifyWhenUnmatchedStubs(t *testing.T) {
	mt := &mockTest{}
	prompter := NewStubPrompter(NewTestIOStreams())
	prompter.RegisterStub(
		MatchConfirm("Proceed?"),
		RespondBool(true),
	)

	prompter.VerifyStubs(mt)
	assert.Equal(t, true, mt.HelperCalled)
	assert.Equal(t, true, mt.ErrorfCalled)
	assert.Equal(t, "found 1 unmatched stub(s)", mt.Msg)
}

type mockTest struct {
	Msg          string
	HelperCalled bool
	ErrorfCalled bool
}

func (mt *mockTest) Helper() {
	mt.HelperCalled = true
}
func (mt *mockTest) Errorf(line string, args ...interface{}) {
	mt.ErrorfCalled = true
	mt.Msg = fmt.Sprintf(line, args...)
}
