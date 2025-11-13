package testrun

import (
	"fmt"
	"sort"
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

// TestRun captures the metadata needed to execute a run of tests
type TestRun struct {
	Id        int      // Unique identifier for the test run
	Filepaths []string // Specific test file patterns to execute
	Mode      string   // High-level mode describing how the run was initiated (optional)
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
func (tr *TestRuns) Add(filepaths []string, mode Mode) (TestRun, error) {
	testRun := TestRun{
		Id:        generateTestRunID(),
		Filepaths: filepaths,
		Mode:      string(mode),
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
