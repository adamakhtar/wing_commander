package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/adamakhtar/wing_commander/internal/editor"
	"github.com/adamakhtar/wing_commander/internal/runner"
	"github.com/adamakhtar/wing_commander/internal/types"
)

// Model represents the UI state
type Model struct {
	// Data
	failureGroups []types.FailureGroup
	testResults   []types.TestResult

	// Selection state
	selectedGroup int
	selectedTest  int
	selectedFrame int
	activePane    int // 0=groups, 1=tests, 2=backtrace

	// UI state
	width         int
	height        int
	showFullFrames bool

	// Execution result
	result *runner.TestExecutionResult

	// Services
	editor *editor.Editor
	runner *runner.TestRunner
}

// NewModel creates a new UI model
func NewModel(result *runner.TestExecutionResult, testRunner *runner.TestRunner) Model {
	return Model{
		failureGroups:   result.FailureGroups,
		testResults:     result.TestResults,
		selectedGroup:   0,
		selectedTest:    0,
		selectedFrame:   0,
		activePane:      0,
		showFullFrames:  false,
		result:          result,
		editor:          editor.NewEditor(),
		runner:          testRunner,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "tab":
			m.activePane = (m.activePane + 1) % 3
			return m, nil

		case "shift+tab":
			m.activePane = (m.activePane + 2) % 3
			return m, nil

		case "up", "k":
			return m.handleUpKey(), nil

		case "down", "j":
			return m.handleDownKey(), nil

		case "f":
			m.showFullFrames = !m.showFullFrames
			return m, nil

		case "o":
			return m, m.handleOpenFile()

		case "r":
			return m, m.handleReRunTests()
		}

	case OpenFileSuccessMsg:
		// File opened successfully - no UI update needed
		return m, nil

	case OpenFileErrorMsg:
		// File opening failed - could show error message in future
		return m, nil

	case ReRunSuccessMsg:
		// Tests re-run successfully - update the model with new results
		m.result = msg.Result
		m.failureGroups = msg.Result.FailureGroups
		m.testResults = msg.Result.TestResults
		// Reset selections to avoid out-of-bounds
		if m.selectedGroup >= len(m.failureGroups) {
			m.selectedGroup = 0
		}
		if len(m.failureGroups) > 0 && m.selectedGroup < len(m.failureGroups) {
			if m.selectedTest >= len(m.failureGroups[m.selectedGroup].Tests) {
				m.selectedTest = 0
			}
		}
		return m, nil

	case ReRunErrorMsg:
		// Test re-run failed - could show error message in future
		return m, nil
	}

	return m, nil
}

// handleUpKey handles up arrow navigation
func (m Model) handleUpKey() Model {
	switch m.activePane {
	case 0: // Groups pane
		if m.selectedGroup > 0 {
			m.selectedGroup--
		}
	case 1: // Tests pane
		if len(m.failureGroups) > 0 && m.selectedGroup < len(m.failureGroups) {
			if m.selectedTest > 0 {
				m.selectedTest--
			}
		}
	case 2: // Backtrace pane
		if len(m.failureGroups) > 0 && m.selectedGroup < len(m.failureGroups) {
			group := m.failureGroups[m.selectedGroup]
			if m.selectedTest < len(group.Tests) {
				test := group.Tests[m.selectedTest]
				_ = m.getCurrentTestFrames(test) // Get frames to validate bounds
				if m.selectedFrame > 0 {
					m.selectedFrame--
				}
			}
		}
	}
	return m
}

// handleDownKey handles down arrow navigation
func (m Model) handleDownKey() Model {
	switch m.activePane {
	case 0: // Groups pane
		if m.selectedGroup < len(m.failureGroups)-1 {
			m.selectedGroup++
		}
	case 1: // Tests pane
		if len(m.failureGroups) > 0 && m.selectedGroup < len(m.failureGroups) {
			group := m.failureGroups[m.selectedGroup]
			if m.selectedTest < len(group.Tests)-1 {
				m.selectedTest++
			}
		}
	case 2: // Backtrace pane
		if len(m.failureGroups) > 0 && m.selectedGroup < len(m.failureGroups) {
			group := m.failureGroups[m.selectedGroup]
			if m.selectedTest < len(group.Tests) {
				test := group.Tests[m.selectedTest]
				frames := m.getCurrentTestFrames(test)
				if m.selectedFrame < len(frames)-1 {
					m.selectedFrame++
				}
			}
		}
	}
	return m
}

