package resultssection

import (
	"github.com/adamakhtar/wing_commander/internal/runner"
	"github.com/adamakhtar/wing_commander/internal/ui/context"
	"github.com/adamakhtar/wing_commander/internal/ui/keys"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	ctx context.Context
	resultsTable table.Model
}

//
// BUILDERS
//================================================

func NewModel() Model {
	columns := getColumnConfiguration(0, 0)

	resultsTable := table.New(table.WithColumns(columns))
	resultsTable.Focus()

	return Model{
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
	// Make your own border
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

func (m *Model) UpdateContext(ctx context.Context) {
	m.ctx = ctx
}

func (m *Model) SetSize(width int, height int) {
	m.resultsTable.SetWidth(width)
	m.resultsTable.SetHeight(height)

	failureCauseWidth := 3
	testNameWidth := width - failureCauseWidth

	columns := getColumnConfiguration(failureCauseWidth, testNameWidth)
	m.resultsTable.SetColumns(columns)
}

func (m *Model) SetRows(testExecutionResult *runner.TestExecutionResult) {
	rows := []table.Row{}
	for _, test := range testExecutionResult.FailedTests {
		rows = append(rows, table.Row{
			test.FailureCause.Abbreviated(),
			test.Name,
		})
	}
	m.resultsTable.SetRows(rows)
}

//
// INTERNAL FUNCTIONS
//================================================

func getColumnConfiguration(failureCauseWidth int, testNameWidth int) []table.Column {
	return []table.Column{
	{Title: "", Width: failureCauseWidth},
	{Title: "Test", Width: testNameWidth},
	}
}
