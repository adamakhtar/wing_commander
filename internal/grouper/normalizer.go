package grouper

import (
	"strings"

	"github.com/adamakhtar/wing_commander/internal/config"
	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/charmbracelet/log"
)

// Normalizer handles backtrace filtering and normalization
type Normalizer struct {
	projectPath string
}

// NewNormalizer creates a new Normalizer with the given config
func NewNormalizer(cfg *config.Config) *Normalizer {
	return &Normalizer{
		projectPath: cfg.ProjectPath,
	}
}

// NormalizeTestResults processes all test results and filters their backtraces
func (n *Normalizer) NormalizeTestResults(results []types.TestResult) []types.TestResult {
	log.Debug("Normalizing test results", "projectPath", n.projectPath)
	normalized := make([]types.TestResult, len(results))

	for i, result := range results {
		normalized[i] = n.normalizeTestResult(result)
	}

	return normalized
}

// normalizeTestResult filters a single test result's backtrace
func (n *Normalizer) normalizeTestResult(result types.TestResult) types.TestResult {
	// Filter the backtrace to include only project frames
	filtered := n.filterBacktrace(result.FullBacktrace)

	result.FilteredBacktrace = filtered
	return result
}

// filterBacktrace removes frames matching exclude patterns
func (n *Normalizer) filterBacktrace(frames []types.StackFrame) []types.StackFrame {
	filtered := make([]types.StackFrame, 0, len(frames))

	for _, frame := range frames {
		if !n.shouldExclude(frame) {
			filtered = append(filtered, frame)
		}
	}

	return filtered
}

// shouldExclude checks if a frame should be excluded (excluded if filepath does not start with projectPath)
func (n *Normalizer) shouldExclude(frame types.StackFrame) bool {
	if n.projectPath == "" {
		return false
	}
	return !strings.HasPrefix(frame.File, n.projectPath)
}


// GetProjectFrames returns only project-level frames (filtered)
func GetProjectFrames(result types.TestResult) []types.StackFrame {
	if len(result.FilteredBacktrace) > 0 {
		return result.FilteredBacktrace
	}
	// Fallback to full backtrace if no filtered version
	return result.FullBacktrace
}

// CountFilteredFrames returns statistics about frame filtering
func CountFilteredFrames(results []types.TestResult) (totalFrames, filteredFrames int) {
	for _, result := range results {
		totalFrames += len(result.FullBacktrace)
		filteredFrames += len(result.FilteredBacktrace)
	}
	return totalFrames, filteredFrames
}
