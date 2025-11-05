package previewsection

import (
	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)



type Model struct {
	width int
	height int
	testResult *types.TestResult
}

func NewModel() Model {
	return Model{
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

	headingStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Margin(0, 0 , 0, 0)
	alertStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("198")).
		Foreground(lipgloss.Color("255")).
		Align(lipgloss.Center).
		Width(innerWidth).
		Padding(1, 1).
		Margin(0, 0 , 1, 0)

	// stackFrameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	// codePreviewStyle := lipgloss.NewStyle().Background(lipgloss.Color("0")).Foreground(lipgloss.Color("15"))

	content := headingStyle.Render("NAME NAME")

	log.Debug("testResult", "testResult", m.testResult)

	if m.testResult.ErrorMessage != "" {
		// errorMsg := "lorem ipsum dolor sit amet,\nconsectetur\nadipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
		log.Debug("errorMsg", "errorMsg", m.testResult.ErrorMessage)
		errorMsg := m.testResult.ErrorMessage
		alert := alertStyle.Render(errorMsg)
		// centeredAlert := lipgloss.PlaceHorizontal(innerWidth, lipgloss.Center, m.testResult.ErrorMessage)
		// centeredAlert := lipgloss.PlaceHorizontal(innerWidth, lipgloss.Center, errorMsg)
		content = lipgloss.JoinVertical(
			lipgloss.Top,
			content,
			alert)
	}

	if m.testResult.FailedAssertionMessage != "" {
		// errorMsg := "lorem ipsum dolor sit amet,\nconsectetur\nadipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
		log.Debug("failedAssertionMessage", "failedAssertionMessage", m.testResult.FailedAssertionMessage)
		failedAssertionMessage := m.testResult.FailedAssertionMessage
		alert := alertStyle.Render(failedAssertionMessage)
		content = lipgloss.JoinVertical(
			lipgloss.Top,
			content,
			alert,
			alert)
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