// handleOpenFile handles opening the selected file in an external editor
func (m Model) handleOpenFile() tea.Cmd {
	return func() tea.Msg {
		if len(m.failureGroups) == 0 || m.selectedGroup >= len(m.failureGroups) {
			return OpenFileErrorMsg{Error: "no group selected"}
		}

		group := m.failureGroups[m.selectedGroup]
		if m.selectedTest >= len(group.Tests) {
			return OpenFileErrorMsg{Error: "no test selected"}
		}

		test := group.Tests[m.selectedTest]
		frames := m.getCurrentTestFrames(test)

		if len(frames) == 0 || m.selectedFrame >= len(frames) {
			return OpenFileErrorMsg{Error: "no frame selected"}
		}

		frame := frames[m.selectedFrame]
		err := m.editor.OpenFile(frame.File, frame.Line)
		if err != nil {
			return OpenFileErrorMsg{Error: err.Error()}
		}

		return OpenFileSuccessMsg{File: frame.File, Line: frame.Line}
	}
}

// handleReRunTests handles re-running tests for the selected group
func (m Model) handleReRunTests() tea.Cmd {
	return func() tea.Msg {
		if len(m.failureGroups) == 0 || m.selectedGroup >= len(m.failureGroups) {
			return ReRunErrorMsg{Error: "no group selected"}
		}

		// For now, we'll re-run all tests
		// In a future enhancement, we could run only specific tests from the group
		result, err := m.runner.ExecuteTests()
		if err != nil {
			return ReRunErrorMsg{Error: err.Error()}
		}

		return ReRunSuccessMsg{Result: result}
	}
}

// getCurrentFrames returns the appropriate frames based on showFullFrames setting
func (m Model) getCurrentFrames(group types.FailureGroup) []types.StackFrame {
	if m.showFullFrames {
		return group.NormalizedBacktrace
	}
	// For now, we'll use the first test's filtered backtrace
	// In a more advanced implementation, we might want to show a combined view
	if len(group.Tests) > 0 {
		return group.Tests[0].FilteredBacktrace
	}
	return group.NormalizedBacktrace
}

// getCurrentTestFrames returns the appropriate frames for the selected test
func (m Model) getCurrentTestFrames(test types.TestResult) []types.StackFrame {
	if m.showFullFrames {
		return test.FullBacktrace
	}
	return test.FilteredBacktrace
}

// getTailFrames returns all frames except the bottom frame (first frame in backtrace)
func (m Model) getTailFrames(test types.TestResult) []types.StackFrame {
	frames := m.getCurrentTestFrames(test)
	if len(frames) <= 1 {
		return []types.StackFrame{}
	}
	// Return all frames except the first one (bottom frame)
	return frames[1:]
}

// formatTailFrames formats tail frames as a chain with arrows
func (m Model) formatTailFrames(frames []types.StackFrame, maxWidth int) string {
	if len(frames) == 0 {
		return "(no additional frames)"
	}

	var parts []string
	for _, frame := range frames {
		parts = append(parts, fmt.Sprintf("%s:%d", frame.File, frame.Line))
	}

	result := strings.Join(parts, " → ")
	if len(result) > maxWidth {
		// Truncate from the end, keeping the beginning
		result = result[:maxWidth-3] + "..."
	}
	return result
}

// getFailureCauseIcon returns an icon for the failure cause
func getFailureCauseIcon(cause types.FailureCause) string {
	switch cause {
	case types.FailureCauseTestDefinition:
		return "🏗️" // Test definition error
	case types.FailureCauseProductionCode:
		return "🚀" // Production code error
	case types.FailureCauseAssertion:
		return "❌" // Assertion failure
	default:
		return "❓" // Unknown
	}
}

