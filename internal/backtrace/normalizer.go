package backtrace

import (
	"strings"

	"github.com/adamakhtar/wing_commander/internal/config"
	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/charmbracelet/log"
)

// Normalizer handles backtrace filtering and normalization.
type Normalizer struct {
	projectPath string
}

// NewNormalizer creates a new Normalizer with the given config.
func NewNormalizer(cfg *config.Config) *Normalizer {
	return &Normalizer{
		projectPath: cfg.ProjectPath,
	}
}

// NormalizeTestResults processes all test results and filters their backtraces.
func (n *Normalizer) NormalizeTestResults(results []types.TestResult) []types.TestResult {
	log.Debug("Normalizing test results", "projectPath", n.projectPath)
	normalized := make([]types.TestResult, len(results))

	for i, result := range results {
		normalized[i] = n.normalizeTestResult(result)
	}

	return normalized
}

func (n *Normalizer) normalizeTestResult(result types.TestResult) types.TestResult {
	filtered := n.filterBacktrace(result.FullBacktrace)

	result.FilteredBacktrace = filtered
	return result
}

func (n *Normalizer) filterBacktrace(frames []types.StackFrame) []types.StackFrame {
	filtered := make([]types.StackFrame, 0, len(frames))

	for _, frame := range frames {
		if !n.shouldExclude(frame) {
			filtered = append(filtered, frame)
		}
	}

	return filtered
}

func (n *Normalizer) shouldExclude(frame types.StackFrame) bool {
	if n.projectPath == "" {
		return false
	}
	return !strings.HasPrefix(frame.File, n.projectPath)
}
