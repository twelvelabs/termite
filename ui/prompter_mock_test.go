package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Dummy tests so we can export the mock w/out taking a coverage hit.
// We're just going to assume that `moq` generates working code :shrug:.
func TestPrompterMock(t *testing.T) {
	mock := NewPrompterMock()

	assert.Panics(t, func() {
		_, _ = mock.Confirm("", false, "")
	})
	assert.Panics(t, func() {
		_, _ = mock.Input("", "", "")
	})
	assert.Panics(t, func() {
		_, _ = mock.MultiSelect("", nil, nil, "")
	})
	assert.Panics(t, func() {
		_, _ = mock.Select("", nil, "", "")
	})

	mock = &PrompterMock{
		ConfirmFunc: func(msg string, value bool, help string) (bool, error) {
			return false, nil
		},
		InputFunc: func(msg, value, help string) (string, error) {
			return "", nil
		},
		MultiSelectFunc: func(msg string, options, values []string, help string) ([]string, error) {
			return []string{}, nil
		},
		SelectFunc: func(msg string, options []string, value, help string) (string, error) {
			return "", nil
		},
	}
	_, _ = mock.Confirm("", false, "")
	_, _ = mock.Input("", "", "")
	_, _ = mock.MultiSelect("", nil, nil, "")
	_, _ = mock.Select("", nil, "", "")

	assert.Equal(t, 1, len(mock.ConfirmCalls()))
	assert.Equal(t, 1, len(mock.InputCalls()))
	assert.Equal(t, 1, len(mock.MultiSelectCalls()))
	assert.Equal(t, 1, len(mock.SelectCalls()))
}