// View renders the UI
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Calculate pane dimensions
	paneWidth := (m.width - 4) / 3 // Account for borders and spacing

	// Render each pane
	groupsPane := m.renderGroupsPane(paneWidth, m.height-2)
	testsPane := m.renderTestsPane(paneWidth, m.height-2)
	backtracePane := m.renderBacktracePane(paneWidth, m.height-2)

	// Combine panes horizontally
	panes := lipgloss.JoinHorizontal(lipgloss.Top, groupsPane, testsPane, backtracePane)

	// Add status bar
	statusBar := m.renderStatusBar()

	// Combine everything
	return lipgloss.JoinVertical(lipgloss.Left, panes, statusBar)
}

// renderGroupsPane renders the groups pane
func (m Model) renderGroupsPane(width, height int) string {
	title := GetPaneTitleStyle().Render("Failure Groups")
	isActive := m.activePane == 0

	// Create content
	var content strings.Builder

	if len(m.failureGroups) == 0 {
		content.WriteString(GetSuccessTextStyle().Render("✅ All tests passed!"))
	} else {
		// Group failures by cause
		groupsByCause := make(map[types.FailureCause][]types.FailureGroup)
		for _, group := range m.failureGroups {
			if len(group.Tests) > 0 {
				cause := group.Tests[0].FailureCause
				groupsByCause[cause] = append(groupsByCause[cause], group)
			}
		}

		// Define order for displaying causes
		causeOrder := []types.FailureCause{
			types.FailureCauseProductionCode,
			types.FailureCauseTestDefinition,
			types.FailureCauseAssertion,
		}

		causeTitles := map[types.FailureCause]string{
			types.FailureCauseProductionCode:  "Prod Error",
			types.FailureCauseTestDefinition:  "Test Error",
			types.FailureCauseAssertion:       "Failed Assertion",
		}

		currentGlobalIndex := 0
		for _, cause := range causeOrder {
			if groups, exists := groupsByCause[cause]; exists && len(groups) > 0 {
				// Add section header with better styling
				title := causeTitles[cause]
				divider := strings.Repeat("─", width-6-len(title))
				header := fmt.Sprintf("── %s %s", title, divider)
				content.WriteString(GetDimmedTextStyle().Render(header))
				content.WriteString("\n")

				// Add groups for this cause
				for _, group := range groups {
					// Get failure cause icon and count (in yellow)
					var causeIcon string
					var countText string
					if len(group.Tests) > 0 {
						causeIcon = getFailureCauseIcon(group.Tests[0].FailureCause)
						countText = GetYellowTextStyle().Render(fmt.Sprintf("%d", group.Count))
					} else {
						causeIcon = "❓"
						countText = GetYellowTextStyle().Render("0")
					}

					// Show error message (in yellow, truncated if necessary)
					errorMsg := GetYellowTextStyle().Render(group.ErrorMessage)
					if len(group.ErrorMessage) > width-8 { // Account for icon and count
						truncated := group.ErrorMessage[:width-11] + "..."
						errorMsg = GetYellowTextStyle().Render(truncated)
					}

					// First line: icon count - error message
					firstLine := fmt.Sprintf("%s %s - %s", causeIcon, countText, errorMsg)

					// Second line: bottom frame (first frame in backtrace)
					location := "Unknown"
					if len(group.NormalizedBacktrace) > 0 {
						frame := group.NormalizedBacktrace[0] // Bottom frame (first frame)
						// Convert absolute path to relative if it contains the project path
						filePath := frame.File
						if projectPath, err := m.runner.GetWorkingDirectory(); err == nil {
							if strings.HasPrefix(filePath, projectPath+"/") {
								filePath = strings.TrimPrefix(filePath, projectPath+"/")
							}
						}
						location = fmt.Sprintf("%s:%d", filePath, frame.Line)
					}

					secondLine := fmt.Sprintf("  %s", location)

					// Apply selection highlighting if this group is selected
					if currentGlobalIndex == m.selectedGroup && isActive {
						firstLine = GetSelectedTextStyle().Render(firstLine)
						secondLine = GetSelectedTextStyle().Render(secondLine)
					}

					content.WriteString(firstLine)
					content.WriteString("\n")
					content.WriteString(secondLine)
					content.WriteString("\n")
					currentGlobalIndex++
				}
				content.WriteString("\n") // Add spacing between sections
			}
		}
	}

	// Apply pane styling
	paneStyle := GetPaneStyle(isActive).Width(width).Height(height)
	return paneStyle.Render(lipgloss.JoinVertical(lipgloss.Left, title, content.String()))
}


