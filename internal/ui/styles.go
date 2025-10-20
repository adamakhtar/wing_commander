package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color definitions
var (
	// Pane styles
	activePaneBorder   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62"))
	inactivePaneBorder = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("8"))

	// Text styles
	normalText     = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	dimmedText     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	selectedText   = lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Bold(true)

	// Change intensity styles (same color, varying boldness)
	changeIntensity0 = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))                    // Normal
	changeIntensity1 = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(false)        // Light
	changeIntensity2 = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)         // Medium
	changeIntensity3 = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)         // Strong

	// Title styles
	paneTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Bold(true)

	// Status bar styles
	statusBar = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Background(lipgloss.Color("0"))
	keyBinding = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)

	// Success/error styles
	successText = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	errorText   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	yellowText  = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
)

// GetPaneStyle returns the appropriate border style for a pane
func GetPaneStyle(isActive bool) lipgloss.Style {
	if isActive {
		return activePaneBorder
	}
	return inactivePaneBorder
}

// GetChangeIntensityStyle returns the appropriate text style for change intensity
func GetChangeIntensityStyle(intensity int) lipgloss.Style {
	switch intensity {
	case 0:
		return changeIntensity0
	case 1:
		return changeIntensity1
	case 2:
		return changeIntensity2
	case 3:
		return changeIntensity3
	default:
		return changeIntensity0
	}
}

// GetPaneTitleStyle returns the style for pane titles
func GetPaneTitleStyle() lipgloss.Style {
	return paneTitle
}

// GetNormalTextStyle returns the style for normal text
func GetNormalTextStyle() lipgloss.Style {
	return normalText
}

// GetDimmedTextStyle returns the style for dimmed text
func GetDimmedTextStyle() lipgloss.Style {
	return dimmedText
}

// GetSelectedTextStyle returns the style for selected text
func GetSelectedTextStyle() lipgloss.Style {
	return selectedText
}

// GetStatusBarStyle returns the style for the status bar
func GetStatusBarStyle() lipgloss.Style {
	return statusBar
}

// GetKeyBindingStyle returns the style for key bindings
func GetKeyBindingStyle() lipgloss.Style {
	return keyBinding
}

// GetSuccessTextStyle returns the style for success messages
func GetSuccessTextStyle() lipgloss.Style {
	return successText
}

// GetErrorTextStyle returns the style for error messages
func GetErrorTextStyle() lipgloss.Style {
	return errorText
}

// GetYellowTextStyle returns the style for yellow text
func GetYellowTextStyle() lipgloss.Style {
	return yellowText
}
