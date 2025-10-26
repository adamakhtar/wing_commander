package results

import (
	"github.com/adamakhtar/wing_commander/internal/editor"
	"github.com/adamakhtar/wing_commander/internal/runner"
)

// Model represents the UI state
type Model struct {
	// Data
	ctx Context
}

// NewModel creates a new UI model
func NewModel(result *runner.TestExecutionResult, testRunner *runner.TestRunner) Model {
	model := Model{
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

func UpdateContext(ctx Context) {
	m.ctx = ctx
}