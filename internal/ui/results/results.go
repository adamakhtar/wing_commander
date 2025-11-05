package results

import (
	"github.com/adamakhtar/wing_commander/internal/runner"
	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/adamakhtar/wing_commander/internal/ui/context"
	"github.com/adamakhtar/wing_commander/internal/ui/keys"
	"github.com/adamakhtar/wing_commander/internal/ui/results/previewsection"
	"github.com/adamakhtar/wing_commander/internal/ui/results/resultssection"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

//
// TYPES
//================================================

type Model struct {
	ctx *context.Context
	testRuns TestRuns
	testRunner *runner.TestRunner
	testExecutionResult *runner.TestExecutionResult
	resultsSection resultssection.Model
	previewSection previewsection.Model
	width int
	height int
	error error
}

//
// BUILDERS
//================================================

func NewModel(ctx *context.Context) Model {
	testRunner := runner.NewTestRunner(ctx.Config)

	model := Model{
		ctx: ctx,
		testRunner: testRunner,
		testRuns: NewTestRuns(),
		resultsSection: resultssection.NewModel(ctx),
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

	selectedTestResult := m.GetSelectedTestResultId()
	m.previewSection.SetTestResult(selectedTestResult)

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
	screen := lipgloss.JoinHorizontal(lipgloss.Top, m.resultsSection.View(), m.previewSection.View())

	return screen
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

		log.Debug("Executing test run", "testRunId", testRunId, "filepaths", testRun.Filepaths, "testResultsPath", m.ctx.Config.TestResultsPath)

		testExecutionResult, err := m.testRunner.ExecuteTests(testRunId, testRun.Filepaths, m.ctx.Config.TestResultsPath)
		if err != nil {
			log.Debug("Error executing test run", "testRunId", testRunId, "error", err)
			return TestExecutionFailedMsg{TestRunId: testRunId, error: err}
		}
		log.Debug("Test execution completed", "testRunId", testRunId, "testExecutionResult", testExecutionResult)
		return TestExecutionCompletedMsg{TestRunId: testRunId, TestExecutionResult: testExecutionResult}
	}
}

// EXTERNAL FUNCTIONS
//================================================

func (m *Model) AddTestRun(filepaths []string) (TestRun, error) {
	testRun, err := m.testRuns.Add(filepaths)
	return testRun, err
}

func (m Model) GetSelectedTestResultId() *types.TestResult {
	testResultId := m.resultsSection.GetSelectedTestResultId()

	if testResultId == -1 {
		return nil
	}

	// TODO extract this to a TestResultCollection type that has a GetById method
	for _, testResult := range m.testExecutionResult.TestResults {
		if testResult.Id == testResultId {
			return &testResult
		}
	}
	return nil
}


func (m *Model) handleTestExecutionCompletion(testExecutionResult *runner.TestExecutionResult) {
	log.Debug("Test execution completed", "testExecutionResult", testExecutionResult)
	m.testExecutionResult = testExecutionResult
	m.resultsSection.SetRows(testExecutionResult)
}

func (m *Model) Prepare() tea.Cmd {
	return nil
}

func (m *Model) SetSize(width int, height int) {
	m.width = width
	m.height = height

	m.resultsSection.SetSize(m.width / 2 , m.height)
	m.previewSection.SetSize(m.width / 2 , m.height)
}