package ui

import (
	"github.com/adamakhtar/wing_commander/internal/ui/context"
	"github.com/adamakhtar/wing_commander/internal/ui/filepicker"
	"github.com/adamakhtar/wing_commander/internal/ui/keys"
	"github.com/adamakhtar/wing_commander/internal/ui/results"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

//
// TYPES
//================================================

type Model struct {
	// Data
	ctx context.Context
	ready bool
	resultsScreen results.Model
	filepickerScreen filepicker.Model
}

//
// BUILDERS
//================================================

func NewModel() Model {
	ctx := context.Context{
		CurrentScreen: context.ResultsScreen,
	}

	model := Model{
		ready: false,
		ctx: ctx,
	}

	model.resultsScreen = results.NewModel(model.ctx)

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
	case tea.WindowSizeMsg:
		m.onWindowResize(msg)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.AppKeys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.ResultsKeys.PickFiles):
			m.ctx.CurrentScreen = context.FilePickerScreen
			return m, nil
		}
	}

	m.syncContext()

	return m, nil
}

func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	currentScreen := m.getCurrentScreen()
	if currentScreen == nil {
		return "Error: No screen found"
	}

	return currentScreen.View()
}

//
// MESSAGES & HANDLERS
//================================================

func (m *Model) onWindowResize(msg tea.WindowSizeMsg) {
	m.ctx.ScreenWidth = msg.Width
	m.ctx.ScreenHeight = msg.Height

	m.ready = true
}

// INTERNAL FUNCTIONS
//================================================

func (m Model) getCurrentScreen() tea.Model {
	switch m.ctx.CurrentScreen {
	case context.ResultsScreen:
		return m.resultsScreen
	case context.FilePickerScreen:
		return m.filepickerScreen
	default:
		return nil
	}
}

func (m *Model) syncContext() {
	m.resultsScreen.UpdateContext(m.ctx)
}