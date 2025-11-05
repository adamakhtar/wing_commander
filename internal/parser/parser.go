package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/joshdk/go-junit"
	"gopkg.in/yaml.v3"
)

// testPathIndicators lists substrings that indicate frames belonging to test code paths
// Limited to Ruby frameworks for now (RSpec and Minitest). Extend as more frameworks are supported.
var testPathIndicators = []string{
    "/spec/",
    "/test/",
    "_spec.rb",
}

// frameworkPathIndicators lists substrings that indicate frames from test frameworks/runners
var frameworkPathIndicators = []string{
    "/rspec/",
    "/minitest/",
    "/junit/",
    "/jest/",
    "/node_modules/",
    "/gems/",
}


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

// testcaseAttrs represents XML attributes for a testcase element
type testcaseAttrs struct {
	File   string `xml:"file,attr"`
	LineNo string `xml:"lineno,attr"`
}

// ParseXML parses JUnit XML test results data
func ParseXML(data []byte) (*ParseResult, error) {
	// First, parse XML directly to extract file and lineno attributes
	testcaseMap := make(map[string]testcaseAttrs)
	decoder := xml.NewDecoder(strings.NewReader(string(data)))
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		if startElem, ok := token.(xml.StartElement); ok && startElem.Name.Local == "testcase" {
			var attrs testcaseAttrs
			var name, classname string
			for _, attr := range startElem.Attr {
				switch attr.Name.Local {
				case "file":
					attrs.File = attr.Value
				case "lineno":
					attrs.LineNo = attr.Value
				case "name":
					name = attr.Value
				case "classname":
					classname = attr.Value
				}
			}
			// Use classname + name as key for uniqueness (matching junit library behavior)
			key := name
			if classname != "" && classname != name {
				key = fmt.Sprintf("%s %s", classname, name)
			}
			if key != "" {
				testcaseMap[key] = attrs
			}
		}
	}

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
	testIdCounter := 0

	for _, suite := range suites {
		totalTests += suite.Totals.Tests
		failedTests += suite.Totals.Failed + suite.Totals.Error
		skippedTests += suite.Totals.Skipped

		for _, test := range suite.Tests {
			testIdCounter++
			// Look up file and lineno from our map
			var attrs testcaseAttrs
			if found, ok := testcaseMap[test.Name]; ok {
				attrs = found
			}
			testResult := convertJUnitTestToTestResult(testIdCounter, test, attrs)
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
func convertJUnitTestToTestResult(id int, test junit.Test, attrs testcaseAttrs) types.TestResult {
	// Determine status based on JUnit test result
	var status types.TestStatus
	var errorMessage string
	var backtrace []string
    var failureCause types.FailureCause

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
            // Library provides an error value; parse its string form
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

	// Extract group name (classname) and test case name (name)
	groupName := test.Classname
	testCaseName := test.Name

    // Classify failure cause if failed
    if status == types.StatusFail {
        failureCause = classifyFailure(errorMessage, fullBacktrace)
    }

	// Extract test file path and line number from XML attributes
	testFilePath := attrs.File
	testLineNumber := 0
	if attrs.LineNo != "" {
		if parsedLine, err := strconv.Atoi(attrs.LineNo); err == nil {
			testLineNumber = parsedLine
		}
	}

	// Fallback to Properties if attributes not found in direct XML parse
	if testFilePath == "" && test.Properties != nil {
		if file, ok := test.Properties["file"]; ok {
			testFilePath = file
		}
		if lineno, ok := test.Properties["lineno"]; ok && testLineNumber == 0 {
			if parsedLine, err := strconv.Atoi(lineno); err == nil {
				testLineNumber = parsedLine
			}
		}
	}

    return types.TestResult{
		Id:                id,
		GroupName:         groupName,
		TestCaseName:      testCaseName,
		Status:            status,
		ErrorMessage:      errorMessage,
		TestFilePath:      testFilePath,
		TestLineNumber:    testLineNumber,
        FailureCause:      failureCause,
		FullBacktrace:     fullBacktrace,
		FilteredBacktrace: make([]types.StackFrame, 0), // Will be populated by normalizer
	}
}

// classifyFailure decides the FailureCause from parsed failure fields using simple heuristics
// Inputs: error message and parsed stack frames (project/app and test paths if present)
// 1) Assertion-like messages -> AssertionFailure
// 2) If no frames or top frames point to test/spec -> TestDefinitionError
// 3) Otherwise -> ProductionCodeError
func classifyFailure(message string, frames []types.StackFrame) types.FailureCause {
    m := strings.ToLower(message)
    if m != "" {
        if strings.Contains(m, "assertionerror") || strings.Contains(m, "expected ") || strings.HasPrefix(m, "expected:") || strings.Contains(m, "expected:") {
            return types.FailureCauseAssertion
        }
    }

    // If we have no frames, treat as test definition error (runner/setup/teardown/unmapped)
    if len(frames) == 0 {
        return types.FailureCauseTestDefinition
    }

    // Check if any of the first few frames clearly reference test code paths
    limit := len(frames)
    if limit > 5 {
        limit = 5
    }
    for i := 0; i < limit; i++ {
        f := frames[i]
        lp := strings.ToLower(f.File)
        for _, ind := range testPathIndicators {
            if strings.Contains(lp, ind) || strings.HasSuffix(lp, ind) {
                return types.FailureCauseTestDefinition
            }
        }
        // Common framework indicators
        for _, ind := range frameworkPathIndicators {
            if strings.Contains(lp, ind) {
                return types.FailureCauseTestDefinition
            }
        }
    }

    return types.FailureCauseProductionCode
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
            continue
        }

        // Fallback: detect embedded file:line anywhere in the line, including bracketed forms
        trimmed := strings.TrimSpace(line)
        if embedded := extractFileLine(trimmed); embedded != "" {
            stacktrace = append(stacktrace, embedded)
        }
	}

	return stacktrace
}

