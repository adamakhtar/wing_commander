package testresult

import (
	"github.com/adamakhtar/wing_commander/internal/projectfs"
	"github.com/charmbracelet/log"
)

// Normalizer handles backtrace filtering and normalization.
type Normalizer struct {
}

// NewNormalizer creates a new Normalizer.
func NewNormalizer() *Normalizer {
	return &Normalizer{}
}

// NormalizeTestResults processes all test results and filters their backtraces.
func (n *Normalizer) NormalizeTestResults(results []TestResult) []TestResult {
	fs := projectfs.GetProjectFS()
	log.Debug("Normalizing test results", "projectPath", fs.RootPath.String())
	normalized := make([]TestResult, len(results))

	for i, result := range results {
		normalized[i] = n.normalizeTestResult(result)
	}

	return normalized
}

func (n *Normalizer) normalizeTestResult(result TestResult) TestResult {
	result.FilteredBacktrace = result.FullBacktrace.FilterProjectStackFramesOnly()
	return result
}
