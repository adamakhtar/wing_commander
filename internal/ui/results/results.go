package results

import (
	"fmt"

	"github.com/adamakhtar/wing_commander/internal/runner"
	"github.com/adamakhtar/wing_commander/internal/testrun"
	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/adamakhtar/wing_commander/internal/ui/context"
	"github.com/adamakhtar/wing_commander/internal/ui/filepicker"
	"github.com/adamakhtar/wing_commander/internal/ui/keys"
	"github.com/adamakhtar/wing_commander/internal/ui/results/previewsection"
	"github.com/adamakhtar/wing_commander/internal/ui/results/resultssection"
	"github.com/adamakhtar/wing_commander/internal/ui/results/testrunssection"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

//
// TYPES
//================================================

type Model struct {
	ctx                 *context.Context
	testRuns            testrun.TestRuns
	testRunner          *runner.TestRunner
	testExecutionResult *runner.TestExecutionResult
	resultsSection      resultssection.Model
	previewSection      previewsection.Model
	testRunsSection     testrunssection.Model
	width               int
	height              int
	error               error
}

//
// BUILDERS
//================================================

func NewModel(ctx *context.Context) Model {
	testRunner := runner.NewTestRunner(ctx.Config)
	testRuns := testrun.NewTestRuns()

	model := Model{
		ctx:             ctx,
		testRunner:      testRunner,
		testRuns:        testRuns,
		resultsSection:  resultssection.NewModel(ctx, true),
		previewSection:  previewsection.NewModel(ctx, false),
		testRunsSection: testrunssection.NewModel(ctx, &testRuns),
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
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case resultssection.RunTestMsg:

		testRun, err := m.testRuns.Add(
			[]string{msg.TestPattern},
			testrun.ModeReRunSingleFailure,
		)
		if err != nil {
			// TODO - handle error
			return m, nil
		}
		return m, m.ExecuteTestRunCmd(testRun.Id)
	case filepicker.TestsSelectedMsg:
		m.ctx.CurrentScreen = context.ResultsScreen
		// TODO - consider running a command here that the results screen listens to and it then
		//  performs the test run
		testRun, err := m.testRuns.Add(msg.Filepaths, testrun.ModeRunSelectedPatterns)
		if err != nil {
			// TODO - handle error
			return m, nil
		}

		return m, m.ExecuteTestRunCmd(testRun.Id)
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
		case key.Matches(msg, keys.ResultsKeys.SwitchSection):
			m.resultsSection.ToggleFocus()
			m.previewSection.ToggleFocus()
		case key.Matches(msg, keys.ResultsKeys.RunAllTests):
			testRun, err := m.testRuns.Add([]string{""}, testrun.ModeRunWholeSuite)
			if err != nil {
				// TODO - handle error
				return m, nil
			}
			return m, m.ExecuteTestRunCmd(testRun.Id)
		case key.Matches(msg, keys.ResultsKeys.RunFailedTests):
			testRun, err := m.AddTestRunForFailedTests()
			if err != nil {
				// TODO - handle error
				return m, nil
			}
			return m, m.ExecuteTestRunCmd(testRun.Id)
		}
	}

	resultsSection, resultsSectionCmd := m.resultsSection.Update(msg)
	m.resultsSection = resultsSection.(resultssection.Model)
	cmds = append(cmds, resultsSectionCmd)

	previewSection, previewSectionCmd := m.previewSection.Update(msg)
	m.previewSection = previewSection.(previewsection.Model)
	cmds = append(cmds, previewSectionCmd)

	selectedTestResult := m.GetSelectedTestResultId()
	m.previewSection.SetTestResult(selectedTestResult)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	screen := lipgloss.JoinHorizontal(lipgloss.Top, m.testRunsSection.View(), m.resultsSection.View(), m.previewSection.View())

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
	TestRunId           int
	TestExecutionResult *runner.TestExecutionResult
}

type TestExecutionFailedMsg struct {
	TestRunId int
	error     error
}

//
// COMMANDS
//================================================

func switchToFilePickerCmd() tea.Msg {
	return OpenFilePickerMsg{}
}

func (m Model) ExecuteTestRunCmd(testRunId int) tea.Cmd {
	return func() tea.Msg {
		testRun, err := m.testRuns.Get(testRunId)
		if err != nil {
			return TestExecutionFailedMsg{TestRunId: testRunId, error: err}
		}

		log.Debugf("Executing tests for test run %d: %v", testRunId, testRun.Filepaths)

		testExecutionResult, err := m.testRunner.ExecuteTests(testRun)
		if err != nil {
			log.Debugf("Failed to execute tests for test run %d: %v", testRunId, err)
			return TestExecutionFailedMsg{TestRunId: testRunId, error: err}
		}
		log.Debugf("Tests executed for test run %d: %v, metrics: %v", testRunId, len(testExecutionResult.TestResults), testExecutionResult.Metrics)

		return TestExecutionCompletedMsg{TestRunId: testRunId, TestExecutionResult: testExecutionResult}
	}
}

// EXTERNAL FUNCTIONS
//================================================

func (m *Model) AddTestRunForFailedTests() (testrun.TestRun, error) {
	// TODO - create a TestResultsCollection type and move this logic to that type
	var filepaths []string

	if m.testExecutionResult == nil {
		return testrun.TestRun{}, fmt.Errorf("no previous test execution available")
	}

	for _, testResult := range m.testExecutionResult.FailedTests {
		filepaths = append(filepaths, testResult.TestFilePath)
	}

	testRun, err := m.testRuns.Add(filepaths, testrun.ModeReRunAllFailures)
	if err != nil {
		return testrun.TestRun{}, err
	}
	return testRun, nil
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
	m.testExecutionResult = testExecutionResult
	m.resultsSection.SetRows(testExecutionResult)
}

func (m *Model) Prepare() tea.Cmd {
	return nil
}

func (m *Model) SetSize(width int, height int) {
	m.width = width
	m.height = height

	m.testRunsSection.SetSize(m.width*2/10, m.height)
	m.resultsSection.SetSize(m.width*4/10, m.height)
	m.previewSection.SetSize(m.width*4/10, m.height)
}
