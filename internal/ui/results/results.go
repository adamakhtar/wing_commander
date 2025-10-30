package results

import (
	"fmt"
	"strings"

	"github.com/adamakhtar/wing_commander/internal/runner"
	"github.com/adamakhtar/wing_commander/internal/ui/context"
	"github.com/adamakhtar/wing_commander/internal/ui/keys"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

//
// TYPES
//================================================

type Model struct {
	ctx context.Context
	testRuns TestRuns
	testRunner *runner.TestRunner
	testExecutionResult *runner.TestExecutionResult
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
	case TestExecutionCompletedMsg:
		m.recordTestExecutionResult(msg.TestRunId, msg.TestExecutionResult)
	case TestExecutionFailedMsg:
		m.error = msg.error
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.ResultsKeys.PickFiles):
			return m, switchToFilePickerCmd
		}
	}
	return m, nil
}

func (m Model) View() string {
	sb := strings.Builder{}
	sb.WriteString("Results Screen!\n")

	testRun, ok := m.testRuns.MostRecent()
	if ok {
		sb.WriteString(fmt.Sprintf("Test Run: %d\n", testRun.Id))
		sb.WriteString(fmt.Sprintf("Filepaths: %v\n", strings.Join(testRun.Filepaths, ", ")))
	}
	if m.error != nil {
		log.Debug("Results Screen", "error", m.error)
		sb.WriteString(fmt.Sprintf("Error: %v\n", m.error.Error()))
	}
	if m.testExecutionResult != nil {
		sb.WriteString(fmt.Sprintf("Test Execution Result: %v\n", m.testExecutionResult.Metrics))
		sb.WriteString(fmt.Sprintf("Test Execution Result: %v\n", m.testExecutionResult.ExecutionTime))
	}
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

func (m *Model) recordTestExecutionResult(testRunId int, testExecutionResult *runner.TestExecutionResult) {
	m.testExecutionResult = testExecutionResult
}

func (m *Model) UpdateContext(ctx context.Context) {
	m.ctx = ctx
}

func (m *Model) Prepare() tea.Cmd {
	return nil
}