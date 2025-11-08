package resultssection

import (
	"github.com/adamakhtar/wing_commander/internal/runner"
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
		Top:          "─",
		Bottom:       "─",
		Left:         "│",
		Right:        "│",
		TopLeft:      "┌",
		TopRight:     "┐",
		BottomLeft:   "└",
		BottomRight:  "┘",
		TopJunction:    "┬",
		LeftJunction:   "├",
		RightJunction:  "┤",
		BottomJunction: "┴",
		InnerJunction: "┼",
		InnerDivider: "│",
}

const (
	columnKeyGroupName    = "group_name"
	columnKeyTestName    = "test_name"
	columnKeyFailureCause = "failure_cause"
	columnKeyMetaId = "test_id"
)

const (
	paddingX = 1
	paddingY = 0
)

type Model struct {
	ctx *context.Context
	focus bool
	resultsTable table.Model
	width int
	height int
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
		ctx: ctx,
		focus: focus,
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
	if m.isBlurred()  {
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


//
// COMMANDS
//================================================


//
// EXTERNAL FUNCTIONS
//================================================

func (m *Model) SetSize(width int, height int) {
	m.width = width
	m.height = height
	m.resultsTable = m.resultsTable.WithTargetWidth(m.width - 2 * paddingX)
	m.resultsTable = m.resultsTable.WithMinimumHeight(m.height - 2 * paddingY)

	failureCauseWidth := 3
	groupNameWidth := 15
	testNameWidth := width - failureCauseWidth - groupNameWidth

	columns := getColumnConfiguration(failureCauseWidth, groupNameWidth, testNameWidth, &m.ctx.Styles)
	m.resultsTable.WithColumns(columns)
}

func (m *Model) SetRows(testExecutionResult *runner.TestExecutionResult) {
	rows := []table.Row{}
	for _, test := range testExecutionResult.FailedTests {
		row := table.NewRow(table.RowData{
			columnKeyFailureCause: renderFailureType(test.FailureCause, &m.ctx.Styles),
			columnKeyTestName: test.GroupName + " " + test.TestCaseName,
			columnKeyMetaId: test.Id,
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
		table.NewColumn(columnKeyFailureCause, "", failureCauseWidth).WithStyle(headerStyle.Align(lipgloss.Center)),
		// table.NewFlexColumn(columnKeyGroupName, "Group name", groupNameWidth).WithStyle(headerStyle),
		table.NewFlexColumn(columnKeyTestName, "Test name", testNameWidth).WithStyle(headerStyle),
	}
}

func renderFailureType(failureCause types.FailureCause, styles *styles.Styles) string {
	switch failureCause {
	case types.FailureCauseTestDefinition:
		return styles.TestDefinitionErrorBadge.Width(3).Render(failureCause.Abbreviated())
	case types.FailureCauseProductionCode:
		return styles.ProductionCodeErrorBadge.Width(3).Render(failureCause.Abbreviated())
	case types.FailureCauseAssertion:
		return styles.AssertionErrorBadge.Width(3).Render(failureCause.Abbreviated())
	default:
		return ""
	}
}
