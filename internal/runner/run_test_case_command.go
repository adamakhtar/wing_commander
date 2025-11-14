package runner

import (
	"strings"

	"github.com/adamakhtar/wing_commander/internal/testrun"
)

// BuildRunTestCaseCommand builds a command by appending space-delimited pattern strings to the command.
// Each pattern is converted to its string representation via TestPattern.String().
func BuildRunTestCaseCommand(command string, patterns []testrun.TestPattern) string {
	if command == "" {
		return ""
	}

	if len(patterns) == 0 {
		return command
	}

	patternStrings := make([]string, len(patterns))
	for i, pattern := range patterns {
		patternStrings[i] = pattern.String()
	}

	return command + " " + "\"" + strings.Join(patternStrings, " ") + "\""
}
