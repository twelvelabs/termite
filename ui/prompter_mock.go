// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package ui

import (
	"sync"
)

// Ensure, that PrompterMock does implement Prompter.
// If this is not the case, regenerate this file with moq.
var _ Prompter = &PrompterMock{}

// PrompterMock is a mock implementation of Prompter.
//
//	func TestSomethingThatUsesPrompter(t *testing.T) {
//
//		// make and configure a mocked Prompter
//		mockedPrompter := &PrompterMock{
//			ConfirmFunc: func(msg string, value bool, help string) (bool, error) {
//				panic("mock out the Confirm method")
//			},
//			InputFunc: func(msg string, value string, help string) (string, error) {
//				panic("mock out the Input method")
//			},
//			MultiSelectFunc: func(msg string, options []string, values []string, help string) ([]string, error) {
//				panic("mock out the MultiSelect method")
//			},
//			SelectFunc: func(msg string, options []string, value string, help string) (string, error) {
//				panic("mock out the Select method")
//			},
//		}
//
//		// use mockedPrompter in code that requires Prompter
//		// and then make assertions.
//
//	}
type PrompterMock struct {
	// ConfirmFunc mocks the Confirm method.
	ConfirmFunc func(msg string, value bool, help string) (bool, error)

	// InputFunc mocks the Input method.
	InputFunc func(msg string, value string, help string) (string, error)

	// MultiSelectFunc mocks the MultiSelect method.
	MultiSelectFunc func(msg string, options []string, values []string, help string) ([]string, error)

	// SelectFunc mocks the Select method.
	SelectFunc func(msg string, options []string, value string, help string) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// Confirm holds details about calls to the Confirm method.
		Confirm []struct {
			// Msg is the msg argument value.
			Msg string
			// Value is the value argument value.
			Value bool
			// Help is the help argument value.
			Help string
		}
		// Input holds details about calls to the Input method.
		Input []struct {
			// Msg is the msg argument value.
			Msg string
			// Value is the value argument value.
			Value string
			// Help is the help argument value.
			Help string
		}
		// MultiSelect holds details about calls to the MultiSelect method.
		MultiSelect []struct {
			// Msg is the msg argument value.
			Msg string
			// Options is the options argument value.
			Options []string
			// Values is the values argument value.
			Values []string
			// Help is the help argument value.
			Help string
		}
		// Select holds details about calls to the Select method.
		Select []struct {
			// Msg is the msg argument value.
			Msg string
			// Options is the options argument value.
			Options []string
			// Value is the value argument value.
			Value string
			// Help is the help argument value.
			Help string
		}
	}
	lockConfirm     sync.RWMutex
	lockInput       sync.RWMutex
	lockMultiSelect sync.RWMutex
	lockSelect      sync.RWMutex
}

// Confirm calls ConfirmFunc.
func (mock *PrompterMock) Confirm(msg string, value bool, help string) (bool, error) {
	if mock.ConfirmFunc == nil {
		panic("PrompterMock.ConfirmFunc: method is nil but Prompter.Confirm was just called")
	}
	callInfo := struct {
		Msg   string
		Value bool
		Help  string
	}{
		Msg:   msg,
		Value: value,
		Help:  help,
	}
	mock.lockConfirm.Lock()
	mock.calls.Confirm = append(mock.calls.Confirm, callInfo)
	mock.lockConfirm.Unlock()
	return mock.ConfirmFunc(msg, value, help)
}

// ConfirmCalls gets all the calls that were made to Confirm.
// Check the length with:
//
//	len(mockedPrompter.ConfirmCalls())
func (mock *PrompterMock) ConfirmCalls() []struct {
	Msg   string
	Value bool
	Help  string
} {
	var calls []struct {
		Msg   string
		Value bool
		Help  string
	}
	mock.lockConfirm.RLock()
	calls = mock.calls.Confirm
	mock.lockConfirm.RUnlock()
	return calls
}

// Input calls InputFunc.
func (mock *PrompterMock) Input(msg string, value string, help string) (string, error) {
	if mock.InputFunc == nil {
		panic("PrompterMock.InputFunc: method is nil but Prompter.Input was just called")
	}
	callInfo := struct {
		Msg   string
		Value string
		Help  string
	}{
		Msg:   msg,
		Value: value,
		Help:  help,
	}
	mock.lockInput.Lock()
	mock.calls.Input = append(mock.calls.Input, callInfo)
	mock.lockInput.Unlock()
	return mock.InputFunc(msg, value, help)
}

// InputCalls gets all the calls that were made to Input.
// Check the length with:
//
//	len(mockedPrompter.InputCalls())
func (mock *PrompterMock) InputCalls() []struct {
	Msg   string
	Value string
	Help  string
} {
	var calls []struct {
		Msg   string
		Value string
		Help  string
	}
	mock.lockInput.RLock()
	calls = mock.calls.Input
	mock.lockInput.RUnlock()
	return calls
}

// MultiSelect calls MultiSelectFunc.
func (mock *PrompterMock) MultiSelect(msg string, options []string, values []string, help string) ([]string, error) {
	if mock.MultiSelectFunc == nil {
		panic("PrompterMock.MultiSelectFunc: method is nil but Prompter.MultiSelect was just called")
	}
	callInfo := struct {
		Msg     string
		Options []string
		Values  []string
		Help    string
	}{
		Msg:     msg,
		Options: options,
		Values:  values,
		Help:    help,
	}
	mock.lockMultiSelect.Lock()
	mock.calls.MultiSelect = append(mock.calls.MultiSelect, callInfo)
	mock.lockMultiSelect.Unlock()
	return mock.MultiSelectFunc(msg, options, values, help)
}

// MultiSelectCalls gets all the calls that were made to MultiSelect.
// Check the length with:
//
//	len(mockedPrompter.MultiSelectCalls())
func (mock *PrompterMock) MultiSelectCalls() []struct {
	Msg     string
	Options []string
	Values  []string
	Help    string
} {
	var calls []struct {
		Msg     string
		Options []string
		Values  []string
		Help    string
	}
	mock.lockMultiSelect.RLock()
	calls = mock.calls.MultiSelect
	mock.lockMultiSelect.RUnlock()
	return calls
}

// Select calls SelectFunc.
func (mock *PrompterMock) Select(msg string, options []string, value string, help string) (string, error) {
	if mock.SelectFunc == nil {
		panic("PrompterMock.SelectFunc: method is nil but Prompter.Select was just called")
	}
	callInfo := struct {
		Msg     string
		Options []string
		Value   string
		Help    string
	}{
		Msg:     msg,
		Options: options,
		Value:   value,
		Help:    help,
	}
	mock.lockSelect.Lock()
	mock.calls.Select = append(mock.calls.Select, callInfo)
	mock.lockSelect.Unlock()
	return mock.SelectFunc(msg, options, value, help)
}

// SelectCalls gets all the calls that were made to Select.
// Check the length with:
//
//	len(mockedPrompter.SelectCalls())
func (mock *PrompterMock) SelectCalls() []struct {
	Msg     string
	Options []string
	Value   string
	Help    string
} {
	var calls []struct {
		Msg     string
		Options []string
		Value   string
		Help    string
	}
	mock.lockSelect.RLock()
	calls = mock.calls.Select
	mock.lockSelect.RUnlock()
	return calls
}
