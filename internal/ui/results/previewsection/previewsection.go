package previewsection

import (
	"fmt"

	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/adamakhtar/wing_commander/internal/ui/context"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)



type Model struct {
	ctx *context.Context
	width int
	height int
	testResult *types.TestResult
}

func NewModel(ctx *context.Context) Model {
	return Model{
		ctx: ctx,
		width: 0,
		height: 0,
		testResult: nil,
	}
}

func (m Model) View() string {
	if m.testResult == nil {
		return "No Test Result Selected\n"
	}

	var paddingX = 1
	var paddingY = 0
	var innerWidth = m.width - (2 * paddingX)

	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(paddingY, paddingX)

	headingStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Margin(0, 0 , 1, 0)
	alertStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("198")).
		Foreground(lipgloss.Color("255")).
		Align(lipgloss.Left).
		Width(innerWidth).
		Padding(1, 1).
		Margin(0, 0 , 1, 0)

	// stackFrameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	// codePreviewStyle := lipgloss.NewStyle().Background(lipgloss.Color("0")).Foreground(lipgloss.Color("15"))

	content := headingStyle.Render(m.testResult.GroupName + " " + m.testResult.TestCaseName)

	log.Debug("testResult", "FailedAssertionMessage", m.testResult.FailedAssertionMessage, "ErrorMessage", m.testResult.ErrorMessage)

	switch {
	case m.testResult.ErrorMessage != "":
			errorAlert := alertStyle.Render(m.testResult.ErrorMessage)

			content = lipgloss.JoinVertical(
				lipgloss.Top,
				content,
				errorAlert,
			)
	case m.testResult.FailedAssertionMessage != "":
			assertionAlert := alertStyle.Render(m.testResult.FailedAssertionMessage)
			content = lipgloss.JoinVertical(
				lipgloss.Top,
				content,
				assertionAlert)
	}

	for _, frame := range m.testResult.FilteredBacktrace {
		stackFrameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
		content = lipgloss.JoinVertical(
			lipgloss.Top,
			content,
			stackFrameStyle.Render(frame.RelativeFilePath(m.ctx.Config.ProjectPath) + ":" + fmt.Sprintf("%d", frame.Line)))
	}

	return panelStyle.Render(content)
}

func (m *Model) SetSize(width int, height int) {
	m.width = width
	m.height = height
}

func (m *Model) SetTestResult(testResult *types.TestResult) {
	m.testResult = testResult
}
