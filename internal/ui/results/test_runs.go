package results

import (
	"fmt"

	"github.com/adamakhtar/wing_commander/internal/runner"
)

var testRunIdCounter int = 0

func generateTestRunId() int {
	testRunIdCounter++
	return testRunIdCounter
}

type TestRuns struct {
	testRuns map[int]TestRun
}

type TestRun struct {
	Id int
	Filepaths []string
	Result *runner.TestExecutionResult
	Error string
}

func NewTestRuns() TestRuns {
	return TestRuns{
		testRuns: make(map[int]TestRun),
	}
}

func (tr *TestRuns) Add(filepaths []string) (TestRun, error) {
	testRun := TestRun{
		Id: generateTestRunId(),
		Filepaths: filepaths,
	}
	_, ok := tr.testRuns[testRun.Id]
	if ok {
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

func (tr *TestRuns) UpdateError(id int, error string) (TestRun, error) {
	testRun, err := tr.Get(id)
	if err != nil {
		return TestRun{}, err
	}
	testRun.Error = error
	return testRun, nil
}

func (tr *TestRuns) UpdateResult(id int, result *runner.TestExecutionResult) (TestRun, error) {
	testRun, err := tr.Get(id)
	if err != nil {
		return TestRun{}, err
	}
	testRun.Result = result
	return testRun, nil
}

func (tr *TestRuns) MostRecent() (TestRun, bool) {
	if len(tr.testRuns) == 0 {
		return TestRun{}, false
	}

	maxId := 0
	for id := range tr.testRuns {
		if id > maxId {
			maxId = id
		}
	}

	testRun, ok := tr.testRuns[maxId]
	return testRun, ok
}