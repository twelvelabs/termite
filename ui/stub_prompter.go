package ui

import (
	"fmt"
	"sync"
)

// Stub is a stubbed Prompter invocation.
type Stub struct {
	Matched   bool
	Matcher   Matcher
	Responder Responder
}

func NewStubPrompter() *StubPrompter {
	return &StubPrompter{
		stubs: []*Stub{},
	}
}

// StubPrompter is an implementation of Prompter that invokes stubbed prompts.
type StubPrompter struct {
	mu    sync.Mutex
	stubs []*Stub
}

var (
	_ Prompter = &StubPrompter{}
)

// Confirm prompts for a boolean yes/no value.
func (sp *StubPrompter) Confirm(msg string, value bool, help string) (bool, error) {
	prompt := Prompt{
		Type:    PromptTypeConfirm,
		Message: msg,
	}
	stub, err := sp.match(prompt)
	if err != nil {
		return false, err
	}
	response, err := stub.Responder(prompt)
	if err != nil {
		return false, err
	}
	return response.(bool), nil
}

// Input prompts for single string value.
func (sp *StubPrompter) Input(msg string, value string, help string) (string, error) {
	prompt := Prompt{
		Type:    PromptTypeInput,
		Message: msg,
	}
	stub, err := sp.match(prompt)
	if err != nil {
		return "", err
	}
	response, err := stub.Responder(prompt)
	if err != nil {
		return "", err
	}
	return response.(string), nil
}

// MultiSelect prompts for a slice of string values w/ a fixed set of options.
func (sp *StubPrompter) MultiSelect(msg string, options []string, values []string, help string) ([]string, error) {
	prompt := Prompt{
		Type:    PromptTypeMultiSelect,
		Message: msg,
	}
	stub, err := sp.match(prompt)
	if err != nil {
		return nil, err
	}
	response, err := stub.Responder(prompt)
	if err != nil {
		return nil, err
	}
	return response.([]string), nil
}

// Select prompts for single string value w/ a fixed set of options.
func (sp *StubPrompter) Select(msg string, options []string, value string, help string) (string, error) {
	prompt := Prompt{
		Type:    PromptTypeSelect,
		Message: msg,
	}
	stub, err := sp.match(prompt)
	if err != nil {
		return "", err
	}
	response, err := stub.Responder(prompt)
	if err != nil {
		return "", err
	}
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
