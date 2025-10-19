package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
}

// NewModel creates a new UI model
func NewModel(result *runner.TestExecutionResult) Model {
	return Model{
		failureGroups:   result.FailureGroups,
		testResults:     result.TestResults,
		selectedGroup:   0,
		selectedTest:    0,
		selectedFrame:   0,
		activePane:      0,
		showFullFrames:  false,
		result:          result,
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
		}
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
			_ = m.getCurrentFrames(group) // Get frames to validate bounds
			if m.selectedFrame > 0 {
				m.selectedFrame--
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
			frames := m.getCurrentFrames(group)
			if m.selectedFrame < len(frames)-1 {
				m.selectedFrame++
			}
		}
	}
	return m
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
		for i, group := range m.failureGroups {
			style := GetNormalTextStyle()
			if i == m.selectedGroup && isActive {
				style = GetSelectedTextStyle()
			}

			// Show error location (file:line) and count
			location := "Unknown"
			if len(group.NormalizedBacktrace) > 0 {
				frame := group.NormalizedBacktrace[len(group.NormalizedBacktrace)-1] // Bottom frame
				location = fmt.Sprintf("%s:%d", frame.File, frame.Line)
			}

			line := fmt.Sprintf("%s (%d failures)", location, group.Count)
			content.WriteString(style.Render(line))
			content.WriteString("\n")
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

			// Truncate long test names
			name := test.Name
			if len(name) > width-4 {
				name = name[:width-7] + "..."
			}

			content.WriteString(style.Render(name))
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
		frames := m.getCurrentFrames(group)

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

	paneStyle := GetPaneStyle(isActive).Width(width).Height(height)
	return paneStyle.Render(lipgloss.JoinVertical(lipgloss.Left, title, content.String()))
}

// renderStatusBar renders the status bar with keybindings
func (m Model) renderStatusBar() string {
	keyBindings := []string{
		GetKeyBindingStyle().Render("↑↓") + " navigate",
		GetKeyBindingStyle().Render("Tab") + " switch panes",
		GetKeyBindingStyle().Render("f") + " toggle frames",
		GetKeyBindingStyle().Render("q") + " quit",
	}

	statusText := strings.Join(keyBindings, " • ")
	return GetStatusBarStyle().Width(m.width).Render(statusText)
}
