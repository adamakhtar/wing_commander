package results

import (
	"github.com/adamakhtar/wing_commander/internal/ui/context"
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
	return m, nil
}

func (m Model) View() string {
	return "Results Screen!"
}

//
// EXTERNAL FUNCTIONS
//================================================

func (m *Model) UpdateContext(ctx context.Context) {
	m.ctx = ctx
}