package ui

import (
	"github.com/adamakhtar/wing_commander/internal/config"
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

func NewModel(cfg *config.Config) tea.Model {
	ctx := context.Context{
		CurrentScreen: context.ResultsScreen,
		Config: cfg,
	}

	model := Model{
		ready: false,
		ctx: ctx,
	}

	model.resultsScreen = results.NewModel(model.ctx)
	model.filepickerScreen = filepicker.NewModel(model.ctx)

	return model
}

//
// BUBBLETEA
//================================================

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case results.OpenFilePickerMsg:
		m.ctx.CurrentScreen = context.FilePickerScreen
		cmd = m.filepickerScreen.Prepare()
		return m, cmd
	case filepicker.CancelMsg:
		m.ctx.CurrentScreen = context.ResultsScreen
		cmd = m.resultsScreen.Prepare()
		return m, cmd
	case filepicker.TestsSelectedMsg:
		m.ctx.CurrentScreen = context.ResultsScreen
		// TODO - consider running a command here that the results screen listens to and it then
		//  performs the test run
		testRun, err := m.resultsScreen.AddTestRun(msg.Filepaths)
		if err != nil {
			// TODO - handle error
			return m, nil
		}

		cmd = m.resultsScreen.ExecuteTestRunCmd(testRun.Id)
		return m, cmd
	case tea.WindowSizeMsg:
		m.ctx.ScreenWidth = msg.Width
		m.ctx.ScreenHeight = msg.Height
		m.ready = true
	case tea.KeyMsg:
		if key.Matches(msg, keys.AppKeys.Quit) {
			return m, tea.Quit
		}
	}

	// process child screen messages
	switch m.ctx.CurrentScreen {
	case context.ResultsScreen:
		resultsScreen, resultsCmd := m.resultsScreen.Update(msg)
		m.resultsScreen = resultsScreen.(results.Model)
		cmd = resultsCmd
	case context.FilePickerScreen:
		filepickerScreen, filepickerCmd := m.filepickerScreen.Update(msg)
		m.filepickerScreen = filepickerScreen.(filepicker.Model)
		cmd = filepickerCmd
	}

	m.syncContext()

	return m, cmd
}

func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	currentScreen := m.getCurrentScreen()
	if currentScreen == nil {
		return "Error: No screen found"
	} else {
		return currentScreen.View()
	}
}

//
// MESSAGES & HANDLERS
//================================================

func (m *Model) onWindowResize(msg tea.WindowSizeMsg) {
	m.ctx.ScreenWidth = msg.Width
	m.ctx.ScreenHeight = msg.Height

	m.ready = true
}

//
// COMMANDS
//================================================


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
	m.filepickerScreen.UpdateContext(m.ctx)
}