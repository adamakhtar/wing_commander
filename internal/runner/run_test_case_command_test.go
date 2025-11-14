package runner

import (
	"testing"

	"github.com/adamakhtar/wing_commander/internal/testrun"
	"github.com/stretchr/testify/assert"
)

func TestBuildRunTestCaseCommand(t *testing.T) {
	testCaseName := "test_user_creation"
	lineNumber := 42
	cmd := BuildRunTestCaseCommand(
		"bundle exec ruby -Itest %{test_case_name}:%{line_number}",
		testrun.TestPattern{
			Path:         "test/user_test.rb",
			TestCaseName: &testCaseName,
			LineNumber:   &lineNumber,
			TestGroupName: nil,
		},
	)

	assert.Equal(t, "bundle exec ruby -Itest test_user_creation:42", cmd)
}

func TestBuildRunTestCaseCommandWithFilePath(t *testing.T) {
	testCaseName := "test_user_creation"
	lineNumber := 42
	cmd := BuildRunTestCaseCommand(
		"bundle exec rake test %{test_file_path} -n '/%{test_case_name}/'",
		testrun.TestPattern{
			Path:         "test/user_test.rb",
			TestCaseName: &testCaseName,
			LineNumber:   &lineNumber,
			TestGroupName: nil,
		},
	)

	assert.Equal(t, "bundle exec rake test test/user_test.rb -n '/test_user_creation/'", cmd)
}

func TestBuildRunTestCaseCommandEmptyTemplate(t *testing.T) {
	testCaseName := "test_user_creation"
	lineNumber := 42
	cmd := BuildRunTestCaseCommand("", testrun.TestPattern{
		Path:         "test/user_test.rb",
		TestCaseName: &testCaseName,
		LineNumber:   &lineNumber,
		TestGroupName: nil,
	})

	assert.Equal(t, "", cmd)
}

func TestBuildRunTestCaseCommandMissingValues(t *testing.T) {
	cmd := BuildRunTestCaseCommand("bundle exec ruby %{test_case_name}:%{line_number}", testrun.TestPattern{
		Path:         "test/user_test.rb",
		TestCaseName: nil,
		LineNumber:   nil,
		TestGroupName: nil,
	})

	assert.Equal(t, "bundle exec ruby :", cmd)
}

func TestBuildRunTestCaseCommandNilTestCaseName(t *testing.T) {
	lineNumber := 42
	cmd := BuildRunTestCaseCommand(
		"bundle exec ruby %{test_file_path}:%{line_number}",
		testrun.TestPattern{
			Path:         "test/user_test.rb",
			TestCaseName: nil,
			LineNumber:   &lineNumber,
			TestGroupName: nil,
		},
	)

	assert.Equal(t, "bundle exec ruby test/user_test.rb:42", cmd)
}

func TestBuildRunTestCaseCommandNilLineNumber(t *testing.T) {
	testCaseName := "test_user_creation"
	cmd := BuildRunTestCaseCommand(
		"bundle exec ruby %{test_file_path} -n '/%{test_case_name}/'",
		testrun.TestPattern{
			Path:         "test/user_test.rb",
			TestCaseName: &testCaseName,
			LineNumber:   nil,
			TestGroupName: nil,
		},
	)

	assert.Equal(t, "bundle exec ruby test/user_test.rb -n '/test_user_creation/'", cmd)
}
