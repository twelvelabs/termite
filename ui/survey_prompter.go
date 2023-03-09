package ui

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2" // spell: disable-line
)

var (
	// for test stubbing
	surveyAsk = survey.Ask

	// interface assertions
	_ Prompter           = &SurveyPrompter{}
	_ survey.Transformer = trimSpace
)

func NewSurveyPrompter(ios *IOStreams) *SurveyPrompter {
	return &SurveyPrompter{
		ios: ios,
	}
}

// SurveyPrompter is a light wrapper around the [survey](https://github.com/go-survey/survey) library.
type SurveyPrompter struct {
	ios *IOStreams
}

// Confirm prompts for a boolean yes/no value.
func (p *SurveyPrompter) Confirm(prompt string, defaultValue bool, help string) (bool, error) {
	result := defaultValue
	err := p.ask(&survey.Confirm{
		Message: prompt,
		Help:    help,
		Default: defaultValue,
	}, &result)
	return result, err
}

// Input prompts for single string value.
func (p *SurveyPrompter) Input(prompt string, defaultValue string, help string) (string, error) {
	result := defaultValue
	err := p.ask(&survey.Input{
		Message: prompt,
		Help:    help,
		Default: defaultValue,
	}, &result)
	return result, err
}

// MultiSelect prompts for a slice of string values w/ a fixed set of options.
func (p *SurveyPrompter) MultiSelect(
	prompt string, options []string, defaultValues []string, help string,
) ([]string, error) {
	result := defaultValues
	err := p.ask(&survey.MultiSelect{
		Message: prompt,
		Help:    help,
		Options: options,
		Default: defaultValues,
	}, &result)
	return result, err
}

// Select prompts for single string value w/ a fixed set of options.
func (p *SurveyPrompter) Select(
	prompt string, options []string, defaultValue string, help string,
) (string, error) {
	result := defaultValue
	err := p.ask(&survey.Select{
		Message: prompt,
		Help:    help,
		Options: options,
		Default: defaultValue,
	}, &result)
	return result, err
}

func (p *SurveyPrompter) ask(q survey.Prompt, response interface{}) error {
	if !p.ios.IsInteractive() {
		return nil
	}
	// survey.AskOne() doesn't allow passing in a transform func,
	// so we need to call survey.Ask().
	qs := []*survey.Question{
		{
			Prompt:    q,
			Transform: trimSpace,
		},
	}
	err := surveyAsk(qs, response, survey.WithStdio(p.ios.In, p.ios.Out, p.ios.Err))
	if err == nil {
		return nil
	}
	return fmt.Errorf("prompt error: %w", err)
}

func trimSpace(val any) any {
	if str, ok := val.(string); ok {
		return strings.TrimSpace(str)
	}
	return val
}
