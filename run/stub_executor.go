package run

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
)

var (
	_ Executor = &StubExecutor{}
)

// Stub is a stubbed Cmd execution.
type Stub struct {
	Matched   bool
	Matcher   Matcher
	Responder Responder
}

func NewStubExecutor() *StubExecutor {
	return &StubExecutor{
		Commands: []*Cmd{},
		stubs:    []*Stub{},
	}
}

// StubExecutor is an implementation of Executor that executes stubbed commands.
type StubExecutor struct {
	Commands []*Cmd

	mu    sync.Mutex
	stubs []*Stub
}

// RegisterStub registers a new stub for the given matcher/responder pair.
func (e *StubExecutor) RegisterStub(matcher Matcher, responder Responder) *StubExecutor {
	e.mu.Lock()
	e.stubs = append(e.stubs, &Stub{
		Matcher:   matcher,
		Responder: responder,
	})
	e.mu.Unlock()
	return e
}

func (e *StubExecutor) ExitCode(cmd *Cmd) int {
	return cmd.exitCode
}

func (e *StubExecutor) Output(cmd *Cmd) ([]byte, error) {
	if cmd.Stdout != nil {
		return nil, errors.New("run: Stdout already set")
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout

	captureErr := cmd.Stderr == nil
	if captureErr {
		cmd.Stderr = &stderr
	}

	err := e.Run(cmd)
	if err != nil && captureErr {
		var exitErr *ExitError
		if errors.As(err, &exitErr) {
			exitErr.Stderr = stderr.Bytes()
		}
	}
	return stdout.Bytes(), err
}

func (e *StubExecutor) Run(cmd *Cmd) error {
	e.mu.Lock()
	var stub *Stub
	var matches []*Stub

	for _, s := range e.stubs {
		if s.Matcher(cmd) {
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
			return fmt.Errorf("no registered stubs matching: %s", cmd.DebugString())
		} else {
			return fmt.Errorf("wanted %d of only %d stubs matching: %s", n+1, n, cmd.DebugString())
		}
	}

	e.Commands = append(e.Commands, cmd)
	e.mu.Unlock()

	stdout, _, err := stub.Responder(cmd)
	if cmd.Stdout != nil {
		_, we := cmd.Stdout.Write(stdout)
		if we != nil {
			panic(we)
		}
	}

	var exitErr *ExitError
	if errors.As(err, &exitErr) {
		cmd.exitCode = exitErr.Code()
	} else if err == nil {
		cmd.exitCode = 0
	}

	return err
}

// VerifyStubs fails the test if there are unmatched stubs.
func (e *StubExecutor) VerifyStubs(test testable) {
	test.Helper()

	n := 0
	for _, s := range e.stubs {
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
