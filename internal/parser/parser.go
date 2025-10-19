package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/joshdk/go-junit"
)


// InputTestResult represents a test result from JUnit XML input
type InputTestResult struct {
	Name        string
	ClassName   string
	Status      string
	Message     string
	Backtrace   []string
	Duration    float64
	File        string
	Line        int
}

// ParseResult contains parsed test results and metadata
type ParseResult struct {
	Tests   []types.TestResult
	Summary TestSummary
}

// TestSummary contains test run statistics
type TestSummary struct {
	Total   int
	Passed  int
	Failed  int
	Skipped int
}

// ParseFile parses a JUnit XML test results file
func ParseFile(filePath string) (*ParseResult, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return ParseXML(data)
}

// ParseXML parses JUnit XML test results data
func ParseXML(data []byte) (*ParseResult, error) {
	suites, err := junit.Ingest(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JUnit XML: %w", err)
	}

	result := &ParseResult{
		Summary: TestSummary{},
	}

	totalTests := 0
	failedTests := 0
	skippedTests := 0

	for _, suite := range suites {
		totalTests += suite.Totals.Tests
		failedTests += suite.Totals.Failed + suite.Totals.Error
		skippedTests += suite.Totals.Skipped

		for _, test := range suite.Tests {
			testResult := convertJUnitTestToTestResult(test)
			result.Tests = append(result.Tests, testResult)
		}
	}

	result.Summary = TestSummary{
		Total:   totalTests,
		Passed:  totalTests - failedTests - skippedTests,
		Failed:  failedTests,
		Skipped: skippedTests,
	}

	return result, nil
}

// convertJUnitTestToTestResult converts JUnit test to our domain type
func convertJUnitTestToTestResult(test junit.Test) types.TestResult {
	// Determine status based on JUnit test result
	var status types.TestStatus
	var errorMessage string
	var backtrace []string

	switch test.Status {
	case junit.StatusPassed:
		status = types.StatusPass
	case junit.StatusSkipped:
		status = types.StatusSkip
		errorMessage = test.Message
	case junit.StatusFailed, junit.StatusError:
		status = types.StatusFail
		errorMessage = test.Message
		// Parse stacktrace from error output
		if test.Error != nil {
			backtrace = parseStacktraceFromError(test.Error.Error())
		}
		// Also check SystemErr for additional stacktrace info
		if test.SystemErr != "" {
			additionalBacktrace := parseStacktraceFromError(test.SystemErr)
			backtrace = append(backtrace, additionalBacktrace...)
		}
	default:
		status = types.StatusFail
		errorMessage = "Unknown test status"
	}

	// Parse backtrace into StackFrames
	var fullBacktrace []types.StackFrame
	for _, frameStr := range backtrace {
		frame := parseStackFrame(frameStr)
		if frame.File != "" { // Only add frames with file info
			fullBacktrace = append(fullBacktrace, frame)
		}
	}

	// Cap at 50 frames
	if len(fullBacktrace) > 50 {
		fullBacktrace = fullBacktrace[:50]
	}

	// Create test name combining classname and name if available
	testName := test.Name
	if test.Classname != "" && test.Classname != test.Name {
		testName = fmt.Sprintf("%s.%s", test.Classname, test.Name)
	}

	return types.TestResult{
		Name:              testName,
		Status:            status,
		ErrorMessage:      errorMessage,
		FullBacktrace:     fullBacktrace,
		FilteredBacktrace: make([]types.StackFrame, 0), // Will be populated by normalizer
	}
}

// parseStacktraceFromError extracts stacktrace lines from error output
func parseStacktraceFromError(output string) []string {
	if output == "" {
		return []string{}
	}

	lines := strings.Split(output, "\n")
	var stacktrace []string

	for _, line := range lines {
		// Don't trim the line yet - we need to check indentation
		// Skip empty lines and common non-stacktrace lines
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "Error:") || strings.HasPrefix(strings.TrimSpace(line), "Failure:") {
			continue
		}
		// Look for lines that look like stack frames (indented with 4 spaces and contain file:line pattern)
		if strings.HasPrefix(line, "    ") && strings.Contains(line, ":") &&
		   (strings.Contains(line, ".rb:") || strings.Contains(line, ".py:") || strings.Contains(line, ".js:")) {
			stacktrace = append(stacktrace, strings.TrimSpace(line))
		}
	}

	return stacktrace
}

// parseStackFrame parses a backtrace frame string into a StackFrame
func parseStackFrame(frameStr string) types.StackFrame {
	// Common formats:
	// "app/models/user.rb:42:in `create_user'"
	// "app/models/user.rb:42"
	// "File \"app/models/user.rb\", line 42, in create_user"

	// Handle Python format first
	if strings.HasPrefix(frameStr, "File \"") {
		return types.StackFrame{
			File:     frameStr,
			Line:     0,
			Function: "",
		}
	}

	parts := strings.Split(frameStr, ":")
	if len(parts) < 2 {
		return types.StackFrame{File: frameStr}
	}

	file := parts[0]

	// Try to extract line number
	var line int
	var function string

	if len(parts) >= 2 {
		// Parse line number
		if _, err := fmt.Sscanf(parts[1], "%d", &line); err != nil {
			return types.StackFrame{File: file}
		}
	}

	// Try to extract function name
	if len(parts) >= 3 {
		funcPart := parts[2]
		// Remove "in `" and "`" wrapper, or "in '" and "'" wrapper
		if strings.HasPrefix(funcPart, "in `") && strings.HasSuffix(funcPart, "'") {
			function = funcPart[4 : len(funcPart)-1]
		} else if strings.HasPrefix(funcPart, "in '") && strings.HasSuffix(funcPart, "'") {
			function = funcPart[4 : len(funcPart)-1]
		} else if strings.HasPrefix(funcPart, "in ") {
			function = funcPart[3:]
		}
	}

	return types.StackFrame{
		File:     file,
		Line:     line,
		Function: function,
	}
}
