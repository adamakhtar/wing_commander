package runner

import (
	"fmt"
	"strings"

	"github.com/adamakhtar/wing_commander/internal/testrun"
)

// shellEscape escapes a string for safe use in shell commands.
// Wraps the string in single quotes and escapes any single quotes within it.
func shellEscape(s string) string {
	if s == "" {
		return "''"
	}
	// Replace single quotes with: '\'' (end quote, escaped quote, start quote)
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func shellEscapeList(strings []string) []string {
	escapedStrings := make([]string, len(strings))
	for i, s := range strings {
		escapedStrings[i] = shellEscape(s)
	}
	return escapedStrings
}

// BuildRunTestCaseCommand builds a command by appending space-delimited pattern strings to the command.
// Each pattern is converted to its string representation via TestPattern.String().
func BuildRunTestCaseCommand(command string, testRun testrun.TestRun) (string, error) {
	if command == "" {
		return "", fmt.Errorf("command is required")
	}

	switch testRun.Mode {
	case string(testrun.ModeRunWholeSuite):
		return command, nil

	case string(testrun.ModeRunSelectedPatterns):
		filePathStrings := testRun.PatternsToFilePaths()
		escapedPaths := shellEscapeList(filePathStrings)
		return command + " " + strings.Join(escapedPaths, " "), nil

	case string(testrun.ModeReRunSingleFailure), string(testrun.ModeReRunAllFailures):
		testCaseStrings := testRun.PatternsToTestCaseIdentifiers()
		commaSeparatedTestCases := strings.Join(testCaseStrings, ",")
		escapedTestCases := shellEscape(commaSeparatedTestCases)
		return command + " --test-cases " + escapedTestCases, nil

	default:
		return "", fmt.Errorf("invalid mode: %s", testRun.Mode)
	}
}
