package keys

import (
	"github.com/charmbracelet/bubbles/key"
)

type AppKeysMap struct {
	Quit key.Binding
}

var AppKeys = AppKeysMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+q"),
		key.WithHelp("ctrl+q", "quit"),
	),
};
