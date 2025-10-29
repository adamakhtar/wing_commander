package results

import (
	"fmt"
	"strings"

	"github.com/adamakhtar/wing_commander/internal/runner"
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
	testRuns TestRuns
	testRunner *runner.TestRunner
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
		m.recordTestRunResult(msg.TestRunId, msg.Result)
	case TestExecutionFailedMsg:
		m.recordTestRunError(msg.TestRunId, msg.Error)
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
		sb.WriteString(fmt.Sprintf("Error: %v\n", testRun.Error))
		sb.WriteString(fmt.Sprintf("Filepaths: %v\n", strings.Join(testRun.Filepaths, ", ")))
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
	Result    *runner.TestExecutionResult
}

type TestExecutionFailedMsg struct {
	TestRunId int
	Error string
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
			return TestExecutionFailedMsg{TestRunId: testRunId, Error: err.Error()}
		}

		result, err := m.testRunner.ExecuteTests(testRun.Filepaths)
		if err != nil {
			return TestExecutionFailedMsg{TestRunId: testRunId, Error: err.Error()}
		}
		return TestExecutionCompletedMsg{TestRunId: testRunId, Result: result}
	}
}

// EXTERNAL FUNCTIONS
//================================================

func (m *Model) AddTestRun(filepaths []string) (TestRun, error) {
	testRun, err := m.testRuns.Add(filepaths)
	return testRun, err
}

func (m *Model) recordTestRunResult(testRunId int, result *runner.TestExecutionResult) {
	_, err := m.testRuns.UpdateResult(testRunId, result)
	if err != nil {
		// TODO - handle error
		return
	}
}

func (m *Model) recordTestRunError(testRunId int, error string) {
	_, err := m.testRuns.UpdateError(testRunId, error)
	if err != nil {
		// TODO - handle error
		return
	}
}

func (m *Model) UpdateContext(ctx context.Context) {
	m.ctx = ctx
}

func (m *Model) Prepare() tea.Cmd {
	return nil
}