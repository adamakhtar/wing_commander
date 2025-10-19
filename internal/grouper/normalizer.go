package grouper

import (
	"strings"

	"github.com/adamakhtar/wing_commander/internal/config"
	"github.com/adamakhtar/wing_commander/internal/types"
)

// Normalizer handles backtrace filtering and normalization
type Normalizer struct {
	excludePatterns []string
}

// NewNormalizer creates a new Normalizer with the given exclude patterns
func NewNormalizer(cfg *config.Config) *Normalizer {
	return &Normalizer{
		excludePatterns: cfg.ExcludePatterns,
	}
}

// NormalizeTestResults processes all test results and filters their backtraces
func (n *Normalizer) NormalizeTestResults(results []types.TestResult) []types.TestResult {
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

// shouldExclude checks if a frame matches any exclude pattern
func (n *Normalizer) shouldExclude(frame types.StackFrame) bool {
	for _, pattern := range n.excludePatterns {
		if strings.Contains(frame.File, pattern) {
			return true
		}
	}
	return false
}

// NormalizeFrameForGrouping creates a normalized string for grouping
// This removes line numbers and focuses on file + method
func NormalizeFrameForGrouping(frame types.StackFrame) string {
	// For grouping, we use file + method (ignore line number)
	if frame.Function != "" {
		return frame.File + "::" + frame.Function
	}
	// If no function name, use just the file
	return frame.File
}

// NormalizeBacktraceForGrouping creates a signature for a backtrace
// This is used to group tests with similar call stacks
func NormalizeBacktraceForGrouping(frames []types.StackFrame) string {
	if len(frames) == 0 {
		return ""
	}

	// Build signature from normalized frames
	parts := make([]string, 0, len(frames))
	for _, frame := range frames {
		parts = append(parts, NormalizeFrameForGrouping(frame))
	}

	return strings.Join(parts, "|")
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
