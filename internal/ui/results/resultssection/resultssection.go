package resultssection

import (
	"sort"
	"strings"

	"github.com/adamakhtar/wing_commander/internal/runner"
	"github.com/adamakhtar/wing_commander/internal/testrun"
	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/adamakhtar/wing_commander/internal/ui/context"
	"github.com/adamakhtar/wing_commander/internal/ui/keys"
	"github.com/adamakhtar/wing_commander/internal/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

var customBorder = table.Border{
	Top:            "─",
	Bottom:         "─",
	Left:           "│",
	Right:          "│",
	TopLeft:        "┌",
	TopRight:       "┐",
	BottomLeft:     "└",
	BottomRight:    "┘",
	TopJunction:    "┬",
	LeftJunction:   "├",
	RightJunction:  "┤",
	BottomJunction: "┴",
	InnerJunction:  "┼",
	InnerDivider:   "│",
}

const (
	columnKeyTestName        = "test_name"
	columnKeyTestResult    = "test_result"
	columnKeyMetaId          = "test_id"
	columnKeyMetaTestPattern = "test_pattern"
)

const (
	paddingX = 1
	paddingY = 0
)

type Model struct {
	ctx          *context.Context
	focus        bool
	resultsTable table.Model
	width        int
	height       int
}

//
// BUILDERS
//================================================

func NewModel(ctx *context.Context, focus bool) Model {
	columns := getColumnConfiguration(3, 15, 15, &ctx.Styles)

	resultsTable := table.New(columns).
		Focused(true).
		WithBaseStyle(ctx.Styles.ResultsSection.TableBaseStyle).
		HighlightStyle(ctx.Styles.ResultsSection.TableHighlight).
		Border(customBorder)

	return Model{
		ctx:          ctx,
		focus:        focus,
		resultsTable: resultsTable,
	}
}

//
// BUBBLETEA
//================================================

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.isBlurred() {
		return m, nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	// case table.UserEventHighlightedIndexChanged:
	// 	m.handleResultTableHighlightedRowChange(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.ResultsSectionKeys.LineUp):
			m.resultsTable, cmd = m.resultsTable.Update(msg)
		case key.Matches(msg, keys.ResultsSectionKeys.LineDown):
			m.resultsTable, cmd = m.resultsTable.Update(msg)
		case key.Matches(msg, keys.ResultsSectionKeys.RunSelectedTest):
			testPattern, ok := m.getSelectedTestPattern()
			if ok {
				cmd = runTestCmd(testPattern)
			}
		}
	}
	return m, cmd
}

func (m Model) View() string {
	panelStyle := m.ctx.Styles.Border.Padding(0, 1).Height(m.height).Width(m.width)
	if m.isFocused() {
		panelStyle = panelStyle.Inherit(m.ctx.Styles.BorderActive)
	} else {
		panelStyle = panelStyle.Inherit(m.ctx.Styles.BorderMuted)
	}

	return panelStyle.Render(m.resultsTable.View())
}

//
// MESSAGES & HANDLERS
//================================================

type RunTestMsg struct {
	TestPattern testrun.TestPattern
}

//
// COMMANDS
//================================================

func runTestCmd(testPattern testrun.TestPattern) tea.Cmd {
	return func() tea.Msg {
		return RunTestMsg{TestPattern: testPattern}
	}
}

//
// EXTERNAL FUNCTIONS
//================================================

func (m *Model) SetSize(width int, height int) {
	m.width = width
	m.height = height
	m.resultsTable = m.resultsTable.WithTargetWidth(m.width - 2*paddingX)
	m.resultsTable = m.resultsTable.WithMinimumHeight(m.height - 2*paddingY)

	failureCauseWidth := 3
	groupNameWidth := 15
	testNameWidth := width - failureCauseWidth - groupNameWidth

	columns := getColumnConfiguration(failureCauseWidth, groupNameWidth, testNameWidth, &m.ctx.Styles)
	m.resultsTable.WithColumns(columns)
}

