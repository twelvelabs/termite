package ui

//go:generate moq -rm -out prompter_test.go . Prompter

// Prompter is an interface for types that prompt for user input.
type Prompter interface {
	// Confirm prompts for a boolean yes/no value.
	Confirm(msg string, value bool, help string) (bool, error)
	// Input prompts for single string value.
	Input(msg string, value string, help string) (string, error)
	// MultiSelect prompts for a slice of string values w/ a fixed set of options.
	MultiSelect(msg string, options []string, values []string, help string) ([]string, error)
	// Select prompts for single string value w/ a fixed set of options.
	Select(msg string, options []string, value string, help string) (string, error)
}
