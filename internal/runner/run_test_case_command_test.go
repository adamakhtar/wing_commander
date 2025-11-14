package runner

import (
	"testing"

	"github.com/adamakhtar/wing_commander/internal/testrun"
	"github.com/stretchr/testify/assert"
)

func TestBuildRunTestCaseCommand(t *testing.T) {
	cmd := BuildRunTestCaseCommand(
		"bundle exec rake test",
		[]testrun.TestPattern{},
	)

	assert.Equal(t, "bundle exec rake test", cmd)
}

func TestBuildRunTestCaseCommandEmptyCommand(t *testing.T) {
	cmd := BuildRunTestCaseCommand("", []testrun.TestPattern{
		{Path: "test/user_test.rb"},
	})

	assert.Equal(t, "", cmd)
}

func TestBuildRunTestCaseCommandWithOnlyPath(t *testing.T) {
	cmd := BuildRunTestCaseCommand(
		"bundle exec rake test",
		[]testrun.TestPattern{
			{Path: "test/worker_test.rb"},
		},
	)

	assert.Equal(t, "bundle exec rake test test/worker_test.rb", cmd)
}

func TestBuildRunTestCaseCommandWithPathAndTestCase(t *testing.T) {
	testCaseName := "test_success"
	cmd := BuildRunTestCaseCommand(
		"bundle exec rake test",
		[]testrun.TestPattern{
			{
				Path:         "test/worker_test.rb",
				TestCaseName: &testCaseName,
			},
		},
	)

	assert.Equal(t, "bundle exec rake test test/worker_test.rb:test_success", cmd)
}

func TestBuildRunTestCaseCommandWithPathTestCaseAndGroup(t *testing.T) {
	testCaseName := "test_success"
	groupName := "AGroupName"
	cmd := BuildRunTestCaseCommand(
		"bundle exec rake test",
		[]testrun.TestPattern{
			{
				Path:          "test/worker_test.rb",
				TestCaseName:  &testCaseName,
				TestGroupName: &groupName,
			},
		},
	)

	assert.Equal(t, "bundle exec rake test test/worker_test.rb:AGroupName#test_success", cmd)
}

func TestBuildRunTestCaseCommandMultiplePatterns(t *testing.T) {
	testCaseName1 := "test_one"
	testCaseName2 := "test_two"
	groupName := "MyGroup"
	cmd := BuildRunTestCaseCommand(
		"bundle exec rake test",
		[]testrun.TestPattern{
			{Path: "test/worker_test.rb"},
			{
				Path:         "test/user_test.rb",
				TestCaseName: &testCaseName1,
			},
			{
				Path:          "test/admin_test.rb",
				TestCaseName:  &testCaseName2,
				TestGroupName: &groupName,
			},
		},
	)

	assert.Equal(t, "bundle exec rake test test/worker_test.rb test/user_test.rb:test_one test/admin_test.rb:MyGroup#test_two", cmd)
}
