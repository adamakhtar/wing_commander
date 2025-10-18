package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/adamakhtar/wing_commander/internal/types"
)


// InputTestResult represents a test result from JSON input
type InputTestResult struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Message     string `json:"message"`
	Backtrace   []string `json:"backtrace"`
	Duration    float64 `json:"duration,omitempty"`
	File        string `json:"file,omitempty"`
	Line        int    `json:"line,omitempty"`
}

// InputTestSuite represents the complete test suite output
type InputTestSuite struct {
	Tests []InputTestResult `json:"tests"`
	Summary struct {
		Total   int `json:"total"`
		Passed  int `json:"passed"`
		Failed  int `json:"failed"`
		Skipped int `json:"skipped"`
	} `json:"summary"`
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

// ParseFile parses a JSON test results file
func ParseFile(filePath string) (*ParseResult, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return ParseJSON(data)
}

// ParseJSON parses JSON test results data
func ParseJSON(data []byte) (*ParseResult, error) {
	var suite InputTestSuite

	// Try to parse as complete test suite first
	if err := json.Unmarshal(data, &suite); err != nil {
		// If that fails, try parsing as array of tests
		var tests []InputTestResult
		if err := json.Unmarshal(data, &tests); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
		suite.Tests = tests
	}

	// Convert to our domain types
	result := &ParseResult{
		Summary: TestSummary{
			Total:   suite.Summary.Total,
			Passed:  suite.Summary.Passed,
			Failed:  suite.Summary.Failed,
			Skipped: suite.Summary.Skipped,
		},
	}

	for _, inputTest := range suite.Tests {
		testResult := convertToTestResult(inputTest)
		result.Tests = append(result.Tests, testResult)
	}

	return result, nil
}

// convertToTestResult converts input format to our domain type
func convertToTestResult(input InputTestResult) types.TestResult {
	// Convert status string to our enum
	var status types.TestStatus
	switch strings.ToLower(input.Status) {
	case "pass", "passed", "success":
		status = types.StatusPass
	case "fail", "failed", "failure":
		status = types.StatusFail
	case "skip", "skipped", "pending":
		status = types.StatusSkip
	default:
		status = types.StatusFail // Default to fail for unknown statuses
	}

	// Parse backtrace into StackFrames
	var fullBacktrace []types.StackFrame
	for _, frameStr := range input.Backtrace {
		frame := parseStackFrame(frameStr)
		if frame.File != "" { // Only add frames with file info
			fullBacktrace = append(fullBacktrace, frame)
		}
	}

	// Cap at 50 frames
	if len(fullBacktrace) > 50 {
		fullBacktrace = fullBacktrace[:50]
	}

	return types.TestResult{
		Name:              input.Name,
		Status:            status,
		ErrorMessage:      input.Message,
		FullBacktrace:     fullBacktrace,
		FilteredBacktrace: make([]types.StackFrame, 0), // Will be populated by normalizer
	}
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
