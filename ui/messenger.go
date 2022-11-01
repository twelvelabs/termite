package ui

import (
	"fmt"

	"github.com/twelvelabs/termite/ioutil"
)

// NewMessenger returns a new IconLogger.
func NewMessenger(ios *ioutil.IOStreams) *Messenger {
	return &Messenger{
		ios:    ios,
		format: ios.Formatter(),
	}
}

// Messenger is a generic logger that prefixes lines with status icons.
type Messenger struct {
	ios    *ioutil.IOStreams
	format *ioutil.Formatter
}

// Info prints line to StdOut prefixed with "•".
func (m *Messenger) Info(line string, args ...any) {
	m.out(m.format.InfoIcon(), line, args...)
}

// InfoErr prints line to StdErr prefixed with "•".
func (m *Messenger) InfoErr(line string, args ...any) {
	m.err(m.format.InfoIcon(), line, args...)
}

// InfoTag prints line to StdOut prefixed with "• [tag] ".
func (m *Messenger) InfoTag(tag string, line string, args ...any) {
	line = "[" + m.format.Info(tag) + "] " + line
	m.out(m.format.InfoIcon(), line, args...)
}

// InfoTagErr prints line to StdErr prefixed with "• [tag] ".
func (m *Messenger) InfoTagErr(tag string, line string, args ...any) {
	line = "[" + m.format.Info(tag) + "] " + line
	m.err(m.format.InfoIcon(), line, args...)
}

// Success prints line to StdOut prefixed with "✓".
func (m *Messenger) Success(line string, args ...any) {
	m.out(m.format.SuccessIcon(), line, args...)
}

// SuccessErr prints line to StdErr prefixed with "✓".
func (m *Messenger) SuccessErr(line string, args ...any) {
	m.err(m.format.SuccessIcon(), line, args...)
}

// SuccessTag prints line to StdOut prefixed with "✓ [tag] ".
func (m *Messenger) SuccessTag(tag string, line string, args ...any) {
	line = "[" + m.format.Success(tag) + "] " + line
	m.out(m.format.SuccessIcon(), line, args...)
}

// SuccessTagErr prints line to StdErr prefixed with "✓ [tag] ".
func (m *Messenger) SuccessTagErr(tag string, line string, args ...any) {
	line = "[" + m.format.Success(tag) + "] " + line
	m.err(m.format.SuccessIcon(), line, args...)
}

// Warning prints line to StdOut prefixed with "!".
func (m *Messenger) Warning(line string, args ...any) {
	m.out(m.format.WarningIcon(), line, args...)
}

// WarningErr prints line to StdErr prefixed with "!".
func (m *Messenger) WarningErr(line string, args ...any) {
	m.err(m.format.WarningIcon(), line, args...)
}

// WarningTag prints line to StdOut prefixed with "! [tag] ".
func (m *Messenger) WarningTag(tag string, line string, args ...any) {
	line = "[" + m.format.Warning(tag) + "] " + line
	m.out(m.format.WarningIcon(), line, args...)
}

// WarningTagErr prints line to StdOut prefixed with "! [tag] ".
func (m *Messenger) WarningTagErr(tag string, line string, args ...any) {
	line = "[" + m.format.Warning(tag) + "] " + line
	m.err(m.format.WarningIcon(), line, args...)
}

// Failure prints line to StdOut prefixed with "✖".
func (m *Messenger) Failure(line string, args ...any) {
	m.out(m.format.FailureIcon(), line, args...)
}

// FailureErr prints line to StdErr prefixed with "✖".
func (m *Messenger) FailureErr(line string, args ...any) {
	m.err(m.format.FailureIcon(), line, args...)
}

// FailureTag prints line to StdOut prefixed with "✖ [tag] ".
func (m *Messenger) FailureTag(tag string, line string, args ...any) {
	line = "[" + m.format.Failure(tag) + "] " + line
	m.out(m.format.FailureIcon(), line, args...)
}

// FailureTagErr prints line to StdErr prefixed with "✖ [tag] ".
func (m *Messenger) FailureTagErr(tag string, line string, args ...any) {
	line = "[" + m.format.Failure(tag) + "] " + line
	m.err(m.format.FailureIcon(), line, args...)
}

// prints to StdOut
func (m *Messenger) out(icon string, line string, args ...any) {
	fmt.Fprintf(m.ios.Out, (icon + " " + line), args...)
}

// prints to StdErr
func (m *Messenger) err(icon string, line string, args ...any) {
	fmt.Fprintf(m.ios.Err, (icon + " " + line), args...)
}