// extractFileLine attempts to find a substring like "path/to/file.rb:123" inside an arbitrary line
// and returns it; returns empty string if not found.
func extractFileLine(s string) string {
    // Quick exits for common non-frame lines
    if strings.HasPrefix(s, "Expected:") || strings.HasPrefix(s, "Actual:") {
        return ""
    }

    // Find a ruby/python/js file marker
    idx := strings.Index(s, ".rb:")
    ext := ".rb:"
    if idx == -1 {
        idx = strings.Index(s, ".py:")
        ext = ".py:"
    }
    if idx == -1 {
        idx = strings.Index(s, ".js:")
        ext = ".js:"
    }
    if idx == -1 {
        return ""
    }

    // Walk backwards to find start of path (stop at whitespace or '[')
    start := idx
    for start > 0 {
        c := s[start-1]
        if c == ' ' || c == '\t' || c == '[' || c == '(' {
            break
        }
        start--
    }

    // Walk forwards from after the colon to consume digits of the line number
    end := idx + len(ext)
    for end < len(s) && s[end] >= '0' && s[end] <= '9' {
        end++
    }

    candidate := s[start:end]
    // Basic sanity: ensure there's at least one slash and a colon digit
    if strings.Contains(candidate, "/") && strings.Count(candidate, ":") >= 1 {
        return candidate
    }

    return ""
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

// ParseYAMLFile parses a YAML test results file
func ParseYAMLFile(filePath string) (*ParseResult, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return ParseYAML(data)
}

// ParseYAML parses YAML test results data
func ParseYAML(data []byte) (*ParseResult, error) {
	var yamlTests []map[string]interface{}
	if err := yaml.Unmarshal(data, &yamlTests); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	result := &ParseResult{
		Summary: TestSummary{},
	}

	testIdCounter := 0
	for _, yamlMap := range yamlTests {
		testIdCounter++
		testResult := convertYAMLMapToTestResult(testIdCounter, yamlMap)
		result.Tests = append(result.Tests, testResult)

		// Update summary counts
		result.Summary.Total++
		switch testResult.Status {
		case types.StatusPass:
			result.Summary.Passed++
		case types.StatusFail:
			result.Summary.Failed++
		case types.StatusSkip:
			result.Summary.Skipped++
		}
	}

	return result, nil
}

// extractString safely extracts a string value from a map
func extractString(m map[string]interface{}, key string) string {
	val, ok := m[key]
	if !ok {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

// extractInt safely extracts an int value from a map (handles both int and string types)
func extractInt(m map[string]interface{}, key string) int {
	val, ok := m[key]
	if !ok {
		return 0
	}
	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case string:
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
		return 0
	default:
		return 0
	}
}

// extractStringSlice safely extracts a []string value from a map
func extractStringSlice(m map[string]interface{}, key string) []string {
	val, ok := m[key]
	if !ok {
		return []string{}
	}
	slice, ok := val.([]interface{})
	if !ok {
		return []string{}
	}
	result := make([]string, 0, len(slice))
	for _, item := range slice {
		if str, ok := item.(string); ok {
			result = append(result, str)
		}
	}
	return result
}

// convertYAMLMapToTestResult converts a YAML map to our domain type
func convertYAMLMapToTestResult(id int, yamlMap map[string]interface{}) types.TestResult {
	// Extract basic fields
	groupName := extractString(yamlMap, "test_group_name")
	testCaseName := extractString(yamlMap, "test_case_name")
	testFilePath := extractString(yamlMap, "test_file_path")
	testLineNumber := extractInt(yamlMap, "test_line_number")

	// Parse status
	statusStr := extractString(yamlMap, "test_status")
	var status types.TestStatus
	switch statusStr {
	case "passed":
		status = types.StatusPass
	case "failed":
		status = types.StatusFail
	case "skipped":
		status = types.StatusSkip
	default:
		status = types.StatusFail
	}

	// Parse duration
	durationStr := extractString(yamlMap, "duration")
	duration := 0.0
	if durationStr != "" {
		if parsed, err := strconv.ParseFloat(durationStr, 64); err == nil {
			duration = parsed
		}
	}

	// Parse failure cause
	failureCauseStr := extractString(yamlMap, "failure_cause")
	var failureCause types.FailureCause
	switch failureCauseStr {
	case "error":
		failureCause = types.FailureCauseProductionCode
	case "failed_assertion":
		failureCause = types.FailureCauseAssertion
	default:
		// Empty for passed/skipped tests
		failureCause = ""
	}

	// Extract error fields
	errorMessage := extractString(yamlMap, "error_message")
	errorFilePath := extractString(yamlMap, "error_file_path")
	errorLineNumber := extractInt(yamlMap, "error_line_number")

	// Extract assertion fields
	failedAssertionMessage := extractString(yamlMap, "failed_assertion_details")
	assertionFilePath := extractString(yamlMap, "assertion_file_path")
	assertionLineNumber := extractInt(yamlMap, "assertion_line_number")

	// Parse backtrace
	backtraceStrings := extractStringSlice(yamlMap, "full_backtrace")
	var fullBacktrace []types.StackFrame
	for _, frameStr := range backtraceStrings {
		frame := parseStackFrame(frameStr)
		// Only add frames that have a valid file:line format (check for colon separator)
		if frame.File != "" && strings.Contains(frameStr, ":") {
			fullBacktrace = append(fullBacktrace, frame)
		}
	}

	// Cap at 50 frames
	if len(fullBacktrace) > 50 {
		fullBacktrace = fullBacktrace[:50]
	}

	return types.TestResult{
		Id:                        id,
		GroupName:                 groupName,
		TestCaseName:              testCaseName,
		Status:                    status,
		ErrorMessage:              errorMessage,
		ErrorFilePath:             errorFilePath,
		ErrorLineNumber:           errorLineNumber,
		FailedAssertionMessage:    failedAssertionMessage,
		FailedAssertionFilePath:   assertionFilePath,
		FailedAssertionLineNumber: assertionLineNumber,
		TestFilePath:              testFilePath,
		TestLineNumber:            testLineNumber,
		FailureCause:              failureCause,
		FullBacktrace:             fullBacktrace,
		FilteredBacktrace:         make([]types.StackFrame, 0),
		Duration:                   duration,
	}
}
