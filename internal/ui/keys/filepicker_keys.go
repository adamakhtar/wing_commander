package keys

import (
	"github.com/charmbracelet/bubbles/key"
)

type FilepickerKeyMap struct {
	Cancel key.Binding
	LineUp key.Binding
	LineDown key.Binding
	Select key.Binding
	ConfirmSelection key.Binding
}

var FilepickerKeys = FilepickerKeyMap{
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	LineUp: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("up", "scroll up"),
	),
	LineDown: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("down", "scroll down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select file"),
	),
	ConfirmSelection: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "run tests"),
	),
}
