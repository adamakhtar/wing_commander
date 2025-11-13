package runner

import (
	"strconv"
	"strings"

	"github.com/adamakhtar/wing_commander/internal/testrun"
)

const (
	placeholderTestCaseName = "%{test_case_name}"
	placeholderLineNumber   = "%{line_number}"
	placeholderTestFilePath = "%{test_file_path}"
)

// BuildRunTestCaseCommand resolves the placeholders in a run-test-case command template.
// Supported placeholders:
//   - %{test_file_path}
//   - %{test_case_name}
//   - %{line_number}
//
// Example: "bundle exec rake test %{test_file_path} -n '/%{test_case_name}/'" =>
//
//	"bundle exec rake test test/user_test.rb -n '/test_user_creation/'"
func BuildRunTestCaseCommand(template string, pattern testrun.TestPattern) string {
	if template == "" {
		return ""
	}

	testCaseName := ""
	if pattern.TestCaseName != nil {
		testCaseName = *pattern.TestCaseName
	}

	lineNumber := ""
	if pattern.LineNumber != nil {
		lineNumber = lineNumberString(*pattern.LineNumber)
	}

	replacer := strings.NewReplacer(
		placeholderTestFilePath, pattern.Path,
		placeholderTestCaseName, testCaseName,
		placeholderLineNumber, lineNumber,
	)

	return replacer.Replace(template)
}

func lineNumberString(line int) string {
	if line <= 0 {
		return ""
	}
	return strconv.Itoa(line)
}
