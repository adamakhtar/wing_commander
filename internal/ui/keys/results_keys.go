package keys

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	PickFiles key.Binding
}

var ResultsKeys = KeyMap{
	PickFiles: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "pick tests to run"),
	),
}
