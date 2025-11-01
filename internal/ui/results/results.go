package results

import (
	"strings"

	"github.com/adamakhtar/wing_commander/internal/runner"
	"github.com/adamakhtar/wing_commander/internal/ui/context"
	"github.com/adamakhtar/wing_commander/internal/ui/keys"
	"github.com/adamakhtar/wing_commander/internal/ui/results/resultssection"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

//
// TYPES
//================================================

type Model struct {
	ctx context.Context
	testRuns TestRuns
	testRunner *runner.TestRunner
	resultsSection resultssection.Model
	width int
	height int
	error error
}

//
// BUILDERS
//================================================

func NewModel(ctx context.Context) Model {
	testRunner := runner.NewTestRunner(ctx.Config)

	model := Model{
		ctx: ctx,
		testRunner: testRunner,
		testRuns: NewTestRuns(),
		resultsSection: resultssection.NewModel(),
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
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case TestExecutionCompletedMsg:
		m.handleTestExecutionCompletion(msg.TestExecutionResult)
		return m, nil
	case TestExecutionFailedMsg:
		m.error = msg.error
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.ResultsKeys.PickFiles):
			return m, switchToFilePickerCmd
		}
	}

	resultsSection, cmd := m.resultsSection.Update(msg)
	m.resultsSection = resultsSection.(resultssection.Model)

	return m, cmd
}

func (m Model) View() string {
	sb := strings.Builder{}
	sb.WriteString("Results Screen!\n")

	sb.WriteString(m.resultsSection.View())
	return sb.String()
}

//
// MESSAGES & HANDLERS
//================================================

type OpenFilePickerMsg struct{}
type TestRunExecutedMsg struct {
	TestRunId int
}

type TestExecutionCompletedMsg struct {
	TestRunId int
	TestExecutionResult    *runner.TestExecutionResult
}

type TestExecutionFailedMsg struct {
	TestRunId int
	error error
}

//
// COMMANDS
//================================================

func switchToFilePickerCmd() tea.Msg {
	return OpenFilePickerMsg{}
}

func (m Model) ExecuteTestRunCmd(testRunId int) tea.Cmd {
	return func () tea.Msg {
		testRun, err := m.testRuns.Get(testRunId)
		if err != nil {
			return TestExecutionFailedMsg{TestRunId: testRunId, error: err}
		}

		testExecutionResult, err := m.testRunner.ExecuteTests(testRunId, testRun.Filepaths, m.ctx.Config.TestResultsPath)
		if err != nil {
			return TestExecutionFailedMsg{TestRunId: testRunId, error: err}
		}
		return TestExecutionCompletedMsg{TestRunId: testRunId, TestExecutionResult: testExecutionResult}
	}
}

// EXTERNAL FUNCTIONS
//================================================

func (m *Model) AddTestRun(filepaths []string) (TestRun, error) {
	testRun, err := m.testRuns.Add(filepaths)
	return testRun, err
}

func (m *Model) handleTestExecutionCompletion(testExecutionResult *runner.TestExecutionResult) {
	m.resultsSection.SetRows(testExecutionResult)
}

func (m *Model) UpdateContext(ctx context.Context) {
	m.ctx = ctx
}

func (m *Model) Prepare() tea.Cmd {
	return nil
}

func (m *Model) SetSize(width int, height int) {
	m.width = width
	m.height = height

	m.resultsSection.SetSize(m.width, m.height)
}