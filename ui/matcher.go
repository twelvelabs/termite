package ui

// Matcher is function that matches prompts.
type Matcher func(p Prompt) bool

// MatchAny is a matcher that matches any prompt.
func MatchAny(p Prompt) bool {
	return true
}

// MatchConfirm returns a matcher that matches confirm prompts.
func MatchConfirm(msg string) Matcher {
	return func(p Prompt) bool {
		return p.Type == PromptTypeConfirm && p.Message == msg
	}
}

// MatchInput returns a matcher that matches input prompts.
func MatchInput(msg string) Matcher {
	return func(p Prompt) bool {
		return p.Type == PromptTypeInput && p.Message == msg
	}
}

// MatchMultiSelect returns a matcher that matches multi-select prompts.
func MatchMultiSelect(msg string) Matcher {
	return func(p Prompt) bool {
		return p.Type == PromptTypeMultiSelect && p.Message == msg
	}
}

// MatchSelect returns a matcher that matches select prompts.
func MatchSelect(msg string) Matcher {
	return func(p Prompt) bool {
		return p.Type == PromptTypeSelect && p.Message == msg
	}
}
