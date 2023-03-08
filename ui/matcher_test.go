package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatcher(t *testing.T) {
	tests := []struct {
		desc    string
		matcher Matcher
		prompt  Prompt
		matches bool
	}{
		{
			desc:    "MatchAny: match",
			matcher: MatchAny,
			prompt:  Prompt{},
			matches: true,
		},

		{
			desc:    "MatchConfirm: match",
			matcher: MatchConfirm("Proceed?"),
			prompt: Prompt{
				Type:    PromptTypeConfirm,
				Message: "Proceed?",
			},
			matches: true,
		},
		{
			desc:    "MatchConfirm: fail message",
			matcher: MatchConfirm("Proceed?"),
			prompt: Prompt{
				Type:    PromptTypeConfirm,
				Message: "Nope?",
			},
			matches: false,
		},
		{
			desc:    "MatchConfirm: fail type",
			matcher: MatchConfirm("Proceed?"),
			prompt: Prompt{
				Type:    PromptTypeSelect,
				Message: "Proceed?",
			},
			matches: false,
		},

		{
			desc:    "MatchInput: match",
			matcher: MatchInput("Name:"),
			prompt: Prompt{
				Type:    PromptTypeInput,
				Message: "Name:",
			},
			matches: true,
		},
		{
			desc:    "MatchInput: fail message",
			matcher: MatchInput("Name:"),
			prompt: Prompt{
				Type:    PromptTypeInput,
				Message: "Nope:",
			},
			matches: false,
		},
		{
			desc:    "MatchInput: fail type",
			matcher: MatchInput("Name:"),
			prompt: Prompt{
				Type:    PromptTypeSelect,
				Message: "Name:",
			},
			matches: false,
		},

		{
			desc:    "MatchMultiSelect: match",
			matcher: MatchMultiSelect("Country:"),
			prompt: Prompt{
				Type:    PromptTypeMultiSelect,
				Message: "Country:",
			},
			matches: true,
		},
		{
			desc:    "MatchMultiSelect: fail message",
			matcher: MatchMultiSelect("Country:"),
			prompt: Prompt{
				Type:    PromptTypeMultiSelect,
				Message: "Nope:",
			},
			matches: false,
		},
		{
			desc:    "MatchMultiSelect: fail type",
			matcher: MatchMultiSelect("Country:"),
			prompt: Prompt{
				Type:    PromptTypeConfirm,
				Message: "Country:",
			},
			matches: false,
		},

		{
			desc:    "MatchSelect: match",
			matcher: MatchSelect("Language:"),
			prompt: Prompt{
				Type:    PromptTypeSelect,
				Message: "Language:",
			},
			matches: true,
		},
		{
			desc:    "MatchSelect: fail message",
			matcher: MatchSelect("Language:"),
			prompt: Prompt{
				Type:    PromptTypeSelect,
				Message: "Nope:",
			},
			matches: false,
		},
		{
			desc:    "MatchSelect: fail type",
			matcher: MatchSelect("Language:"),
			prompt: Prompt{
				Type:    PromptTypeConfirm,
				Message: "Language:",
			},
			matches: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			actual := tt.matcher(tt.prompt)
			assert.Equal(t, tt.matches, actual)
		})
	}
}
