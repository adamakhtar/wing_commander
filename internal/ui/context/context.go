package context

import (
	"github.com/adamakhtar/wing_commander/internal/config"
	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/adamakhtar/wing_commander/internal/ui/styles"
)

type ScreenType string

const (
	ResultsScreen ScreenType = "resultsScreen"
	FilePickerScreen ScreenType = "filePickerScreen"
)

type Context struct {
	Config *config.Config
	Styles styles.Styles
	ScreenWidth int
	ScreenHeight int
	CurrentScreen ScreenType
	SelectedTestResult *types.TestResult
}
