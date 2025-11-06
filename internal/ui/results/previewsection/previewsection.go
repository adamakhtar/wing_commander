package previewsection

import (
	"fmt"

	"github.com/adamakhtar/wing_commander/internal/filesnippet"
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

	content := ""

	content = lipgloss.JoinVertical(
		lipgloss.Top,
		content,
		m.renderTestHeading(innerWidth),
	)

	content = lipgloss.JoinVertical(
		lipgloss.Top,
		content,
		m.renderFailureMessage(innerWidth),
	)

	for _, frame := range m.testResult.FilteredBacktrace {
		backtraceLineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
		content = lipgloss.JoinVertical(
			lipgloss.Top,
			content,
			backtraceLineStyle.Render(frame.RelativeFilePath(m.ctx.Config.ProjectPath) + ":" + fmt.Sprintf("%d", frame.Line)))

		snippet, err := filesnippet.ExtractLines(frame.File, frame.Line, 5)
		if err != nil {
			log.Error("failed to extract lines", "error", err)
			continue
		}

		content = lipgloss.JoinVertical(
			lipgloss.Top,
			content,
			m.renderFileSnippet(snippet, innerWidth))
	}

	return renderPanel(content, paddingX, paddingY)
}

func (m Model) renderTestHeading(innerWidth int) string {
	headingStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("62")).
		Margin(0, 0 , 1, 0)

	testName := m.testResult.GroupName + " " + m.testResult.TestCaseName
	return  headingStyle.Width(innerWidth).Render(testName)
}


func (m Model) renderFailureMessage(innerWidth int) string {
	alertStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("198")).
		Foreground(lipgloss.Color("255")).
		Align(lipgloss.Left).
		Width(innerWidth).
		Padding(1, 1).
		Margin(0, 0 , 1, 0)

		switch {
		case m.testResult.ErrorMessage != "":
				return alertStyle.Render(m.testResult.ErrorMessage)
		case m.testResult.FailedAssertionMessage != "":
				return alertStyle.Render(m.testResult.FailedAssertionMessage)
		default:
			return ""
		}
}

func (m Model) renderAssertionFailure(assertionMessage string, innerWidth int) string {
	if m.testResult.FailedAssertionMessage == "" {
		return ""
	}

	alertStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("198")).
		Foreground(lipgloss.Color("255")).
		Align(lipgloss.Left).
		Width(innerWidth).
		Padding(1, 1).
		Margin(0, 0 , 1, 0)

	return alertStyle.Render(m.testResult.FailedAssertionMessage)
}

func (m Model) renderFileSnippet(snippet *filesnippet.FileSnippet, innerWidth int) string {
	fileSnippetStyle := lipgloss.NewStyle().Margin(0, 0 , 1, 0)
	fileLineStyle := lipgloss.NewStyle().Background(lipgloss.Color("0")).Foreground(lipgloss.Color("15"))
	highlightedFileLineStyle := lipgloss.NewStyle().Background(lipgloss.Color("198")).Foreground(lipgloss.Color("255"))

	content := ""
	for _, line := range snippet.Lines {
		lineStyle := fileLineStyle

		if line.IsCenter {
			lineStyle = highlightedFileLineStyle
		}

		content = lipgloss.JoinVertical(
			lipgloss.Top,
			content,
			lineStyle.Width(innerWidth).Render(fmt.Sprintf("%d: %s", line.Number, line.Content)))
	}

	return fileSnippetStyle.Render(content)
}

func renderPanel(content string, paddingX int, paddingY int) string {
	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(paddingY, paddingX)

	return panelStyle.Render(content)
}

func (m *Model) SetSize(width int, height int) {
	m.width = width
	m.height = height
}

func (m *Model) SetTestResult(testResult *types.TestResult) {
	m.testResult = testResult
}
