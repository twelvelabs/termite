package ui

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2" // spell: disable-line

	"github.com/twelvelabs/termite/validate"
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
func (p *SurveyPrompter) Confirm(msg string, value bool, opts ...PromptOpt) (bool, error) {
	params := GetPromptParams(opts...)

	askOpts := []survey.AskOpt{
		survey.WithValidator(surveyValidator(msg, params.ValidationRules)),
	}

	result := value
	err := p.ask(&survey.Confirm{
		Message: msg,
		Help:    params.Help,
		Default: value,
	}, &result, askOpts...)

	return result, err
}

// Input prompts for single string value.
func (p *SurveyPrompter) Input(msg string, value string, opts ...PromptOpt) (string, error) {
	params := GetPromptParams(opts...)

	askOpts := []survey.AskOpt{
		survey.WithValidator(surveyValidator(msg, params.ValidationRules)),
	}

	result := value
	err := p.ask(&survey.Input{
		Message: msg,
		Help:    params.Help,
		Default: value,
	}, &result, askOpts...)

	return result, err
}

// MultiSelect prompts for a slice of string values w/ a fixed set of options.
func (p *SurveyPrompter) MultiSelect(
	msg string, options []string, values []string, opts ...PromptOpt,
) ([]string, error) {
	params := GetPromptParams(opts...)

	askOpts := []survey.AskOpt{
		survey.WithValidator(surveyValidator(msg, params.ValidationRules)),
	}

	result := values
	err := p.ask(&survey.MultiSelect{
		Message: msg,
		Help:    params.Help,
		Options: options,
		Default: values,
	}, &result, askOpts...)

	return result, err
}

// Select prompts for single string value w/ a fixed set of options.
func (p *SurveyPrompter) Select(
	msg string, options []string, value string, opts ...PromptOpt,
) (string, error) {
	params := GetPromptParams(opts...)

	askOpts := []survey.AskOpt{
		survey.WithValidator(surveyValidator(msg, params.ValidationRules)),
	}

	result := value
	err := p.ask(&survey.Select{
		Message: msg,
		Help:    params.Help,
		Options: options,
		Default: value,
	}, &result, askOpts...)

	return result, err
}

func (p *SurveyPrompter) ask(q survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	if !p.ios.IsInteractive() {
		// Do not prompt when non-interactive (the default value will be returned).
		return nil
	}
	opts = append(opts, survey.WithStdio(p.ios.In, p.ios.Out, p.ios.Err))
	// survey.AskOne() doesn't allow passing in a transform func,
	// so we need to call survey.Ask().
	qs := []*survey.Question{
		{
			Prompt:    q,
			Transform: trimSpace,
		},
	}
	err := surveyAsk(qs, response, opts...)
	if err == nil {
		return nil
	}
	return fmt.Errorf("prompt error: %w", err)
}

// Similar to `survey.TransformString(strings.TrimSpace)`, but
// will ignore (and pass through) non-string values.
// This allows us to use it on all prompts.
func trimSpace(val any) any {
	if str, ok := val.(string); ok {
		return strings.TrimSpace(str)
	}
	return val
}

// Returns a survey.Validator that delegates to value.Validate.
func surveyValidator(key string, rules string) survey.Validator {
	return func(val any) error {
		return validate.KeyVal(key, val, rules)
	}
}
