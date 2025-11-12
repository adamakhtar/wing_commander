package runner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildRunTestCaseCommand(t *testing.T) {
	cmd := BuildRunTestCaseCommand(
		"bundle exec ruby -Itest %{test_case_name}:%{line_number}",
		RunTestCaseParams{
			TestCaseName: "test_user_creation",
			LineNumber:   42,
		},
	)

	assert.Equal(t, "bundle exec ruby -Itest test_user_creation:42", cmd)
}

func TestBuildRunTestCaseCommandEmptyTemplate(t *testing.T) {
	cmd := BuildRunTestCaseCommand("", RunTestCaseParams{
		TestCaseName: "test_user_creation",
		LineNumber:   42,
	})

	assert.Equal(t, "", cmd)
}

func TestBuildRunTestCaseCommandMissingValues(t *testing.T) {
	cmd := BuildRunTestCaseCommand("bundle exec ruby %{test_case_name}:%{line_number}", RunTestCaseParams{
		TestCaseName: "",
		LineNumber:   0,
	})

	assert.Equal(t, "bundle exec ruby :", cmd)
}
