package ui

import (
	"errors"
	"testing"

	"github.com/prashantv/gostub" // spell: disable-line
	"github.com/stretchr/testify/assert"

	"github.com/twelvelabs/termite/ioutil"
)

func TestNewSurveyPrompter(t *testing.T) {
	ios := ioutil.Test()

	sp := NewSurveyPrompter(ios.In, ios.Out, ios.Err, ios)
	assert.IsType(t, &SurveyPrompter{}, sp)
}

func TestSurveyPrompter_Confirm(t *testing.T) {
	stubs := gostub.StubFunc(&surveyAsk, errors.New("boom"))
	defer stubs.Reset()

	ios := ioutil.Test()
	sp := NewSurveyPrompter(ios.In, ios.Out, ios.Err, ios)

	// Non-interactive sessions should not prompt - just return the default value.
	ios.SetInteractive(false)
	value, err := sp.Confirm("Explode?", true, "")
	assert.NoError(t, err)
	assert.Equal(t, true, value)

	// Interactive sessions should prompt.
	ios.SetInteractive(true)
	_, err = sp.Confirm("Explode?", true, "")
	assert.ErrorContains(t, err, "prompt error: boom")

	// One stubbed out happy path for coverage stats.
	// We're assuming that Survey works correctly and not testing actual interaction.
	stubs.StubFunc(&surveyAsk, nil)
	_, err = sp.Confirm("Explode?", true, "")
	assert.NoError(t, err)
	assert.Equal(t, true, value)
}

func TestSurveyPrompter_Input(t *testing.T) {
	stubs := gostub.StubFunc(&surveyAsk, errors.New("boom"))
	defer stubs.Reset()

	ios := ioutil.Test()
	sp := NewSurveyPrompter(ios.In, ios.Out, ios.Err, ios)

	// Non-interactive sessions should not prompt - just return the default value.
	ios.SetInteractive(false)
	value, err := sp.Input("Access Code", "12345", "")
	assert.NoError(t, err)
	assert.Equal(t, "12345", value)

	// Interactive sessions should prompt.
	ios.SetInteractive(true)
	_, err = sp.Input("Access Code", "12345", "")
	assert.ErrorContains(t, err, "prompt error: boom")
}

func TestSurveyPrompter_MultiSelect(t *testing.T) {
	stubs := gostub.StubFunc(&surveyAsk, errors.New("boom"))
	defer stubs.Reset()

	ios := ioutil.Test()
	sp := NewSurveyPrompter(ios.In, ios.Out, ios.Err, ios)

	// Non-interactive sessions should not prompt - just return the default value.
	ios.SetInteractive(false)
	value, err := sp.MultiSelect("Prefixes", []string{"a", "b", "c"}, []string{"a"}, "")
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, value)

	// Interactive sessions should prompt.
	ios.SetInteractive(true)
	_, err = sp.MultiSelect("Prefixes", []string{"a", "b", "c"}, []string{"a"}, "")
	assert.ErrorContains(t, err, "prompt error: boom")
}

func TestSurveyPrompter_Select(t *testing.T) {
	stubs := gostub.StubFunc(&surveyAsk, errors.New("boom"))
	defer stubs.Reset()

	ios := ioutil.Test()
	sp := NewSurveyPrompter(ios.In, ios.Out, ios.Err, ios)

	// Non-interactive sessions should not prompt - just return the default value.
	ios.SetInteractive(false)
	value, err := sp.Select("Prefixes", []string{"a", "b", "c"}, "a", "")
	assert.NoError(t, err)
	assert.Equal(t, "a", value)

	// Interactive sessions should prompt.
	ios.SetInteractive(true)
	_, err = sp.Select("Prefixes", []string{"a", "b", "c"}, "a", "")
	assert.ErrorContains(t, err, "prompt error: boom")
}

func TestTrimSpace(t *testing.T) {
	// Should trim strings
	assert.Equal(t, "123", trimSpace("   123   "))
	// Should pass anything else through
	assert.Equal(t, 123, trimSpace(123))
	assert.Equal(t, []string{"a", "b"}, trimSpace([]string{"a", "b"}))
}
