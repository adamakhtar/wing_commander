package results

import (
	"github.com/adamakhtar/wing_commander/internal/ui/context"
	"github.com/adamakhtar/wing_commander/internal/ui/keys"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

//
// TYPES
//================================================

type Model struct {
	ctx context.Context
}

//
// BUILDERS
//================================================

func NewModel(ctx context.Context) Model {
	model := Model{
		ctx: ctx,
	}
	return model
}

//
// BUBBLETEA
//================================================

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.ResultsKeys.PickFiles):
			return m, switchToFilePickerCmd
		}
	}
	return m, nil
}

func (m Model) View() string {
	return "Results Screen!"
}

//
// MESSAGES & HANDLERS
//================================================

type OpenFilePickerMsg struct{}

//
// COMMANDS
//================================================

func switchToFilePickerCmd() tea.Msg {
	return OpenFilePickerMsg{}
}

// EXTERNAL FUNCTIONS
//================================================

func (m *Model) UpdateContext(ctx context.Context) {
	m.ctx = ctx
}