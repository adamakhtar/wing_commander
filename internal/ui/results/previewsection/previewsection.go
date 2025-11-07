package previewsection

import (
	"fmt"

	"github.com/adamakhtar/wing_commander/internal/filesnippet"
	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/adamakhtar/wing_commander/internal/ui/context"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

const (
	paddingX = 1
	paddingY = 0
)

type Model struct {
	ctx *context.Context
	focus bool
	width int
	height int
	testResult *types.TestResult
	viewport viewport.Model
}

func NewModel(ctx *context.Context, focus bool) Model {
	return Model{
		ctx: ctx,
		focus: focus,
		width: 0,
		height: 0,
		testResult: nil,
		viewport: viewport.New(0, 0),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.isBlurred()  {
		return m, nil
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}


func (m Model) innerDimensions(width, height int) (innerWidth, innerHeight int) {
	innerWidth = width - (2 * paddingX)
	innerHeight = height - 2
	return innerWidth, innerHeight
}

func (m Model) View() string {
	if m.testResult == nil {
		return m.renderPanel("No Test Result Selected\n")
	}

	content := m.viewport.View()
	return m.renderPanel(content)
}

func (m Model) buildContent(innerWidth int) string {
	if m.testResult == nil {
		return "No Test Result Selected\n"
	}

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

	return content
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

func (m Model)renderPanel(content string) string {
	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(paddingY, paddingX)

	if m.isFocused() {
		panelStyle = panelStyle.BorderForeground(lipgloss.Color("5"))
	}

	return panelStyle.Render(content)
}

func (m *Model) SetSize(width int, height int) {
	m.width = width
	m.height = height

	innerWidth, innerHeight := m.innerDimensions(width, height)
	m.viewport.Width = innerWidth
	m.viewport.Height = innerHeight
	m.viewport.SetContent(m.buildContent(innerWidth))
}

func (m *Model) SetTestResult(testResult *types.TestResult) {
	m.testResult = testResult

	innerWidth, _ := m.innerDimensions(m.width, m.height)
	m.viewport.SetContent(m.buildContent(innerWidth))
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