package context

import "github.com/adamakhtar/wing_commander/internal/config"

type ScreenType string

const (
	ResultsScreen ScreenType = "resultsScreen"
	FilePickerScreen ScreenType = "filePickerScreen"
)

type Context struct {
	Config *config.Config
	ScreenWidth int
	ScreenHeight int
	CurrentScreen ScreenType
}
