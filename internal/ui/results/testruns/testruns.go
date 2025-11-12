package testruns

import (
	"fmt"
	"sort"

	"github.com/adamakhtar/wing_commander/internal/runner"
)

var testRunIDCounter int

func generateTestRunID() int {
	testRunIDCounter++
	return testRunIDCounter
}

type TestRuns struct {
	testRuns map[int]TestRun
}

type TestRun struct {
	Id        int
	Filepaths []string
	Result    *runner.TestExecutionResult
	Error     string
}

func NewTestRuns() TestRuns {
	return TestRuns{
		testRuns: make(map[int]TestRun),
	}
}

func (tr *TestRuns) Add(filepaths []string) (TestRun, error) {
	testRun := TestRun{
		Id:        generateTestRunID(),
		Filepaths: filepaths,
	}

	if _, exists := tr.testRuns[testRun.Id]; exists {
		return TestRun{}, fmt.Errorf("test run already exists")
	}

	tr.testRuns[testRun.Id] = testRun
	return testRun, nil
}

func (tr *TestRuns) Get(id int) (TestRun, error) {
	testRun, ok := tr.testRuns[id]
	if !ok {
		return TestRun{}, fmt.Errorf("test run not found")
	}
	return testRun, nil
}

func (tr *TestRuns) UpdateError(id int, errMsg string) (TestRun, error) {
	testRun, err := tr.Get(id)
	if err != nil {
		return TestRun{}, err
	}

	testRun.Error = errMsg
	tr.testRuns[id] = testRun
	return testRun, nil
}

func (tr *TestRuns) UpdateResult(id int, result *runner.TestExecutionResult) (TestRun, error) {
	testRun, err := tr.Get(id)
	if err != nil {
		return TestRun{}, err
	}

	testRun.Result = result
	tr.testRuns[id] = testRun
	return testRun, nil
}

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
