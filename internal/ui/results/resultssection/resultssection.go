package resultssection

import (
	"github.com/adamakhtar/wing_commander/internal/runner"
	"github.com/adamakhtar/wing_commander/internal/ui/context"
	"github.com/adamakhtar/wing_commander/internal/ui/keys"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyGroupName    = "group_name"
	columnKeyTestName    = "test_name"
	columnKeyFailureCause = "failure_cause"
	columnKeyMetaId = "test_id"
)

type Model struct {
	ctx *context.Context
	resultsTable table.Model
}

//
// BUILDERS
//================================================

func NewModel(ctx *context.Context) Model {
	columns := getColumnConfiguration(3, 15, 15)

	resultsTable := table.New(columns).Focused(true).WithBaseStyle(
		lipgloss.NewStyle().Align(lipgloss.Left),
	)

	return Model{
		ctx: ctx,
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
	var panelStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).Padding(1, 1)

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
	m.resultsTable.WithTargetWidth(width)
	m.resultsTable.WithMinimumHeight(height)

	failureCauseWidth := 3
	groupNameWidth := 15
	testNameWidth := width - failureCauseWidth - groupNameWidth

	columns := getColumnConfiguration(failureCauseWidth, groupNameWidth, testNameWidth)
	m.resultsTable.WithColumns(columns)
}

func (m *Model) SetRows(testExecutionResult *runner.TestExecutionResult) {
	rows := []table.Row{}
	for _, test := range testExecutionResult.FailedTests {
		// Combine group name and test case name for display
		// testName := test.TestCaseName
		// if test.GroupName != "" {
		// 	if testName != "" {
		// 		testName = fmt.Sprintf("%s#%s", test.GroupName, test.TestCaseName)
		// 	} else {
		// 		testName = test.GroupName
		// 	}
		// }
		row := table.NewRow(table.RowData{
			columnKeyFailureCause: test.FailureCause.Abbreviated(),
			columnKeyGroupName: test.GroupName,
			columnKeyTestName: test.TestCaseName,
			columnKeyMetaId: test.Id,
		})
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

//
// INTERNAL FUNCTIONS
//================================================

func getColumnConfiguration(failureCauseWidth int, groupNameWidth int, testNameWidth int) []table.Column {
	// TODO - syling - refer to https://github.com/Evertras/bubble-table/blob/main/examples/features/main.go
	return []table.Column{
		table.NewColumn(columnKeyFailureCause, "", failureCauseWidth).WithStyle(lipgloss.NewStyle().Align(lipgloss.Center)),
		table.NewFlexColumn(columnKeyGroupName, "Group name", groupNameWidth),
		table.NewFlexColumn(columnKeyTestName, "Test name", testNameWidth),
	}
}
