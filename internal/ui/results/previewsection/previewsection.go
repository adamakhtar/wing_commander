package previewsection

import (
	"fmt"
	"strings"

	"github.com/adamakhtar/wing_commander/internal/filesnippet"
	"github.com/adamakhtar/wing_commander/internal/projectfs"
	"github.com/adamakhtar/wing_commander/internal/testresult"
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
	ctx        *context.Context
	focus      bool
	width      int
	height     int
	testResult *testresult.TestResult
	viewport   viewport.Model
}

func NewModel(ctx *context.Context, focus bool) Model {
	return Model{
		ctx:        ctx,
		focus:      focus,
		width:      0,
		height:     0,
		testResult: nil,
		viewport:   viewport.New(0, 0),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.isBlurred() {
		return m, nil
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) innerDimensions(width, height int) (innerWidth, innerHeight int) {
	innerWidth = width - (2 * paddingX)
	innerHeight = height - (2 * paddingY)
	return innerWidth, innerHeight
}

func (m Model) View() string {
	content := m.viewport.View()
	return m.renderPanel(content)
}

func (m Model) buildContent(innerWidth int) string {
	if m.testResult == nil {
		return lipgloss.PlaceHorizontal(innerWidth, lipgloss.Center, m.ctx.Styles.BodyTextLight.Render("No Test Result Selected"))
	}

	sb := strings.Builder{}

	sb.WriteString(m.renderTestHeading(innerWidth))
	sb.WriteString("\n")
	sb.WriteString(m.renderTestResult(innerWidth))
	sb.WriteString("\n")
	sb.WriteString(m.renderFailureMessage(innerWidth))
	sb.WriteString("\n")

	for _, frame := range m.testResult.FilteredBacktrace.Frames {
		fs := projectfs.GetProjectFS()
		relPath, err := fs.Rel(frame.FilePath)
		var line string
		if err != nil {
			// Fallback to absolute path
			line = frame.FilePath.String() + ":" + fmt.Sprintf("%d", frame.Line)
		} else {
			line = relPath.String() + ":" + fmt.Sprintf("%d", frame.Line)
		}
		sb.WriteString(m.ctx.Styles.PreviewSection.BacktracePath.Width(innerWidth).Render(line))

		snippet, err := filesnippet.ExtractLines(frame.FilePath.String(), frame.Line, 5)
		if err != nil {
			log.Error("failed to extract lines", "error", err)
			continue
		}

		sb.WriteString(m.renderFileSnippet(snippet, innerWidth))
		sb.WriteString("\n")
	}

	// log.Debug("sb", "sb", sb.String())

	return sb.String()
}

func (m Model) renderTestHeading(innerWidth int) string {
	testName := m.testResult.GroupName + " " + m.testResult.TestCaseName
	testPath := m.testResult.TestFilePath.String() + ":" + fmt.Sprintf("%d", m.testResult.TestLineNumber)

	return lipgloss.JoinVertical(lipgloss.Top,
		m.ctx.Styles.HeadingTextStyle.Width(innerWidth).Render(testName),
		m.ctx.Styles.PreviewSection.BacktracePath.Width(innerWidth).Margin(0, 0, 1).Render(testPath),
	)
}

func (m Model) renderTestResult(innerWidth int) string {
	switch {
	case m.testResult.IsFailed():
		switch m.testResult.FailureCause {
		case testresult.FailureCauseTestDefinition:
			return m.ctx.Styles.TestDefinitionErrorBadge.Width(innerWidth).Render(m.testResult.FailureCause.String())
		case testresult.FailureCauseProductionCode:
			return m.ctx.Styles.ProductionCodeErrorBadge.Width(innerWidth).Render(m.testResult.FailureCause.String())
		case testresult.FailureCauseAssertion:
			return m.ctx.Styles.AssertionErrorBadge.Width(innerWidth).Render(m.testResult.FailureCause.String())
		default:
			return ""
		}
	case m.testResult.IsSkipped():
		return m.ctx.Styles.SkipBadge.Width(innerWidth).Render(string(m.testResult.Status))
	case m.testResult.IsPassed():
		return m.ctx.Styles.PassBadge.Width(innerWidth).Render(string(m.testResult.Status))
	default:
		return ""
	}
}

func (m Model) renderFailureMessage(innerWidth int) string {
	alertStyle := m.ctx.Styles.Preview.AlertStyle
	alertStyle = alertStyle.
		Align(lipgloss.Left).
		Width(innerWidth).
		Padding(1, 4).
		Margin(0, 0, 1, 0)

	if m.testResult.FailureDetails != "" {
		return alertStyle.Render(m.testResult.FailureDetails)
	}
	return ""
}

func (m Model) renderFileSnippet(snippet *filesnippet.FileSnippet, innerWidth int) string {
	content := ""
	for _, line := range snippet.Lines {
		lineStyle := m.ctx.Styles.PreviewSection.CodeLine

		if line.IsCenter {
			lineStyle = m.ctx.Styles.PreviewSection.HighlightedCodeLine
		}

		content = lipgloss.JoinVertical(
			lipgloss.Top,
			content,
			lineStyle.Width(innerWidth).Render(fmt.Sprintf("%d: %s", line.Number, line.Content)))
	}

	return lipgloss.NewStyle().Margin(0, 0, 1, 0).Render(content)
}

func (m Model) renderPanel(content string) string {
	panelStyle := m.ctx.Styles.Border.Padding(paddingY, paddingX)

	if m.isFocused() {
		panelStyle = panelStyle.Inherit(m.ctx.Styles.BorderActive)
	} else {
		panelStyle = panelStyle.Inherit(m.ctx.Styles.BorderMuted)
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

func (m *Model) SetTestResult(testResult *testresult.TestResult) {
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
