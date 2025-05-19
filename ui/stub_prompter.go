package ui

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/cast"
)

// Stub is a stubbed Prompter invocation.
type Stub struct {
	Matched   bool
	Matcher   Matcher
	Responder Responder
}

func NewStubPrompter(ios *IOStreams) *StubPrompter {
	return &StubPrompter{
		ios:   ios,
		stubs: []*Stub{},
	}
}

// StubPrompter is an implementation of Prompter that invokes stubbed prompts.
type StubPrompter struct {
	ios   *IOStreams
	mu    sync.Mutex
	stubs []*Stub
}

var (
	_ Prompter = &StubPrompter{}
)

// Confirm prompts for a boolean yes/no value.
func (sp *StubPrompter) Confirm(msg string, value bool, opts ...PromptOpt) (bool, error) {
	prompt := Prompt{
		Type:    PromptTypeConfirm,
		Message: msg,
		Value:   value,
	}
	stub, err := sp.match(prompt)
	if err != nil {
		return false, err
	}
	response, err := stub.Responder(prompt)
	if err != nil {
		return false, err
	}
	casted := cast.ToBool(response)
	formatted := "No"
	if casted {
		formatted = "Yes"
	}
	_, _ = fmt.Fprintf(sp.ios.Out, "? %s %s\n", msg, formatted)
	return casted, nil
}

// Input prompts for single string value.
func (sp *StubPrompter) Input(msg string, value string, opts ...PromptOpt) (string, error) {
	prompt := Prompt{
		Type:    PromptTypeInput,
		Message: msg,
		Value:   value,
	}
	stub, err := sp.match(prompt)
	if err != nil {
		return "", err
	}
	response, err := stub.Responder(prompt)
	if err != nil {
		return "", err
	}
	casted := cast.ToString(response)
	_, _ = fmt.Fprintf(sp.ios.Out, "? %s %s\n", msg, casted)
	return casted, nil
}

// MultiSelect prompts for a slice of string values w/ a fixed set of options.
func (sp *StubPrompter) MultiSelect(msg string, options []string, values []string, opts ...PromptOpt) ([]string, error) {
	prompt := Prompt{
		Type:    PromptTypeMultiSelect,
		Message: msg,
		Value:   values,
	}
	stub, err := sp.match(prompt)
	if err != nil {
		return nil, err
	}
	response, err := stub.Responder(prompt)
	if err != nil {
		return nil, err
	}
	casted := cast.ToStringSlice(response)
	_, _ = fmt.Fprintf(sp.ios.Out, "? %s %s\n", msg, strings.Join(casted, ", "))
	return casted, nil
}

// Select prompts for single string value w/ a fixed set of options.
func (sp *StubPrompter) Select(msg string, options []string, value string, opts ...PromptOpt) (string, error) {
	prompt := Prompt{
		Type:    PromptTypeSelect,
		Message: msg,
		Value:   value,
	}
	stub, err := sp.match(prompt)
	if err != nil {
		return "", err
	}
	response, err := stub.Responder(prompt)
	if err != nil {
		return "", err
	}
	casted := cast.ToString(response)
	_, _ = fmt.Fprintf(sp.ios.Out, "? %s %s\n", msg, casted)
	return response.(string), nil
}

func (e *StubPrompter) match(prompt Prompt) (*Stub, error) {
	e.mu.Lock()
	var stub *Stub
	var matches []*Stub

	for _, s := range e.stubs {
		if s.Matcher(prompt) {
			matches = append(matches, s)
			if s.Matched {
				continue
			}
			if stub == nil {
				s.Matched = true
				stub = s
			}
		}
	}

	if stub == nil {
		e.mu.Unlock()
		n := len(matches)
		if n == 0 {
			return nil, fmt.Errorf("no registered stubs matching: %s", prompt)
		} else {
			return nil, fmt.Errorf("wanted %d of only %d stubs matching: %s", n+1, n, prompt)
		}
	}

	e.mu.Unlock()
	return stub, nil
}

// RegisterStub registers a new stub for the given matcher/responder pair.
func (sp *StubPrompter) RegisterStub(matcher Matcher, responder Responder) *StubPrompter {
	sp.mu.Lock()
	sp.stubs = append(sp.stubs, &Stub{
		Matcher:   matcher,
		Responder: responder,
	})
	sp.mu.Unlock()
	return sp
}

// VerifyStubs fails the test if there are unmatched stubs.
func (sp *StubPrompter) VerifyStubs(test testable) {
	test.Helper()

	n := 0
	for _, s := range sp.stubs {
		if !s.Matched {
			n++
		}
	}
	if n > 0 {
		test.Errorf("found %d unmatched stub(s)", n)
	}
}

type testable interface {
	Errorf(string, ...interface{})
	Helper()
}
