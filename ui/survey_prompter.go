package ui

import (
	"fmt"
	"io"
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

type fileReader interface {
	io.Reader
	Fd() uintptr
}

type fileWriter interface {
	io.Writer
	Fd() uintptr
}

type session interface {
	IsInteractive() bool
}

func NewSurveyPrompter(in fileReader, out fileWriter, err fileWriter, s session) *SurveyPrompter {
	return &SurveyPrompter{
		stdin:   in,
		stdout:  out,
		stderr:  err,
		session: s,
	}
}

// SurveyPrompter is a light wrapper around the [survey](https://github.com/go-survey/survey) library.
type SurveyPrompter struct {
	stdin   fileReader
	stdout  fileWriter
	stderr  fileWriter
	session session
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
	if !p.session.IsInteractive() {
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
	err := surveyAsk(qs, response, survey.WithStdio(p.stdin, p.stdout, p.stderr))
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