func (m *Model) SetRows(testExecutionResult *runner.TestExecutionResult) {
	if testExecutionResult == nil {
		m.resultsTable = m.resultsTable.WithRows([]table.Row{})
		return
	}

	results := make([]types.TestResult, len(testExecutionResult.TestResults))
	copy(results, testExecutionResult.TestResults)

	sort.SliceStable(results, func(i, j int) bool {
		pi := sortPriority(results[i])
		pj := sortPriority(results[j])
		if pi != pj {
			return pi < pj
		}

		groupI := strings.ToLower(results[i].GroupName)
		groupJ := strings.ToLower(results[j].GroupName)
		if groupI != groupJ {
			return groupI < groupJ
		}

		caseI := strings.ToLower(results[i].TestCaseName)
		caseJ := strings.ToLower(results[j].TestCaseName)
		if caseI != caseJ {
			return caseI < caseJ
		}

		return results[i].Id < results[j].Id
	})

	rows := []table.Row{}
	for _, test := range results {
		testPattern, err := testrun.NewTestPattern(
			test.TestFilePath.String(),
			&test.TestLineNumber,
			&test.TestCaseName,
			&test.GroupName,
		)
		if err != nil {
			continue
		}

		row := table.NewRow(table.RowData{
			columnKeyTestResult:    renderFailureType(test.AbbreviatedResult(), &m.ctx.Styles),
			columnKeyTestName:        test.GroupName + " " + test.TestCaseName,
			columnKeyMetaId:          test.Id,
			columnKeyMetaTestPattern: testPattern,
		}).WithStyle(lipgloss.NewStyle().Foreground(m.ctx.Styles.ResultsSection.TableRowTextColor))
		rows = append(rows, row)
	}

	m.resultsTable = m.resultsTable.WithRows(rows)
}

func (m Model) GetSelectedTestResultId() int {
	highlightedRow := m.resultsTable.HighlightedRow()
	if len(highlightedRow.Data) == 0 {
		return -1
	}
	id, ok := highlightedRow.Data[columnKeyMetaId]
	if !ok {
		return -1
	}
	return id.(int)
}

func (m Model) getSelectedRow() (table.Row, bool) {
	row := m.resultsTable.HighlightedRow()

	if len(row.Data) == 0 {
		return table.Row{}, false
	}
	return row, true
}

func (m Model) getSelectedTestPattern() (testrun.TestPattern, bool) {
	row, ok := m.getSelectedRow()
	if !ok {
		return testrun.TestPattern{}, false
	}
	patternData, ok := row.Data[columnKeyMetaTestPattern]
	if !ok {
		return testrun.TestPattern{}, false
	}
	pattern, ok := patternData.(testrun.TestPattern)
	if !ok {
		return testrun.TestPattern{}, false
	}
	return pattern, true
}

func (m *Model) ToggleFocus() {
	m.focus = !m.focus
}

func (m Model) Focus() bool {
	return m.focus
}

func (m Model) isBlurred() bool {
	return !m.focus
}

func (m Model) isFocused() bool {
	return m.focus
}

//
// INTERNAL FUNCTIONS
//================================================

func getColumnConfiguration(failureCauseWidth int, groupNameWidth int, testNameWidth int, styles *styles.Styles) []table.Column {
	// TODO - syling - refer to https://github.com/Evertras/bubble-table/blob/main/examples/features/main.go
	headerStyle := lipgloss.NewStyle().Foreground(styles.ResultsSection.TableHeaderTextColor)

	return []table.Column{
		table.NewColumn(columnKeyTestResult, "", failureCauseWidth).WithStyle(headerStyle.Align(lipgloss.Center)),
		// table.NewFlexColumn(columnKeyGroupName, "Group name", groupNameWidth).WithStyle(headerStyle),
		table.NewFlexColumn(columnKeyTestName, "Test name", testNameWidth).WithStyle(headerStyle),
	}
}

func renderFailureType(result string, styles *styles.Styles) string {
	switch result {
	case "P":
		return styles.PassBadge.Width(3).Render(result)
	case "S":
		return styles.SkipBadge.Width(3).Render(result)
	case "T":
		return styles.TestDefinitionErrorBadge.Width(3).Render(result)
	case "C":
		return styles.ProductionCodeErrorBadge.Width(3).Render(result)
	case "A":
		return styles.AssertionErrorBadge.Width(3).Render(result)
	default:
		return ""
	}
}

func sortPriority(result types.TestResult) int {
	switch result.Status {
	case types.StatusFail:
		switch result.FailureCause {
		case types.FailureCauseProductionCode:
			return 0
		case types.FailureCauseTestDefinition:
			return 1
		case types.FailureCauseAssertion:
			return 2
		default:
			return 3
		}
	case types.StatusPass:
		return 3
	case types.StatusSkip:
		return 4
	default:
		return 5
	}
}
