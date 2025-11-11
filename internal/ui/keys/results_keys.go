package keys

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	PickFiles key.Binding
	SwitchSection key.Binding
}

var ResultsKeys = KeyMap{
	PickFiles: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "pick tests to run"),
	),
	SwitchSection: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch section"),
	),
}

type ResultsSectionKeyMap struct {
	LineUp key.Binding
	LineDown key.Binding
	RunAllTests key.Binding
	RunSelectedTest key.Binding
}
var ResultsSectionKeys = ResultsSectionKeyMap{
	LineUp: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("up", "scroll up"),
	),
	LineDown: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("down", "scroll down"),
	),
	RunAllTests: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "run all tests"),
	),
	RunSelectedTest: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "run selected test result"),
	),
}