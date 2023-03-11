package ui

// Responder is a function that returns stubbed prompt responses.
type Responder func(prompt Prompt) (response any, err error)

// RespondBool creates a responder that returns the given bool.
func RespondBool(response bool) Responder {
	return func(p Prompt) (any, error) {
		return response, nil
	}
}

// RespondString creates a responder that returns the given string.
func RespondString(response string) Responder {
	return func(p Prompt) (any, error) {
		return response, nil
	}
}

// RespondStringSlice creates a responder that returns the given slice of strings.
func RespondStringSlice(response []string) Responder {
	return func(p Prompt) (any, error) {
		return response, nil
	}
}

// RespondError creates a responder that returns the given error.
func RespondError(err error) Responder {
	return func(p Prompt) (any, error) {
		return nil, err
	}
}

// RespondDefault creates a responder that returns the prompt's default value.
func RespondDefault() Responder {
	return func(p Prompt) (any, error) {
		return p.Value, nil
	}
}
