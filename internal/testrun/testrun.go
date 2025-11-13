package testrun

import (
	"fmt"
	"sort"
	"strings"
)

var testRunIDCounter int

func generateTestRunID() int {
	testRunIDCounter++
	return testRunIDCounter
}

// Mode represents the high-level mode describing how a test run was initiated
type Mode string

const (
	ModeRunWholeSuite       Mode = "run_whole_suite"
	ModeRunSelectedPatterns Mode = "run_selected_patterns"
	ModeReRunSingleFailure  Mode = "rerun_single_failure"
	ModeReRunAllFailures    Mode = "rerun_all_failures"
)

// TestPattern represents a test pattern with optional line number and test case name
type TestPattern struct {
	Path         string  // Required: directory or file path
	LineNumber   *int    // Optional: line number for future RSpec support
	TestCaseName *string // Optional: test case name for future support
}

// NewTestPattern creates a new TestPattern with validation
// Path is required. LineNumber and TestCaseName cannot both be set.
func NewTestPattern(path string, lineNumber *int, testCaseName *string) (TestPattern, error) {
	if path == "" {
		return TestPattern{}, fmt.Errorf("path is required")
	}
	if lineNumber != nil && testCaseName != nil {
		return TestPattern{}, fmt.Errorf("lineNumber and testCaseName cannot both be set")
	}
	return TestPattern{
		Path:         path,
		LineNumber:   lineNumber,
		TestCaseName: testCaseName,
	}, nil
}

// String returns the string representation of the pattern (just the path for now)
func (tp TestPattern) String() string {
	return tp.Path
}

// ParsePatternFromString parses a pattern string like "path.rb:123" and extracts the path
// Line number extraction is deferred for future implementation
func ParsePatternFromString(patternStr string) (TestPattern, error) {
	// For now, just extract the path (before the colon if present)
	parts := strings.Split(patternStr, ":")
	path := parts[0]
	return NewTestPattern(path, nil, nil)
}

// PatternsFromStrings converts a slice of string paths to TestPatterns
func PatternsFromStrings(paths []string) ([]TestPattern, error) {
	patterns := make([]TestPattern, 0, len(paths))
	for _, path := range paths {
		pattern, err := NewTestPattern(path, nil, nil)
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, pattern)
	}
	return patterns, nil
}

// PatternsToStrings converts a slice of TestPatterns to strings for command execution
func PatternsToStrings(patterns []TestPattern) []string {
	result := make([]string, len(patterns))
	for i, pattern := range patterns {
		result[i] = pattern.String()
	}
	return result
}

// TestRun captures the metadata needed to execute a run of tests
type TestRun struct {
	Id       int           // Unique identifier for the test run
	Patterns []TestPattern // Specific test patterns to execute
	Mode     string        // High-level mode describing how the run was initiated (optional)
}

// TestRuns is a collection of test runs
type TestRuns struct {
	testRuns map[int]TestRun
}

// NewTestRuns creates a new TestRuns collection
func NewTestRuns() TestRuns {
	return TestRuns{
		testRuns: make(map[int]TestRun),
	}
}

// Add creates and adds a new test run to the collection
func (tr *TestRuns) Add(patterns []TestPattern, mode Mode) (TestRun, error) {
	testRun := TestRun{
		Id:       generateTestRunID(),
		Patterns: patterns,
		Mode:     string(mode),
	}

	tr.testRuns[testRun.Id] = testRun
	return testRun, nil
}

// Get retrieves a test run by ID
func (tr *TestRuns) Get(id int) (TestRun, error) {
	testRun, ok := tr.testRuns[id]
	if !ok {
		return TestRun{}, fmt.Errorf("test run not found")
	}
	return testRun, nil
}

// MostRecent returns the most recent test run (highest ID)
func (tr *TestRuns) MostRecent() (TestRun, bool) {
	if len(tr.testRuns) == 0 {
		return TestRun{}, false
	}

	maxID := 0
	for id := range tr.testRuns {
		if id > maxID {
			maxID = id
		}
	}

	testRun, ok := tr.testRuns[maxID]
	return testRun, ok
}

// AllRecentFirst returns all test runs ordered by most recent first
func (tr *TestRuns) AllRecentFirst() []TestRun {
	orderedRecentFirst := make([]TestRun, 0, len(tr.testRuns))
	for _, testRun := range tr.testRuns {
		orderedRecentFirst = append(orderedRecentFirst, testRun)
	}

	sort.Slice(orderedRecentFirst, func(i, j int) bool {
		return orderedRecentFirst[i].Id > orderedRecentFirst[j].Id
	})

	return orderedRecentFirst
}
