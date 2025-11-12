package runner

import (
	"strconv"
	"strings"
)

const (
	placeholderTestCaseName = "%{test_case_name}"
	placeholderLineNumber   = "%{line_number}"
)

// RunTestCaseParams describes the values that can be substituted into a run-test-case command template.
type RunTestCaseParams struct {
	TestCaseName string
	LineNumber   int
}

// BuildRunTestCaseCommand resolves the placeholders in a run-test-case command template.
// Supported placeholders:
//   - %{test_case_name}
//   - %{line_number}
//
// Example: "bundle exec rake test -n '/%{test_case_name}/'" =>
//
//	"bundle exec rake test -n '/test_user_creation/'"
func BuildRunTestCaseCommand(template string, params RunTestCaseParams) string {
	if template == "" {
		return ""
	}

	replacer := strings.NewReplacer(
		placeholderTestCaseName, params.TestCaseName,
		placeholderLineNumber, lineNumberString(params.LineNumber),
	)

	return replacer.Replace(template)
}

func lineNumberString(line int) string {
	if line <= 0 {
		return ""
	}
	return strconv.Itoa(line)
}
