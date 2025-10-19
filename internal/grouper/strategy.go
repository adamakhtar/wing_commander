package grouper

import (
	"fmt"

	"github.com/adamakhtar/wing_commander/internal/types"
)

// Strategy defines the interface for grouping strategies
type Strategy interface {
	// GroupKey generates a unique key for grouping based on stack frames
	GroupKey(frames []types.StackFrame) string
}

// ErrorLocationStrategy groups failures by the bottom frame (where the error surfaced)
// This represents the location where the error actually occurred
type ErrorLocationStrategy struct{}

// NewErrorLocationStrategy creates a new ErrorLocationStrategy
func NewErrorLocationStrategy() *ErrorLocationStrategy {
	return &ErrorLocationStrategy{}
}

// GroupKey generates a grouping key based on the bottom frame of the backtrace
// Key format: "{filename}:{line_number}"
// This groups failures that occurred at the same file and line
func (s *ErrorLocationStrategy) GroupKey(frames []types.StackFrame) string {
	if len(frames) == 0 {
		return ""
	}

	// Get the bottom frame (first frame in the backtrace)
	// In most backtraces, the first frame is where the error originated
	// (deepest in the call stack), and subsequent frames show the call path
	bottomFrame := frames[0]

	// Generate key from filename and line number
	// Include line number as it's important for precise error location
	return fmt.Sprintf("%s:%d", bottomFrame.File, bottomFrame.Line)
}
