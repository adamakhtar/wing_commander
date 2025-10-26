package context

type ScreenType string

const (
	ResultsScreen ScreenType = "resultsScreen"
	FilePickerScreen ScreenType = "filePickerScreen"
)

type Context struct {
	ScreenWidth int
	ScreenHeight int
	CurrentScreen ScreenType
}
