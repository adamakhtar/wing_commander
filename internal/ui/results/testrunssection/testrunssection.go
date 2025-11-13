package testrunssection

import (
	"fmt"
	"strings"

	"github.com/adamakhtar/wing_commander/internal/testrun"
	"github.com/adamakhtar/wing_commander/internal/ui/context"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	paddingX = 1
)

type Model struct {
	ctx      *context.Context
	testRuns *testrun.TestRuns
	focus    bool
	width    int
	height   int
}

func NewModel(ctx *context.Context, testRuns *testrun.TestRuns) Model {
	return Model{
		ctx:      ctx,
		testRuns: testRuns,
		focus:    false,
		width:    0,
		height:   0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	innerWidth := m.width - 2*paddingX

	sb := strings.Builder{}
	sb.WriteString(m.ctx.Styles.HeadingTextStyle.Width(innerWidth).Render("Recent Test Runs"))
	sb.WriteString("\n")

	for _, testRun := range m.testRuns.AllRecentFirst() {
		sb.WriteString(m.ctx.Styles.TestRunsSection.Label.Width(innerWidth).Render(Label(testRun)))
	}

	return m.ctx.Styles.Border.Padding(0, paddingX).Inherit(m.ctx.Styles.BorderMuted).Width(m.width).Height(m.height).Render(sb.String())
}

func (m *Model) SetSize(width int, height int) {
	m.width = width
	m.height = height
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

func Label(t testrun.TestRun) string {
	switch testrun.Mode(t.Mode) {
	case testrun.ModeRunWholeSuite:
		return "Run whole suite"
	case testrun.ModeRunSelectedPatterns:
		return formatSelectedPatternsLabel(len(t.Patterns))
	case testrun.ModeReRunSingleFailure:
		return "Re-run failure"
	case testrun.ModeReRunAllFailures:
		return "Re-run all failed"
	default:
		return formatSelectedPatternsLabel(len(t.Patterns))
	}
}

func formatSelectedPatternsLabel(count int) string {
	switch count {
	case 0:
		return "Run 0 patterns"
	case 1:
		return "Run 1 pattern"
	default:
		return fmt.Sprintf("Run %d patterns", count)
	}
}
