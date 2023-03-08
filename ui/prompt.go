package ui

import "fmt"

type PromptType int

func (pt PromptType) String() string {
	switch pt {
	case PromptTypeConfirm:
		return "Confirm"
	case PromptTypeInput:
		return "Input"
	case PromptTypeMultiSelect:
		return "MultiSelect"
	case PromptTypeSelect:
		return "Select"
	default:
		return "Unknown"
	}
}

var (
	PromptTypeUnknown     PromptType = 0
	PromptTypeConfirm     PromptType = 1
	PromptTypeInput       PromptType = 2
	PromptTypeMultiSelect PromptType = 3
	PromptTypeSelect      PromptType = 4
)

type Prompt struct {
	Type    PromptType
	Message string
	Help    string
	Value   any
	Options any
}

func (p Prompt) String() string {
	return fmt.Sprintf("<%s msg=%s>", p.Type.String(), p.Message)
}