// renderTestsPane renders the tests pane
func (m Model) renderTestsPane(width, height int) string {
	title := GetPaneTitleStyle().Render("Tests")
	isActive := m.activePane == 1

	var content strings.Builder

	if len(m.failureGroups) == 0 || m.selectedGroup >= len(m.failureGroups) {
		content.WriteString(GetDimmedTextStyle().Render("No groups selected"))
	} else {
		group := m.failureGroups[m.selectedGroup]
		for i, test := range group.Tests {
			style := GetNormalTextStyle()
			if i == m.selectedTest && isActive {
				style = GetSelectedTextStyle()
			}

			// Show test name on first line
			name := test.Name
			if len(name) > width-4 {
				name = name[:width-7] + "..."
			}
			content.WriteString(style.Render(name))
			content.WriteString("\n")

			// Show tail frames (all frames except the bottom frame) on second line
			tailFrames := m.getTailFrames(test)
			if len(tailFrames) == 0 {
				content.WriteString(style.Render("  (no additional frames)"))
			} else {
				tailStr := m.formatTailFrames(tailFrames, width-4)
				content.WriteString(style.Render("  " + tailStr))
			}
			content.WriteString("\n")
		}
	}

	paneStyle := GetPaneStyle(isActive).Width(width).Height(height)
	return paneStyle.Render(lipgloss.JoinVertical(lipgloss.Left, title, content.String()))
}

// renderBacktracePane renders the backtrace pane
func (m Model) renderBacktracePane(width, height int) string {
	title := GetPaneTitleStyle().Render("Backtrace")
	isActive := m.activePane == 2

	var content strings.Builder

	if len(m.failureGroups) == 0 || m.selectedGroup >= len(m.failureGroups) {
		content.WriteString(GetDimmedTextStyle().Render("No groups selected"))
	} else {
		group := m.failureGroups[m.selectedGroup]
		if m.selectedTest >= len(group.Tests) {
			content.WriteString(GetDimmedTextStyle().Render("No test selected"))
		} else {
			test := group.Tests[m.selectedTest]
			frames := m.getCurrentTestFrames(test)

			for i, frame := range frames {
				style := GetChangeIntensityStyle(frame.ChangeIntensity)
				if i == m.selectedFrame && isActive {
					style = GetSelectedTextStyle()
				}

				// Format frame display
				line := fmt.Sprintf("%s:%d", frame.File, frame.Line)
				if frame.Function != "" {
					line += fmt.Sprintf(" in %s", frame.Function)
				}

				// Add change intensity indicator
				if frame.ChangeIntensity > 0 {
					line += fmt.Sprintf(" [%d]", frame.ChangeIntensity)
				}

				// Truncate if too long
				if len(line) > width-4 {
					line = line[:width-7] + "..."
				}

				content.WriteString(style.Render(line))
				content.WriteString("\n")
			}
		}
	}

	paneStyle := GetPaneStyle(isActive).Width(width).Height(height)
	return paneStyle.Render(lipgloss.JoinVertical(lipgloss.Left, title, content.String()))
}

// renderStatusBar renders the status bar with keybindings
func (m Model) renderStatusBar() string {
	keyBindings := []string{
		GetKeyBindingStyle().Render("↑↓") + " navigate",
		GetKeyBindingStyle().Render("Tab") + " switch panes",
		GetKeyBindingStyle().Render("f") + " toggle frames",
		GetKeyBindingStyle().Render("o") + " open file",
		GetKeyBindingStyle().Render("r") + " re-run tests",
		GetKeyBindingStyle().Render("q") + " quit",
	}

	statusText := strings.Join(keyBindings, " • ")
	return GetStatusBarStyle().Width(m.width).Render(statusText)
}

// Message types for async operations

// OpenFileSuccessMsg is sent when a file is successfully opened
type OpenFileSuccessMsg struct {
	File string
	Line int
}

// OpenFileErrorMsg is sent when opening a file fails
type OpenFileErrorMsg struct {
	Error string
}

// ReRunSuccessMsg is sent when tests are successfully re-run
type ReRunSuccessMsg struct {
	Result *runner.TestExecutionResult
}

// ReRunErrorMsg is sent when re-running tests fails
type ReRunErrorMsg struct {
	Error string
}
