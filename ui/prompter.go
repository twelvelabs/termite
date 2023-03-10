package ui

// Prompter is an interface for types that prompt for user input.
type Prompter interface {
	// Confirm prompts for a boolean yes/no value.
	Confirm(msg string, value bool, opts ...PromptOpt) (bool, error)
	// Input prompts for single string value.
	Input(msg string, value string, opts ...PromptOpt) (string, error)
	// MultiSelect prompts for a slice of string values w/ a fixed set of options.
	MultiSelect(msg string, options []string, values []string, opts ...PromptOpt) ([]string, error)
	// Select prompts for single string value w/ a fixed set of options.
	Select(msg string, options []string, value string, opts ...PromptOpt) (string, error)
}

// PromptParams holds optional params for Prompter methods.
type PromptParams struct {
	// Help text to show when the user presses "?".
	Help string
	// One or more validation rules to be passed
	// to [validate.KeyVal].
	ValidationRules string
}

// PromptOpt allows setting optional prompt params.
type PromptOpt func(params *PromptParams)

// WithHelp sets the help text for a prompt.
func WithHelp(text string) PromptOpt {
	return func(params *PromptParams) {
		params.Help = text
	}
}

// WithValidation sets the validation rules for a prompt.
func WithValidation(rules string) PromptOpt {
	return func(params *PromptParams) {
		params.ValidationRules = rules
	}
}

// GetPromptParams returns the params for the given prompt ops.
func GetPromptParams(opts ...PromptOpt) *PromptParams {
	params := &PromptParams{}
	for _, opt := range opts {
		opt(params)
	}
	return params
}
