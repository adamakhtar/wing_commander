package grouper

import (
	"sort"

	"github.com/adamakhtar/wing_commander/internal/types"
)

// Grouper handles grouping of test failures using a configurable strategy
type Grouper struct {
	strategy Strategy
}

// NewGrouper creates a new Grouper with the specified strategy
func NewGrouper(strategy Strategy) *Grouper {
	return &Grouper{
		strategy: strategy,
	}
}

// GroupFailures groups test results by their failure characteristics using the configured strategy
// Returns groups sorted by count (descending) - most frequent failures first
func (g *Grouper) GroupFailures(results []types.TestResult) []types.FailureGroup {
	// Filter to only failed tests
	failedTests := filterFailedTests(results)
	if len(failedTests) == 0 {
		return []types.FailureGroup{}
	}

	// Group tests by strategy-generated key
	groupMap := make(map[string]*types.FailureGroup)

	for _, test := range failedTests {
		// Use filtered backtrace for grouping (project frames only)
		frames := GetProjectFrames(test)
		groupKey := g.strategy.GroupKey(frames)

		// Skip tests with empty group keys (no valid frames)
		if groupKey == "" {
			continue
		}

		// Get or create group for this key
		group, exists := groupMap[groupKey]
		if !exists {
			group = &types.FailureGroup{
				Hash:                groupKey,
				ErrorMessage:        test.ErrorMessage,
				NormalizedBacktrace: frames,
				Tests:               []types.TestResult{},
				Count:               0,
			}
			groupMap[groupKey] = group
		}

		// Add test to group
		group.Tests = append(group.Tests, test)
		group.Count++
	}

	// Convert map to slice
	groups := make([]types.FailureGroup, 0, len(groupMap))
	for _, group := range groupMap {
		groups = append(groups, *group)
	}

	// Sort groups by count (descending) - most frequent failures first
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Count > groups[j].Count
	})

	return groups
}

// filterFailedTests returns only test results that have failed
func filterFailedTests(results []types.TestResult) []types.TestResult {
	failed := make([]types.TestResult, 0, len(results))
	for _, result := range results {
		if result.Status == types.StatusFail {
			failed = append(failed, result)
		}
	}
	return failed
}
