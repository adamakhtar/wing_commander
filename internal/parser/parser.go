package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/adamakhtar/wing_commander/internal/projectfs"
	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/gobwas/glob"
	"gopkg.in/yaml.v3"
)

// testPathIndicators lists substrings that indicate frames belonging to test code paths.
// Limited to Ruby frameworks for now (RSpec and Minitest). Extend as more frameworks are supported.
var testPathIndicators = []string{
	"/spec/",
	"/test/",
	"_spec.rb",
}

// frameworkPathIndicators lists substrings that indicate frames from test frameworks/runners.
var frameworkPathIndicators = []string{
	"/rspec/",
	"/minitest/",
	"/jest/",
	"/node_modules/",
	"/gems/",
}

// ParseOptions controls optional parsing behaviour.
type ParseOptions struct {
	TestFilePattern string
}

type parseContext struct {
	testFileMatcher glob.Glob
}

func newParseContext(opts *ParseOptions) (*parseContext, error) {
	if opts == nil {
		return &parseContext{}, nil
	}

	ctx := &parseContext{}

	if opts.TestFilePattern != "" {
		pattern := filepath.ToSlash(opts.TestFilePattern)
		compiled, err := glob.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to compile test file pattern %q: %w", opts.TestFilePattern, err)
		}
		ctx.testFileMatcher = compiled
	}

	return ctx, nil
}

func (c *parseContext) matchesTestFile(path string) bool {
	if c == nil || c.testFileMatcher == nil || path == "" {
		return false
	}

	// Try matching as absolute path first
	absPath, err := types.NewAbsPath(path)
	if err == nil {
		if absPath.MatchGlob(c.testFileMatcher) {
			return true
		}

		// Use ProjectFS to convert absolute path to relative for pattern matching
		fs := projectfs.GetProjectFS()
		rel, err := fs.Rel(absPath)
		if err == nil {
			if rel.MatchGlob(c.testFileMatcher) {
				return true
			}
		}
	}

	return false
}

// ParseResult contains parsed test results and metadata.
type ParseResult struct {
	Tests   []types.TestResult
	Summary TestSummary
}

// TestSummary contains test run statistics.
type TestSummary struct {
	Total   int
	Passed  int
	Failed  int
	Skipped int
}

// ParseFile parses a Wing Commander reporter summary file.
func ParseFile(filePath string, opts *ParseOptions) (*ParseResult, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return Parse(data, opts)
}

// Parse parses Wing Commander reporter summary data.
func Parse(data []byte, opts *ParseOptions) (*ParseResult, error) {
	ctx, err := newParseContext(opts)
	if err != nil {
		return nil, err
	}

	var tests []map[string]interface{}
	if err := yaml.Unmarshal(data, &tests); err != nil {
		return nil, fmt.Errorf("failed to parse summary data: %w", err)
	}

	result := &ParseResult{
		Summary: TestSummary{},
	}

	testID := 0
	for _, testMap := range tests {
		testID++
		testResult := convertMapToTestResult(testID, testMap, ctx)
		result.Tests = append(result.Tests, testResult)

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

func firstFrameWithFile(frames []types.StackFrame) *types.StackFrame {
	for _, frame := range frames {
		if frame.File.String() != "" {
			frameCopy := frame
			return &frameCopy
		}
	}
	return nil
}

// classifyFailure decides the FailureCause from parsed failure fields using simple heuristics.
// Inputs: error message and top stack frame (project/app and test paths if present)
// 1) Assertion-like messages -> AssertionFailure
// 2) If no frames or top frames point to test/spec -> TestDefinitionError
// 3) Otherwise -> ProductionCodeError
func classifyFailure(message string, topFrame *types.StackFrame, ctx *parseContext) types.FailureCause {
	m := strings.ToLower(message)
	if m != "" {
		if strings.Contains(m, "assertionerror") || strings.Contains(m, "expected ") || strings.HasPrefix(m, "expected:") || strings.Contains(m, "expected:") {
			return types.FailureCauseAssertion
		}
	}

	// If we have no frames, treat as test definition error (runner/setup/teardown/unmapped)
	if topFrame == nil || topFrame.File.String() == "" {
		return types.FailureCauseTestDefinition
	}

	if ctx != nil && ctx.matchesTestFile(topFrame.File.String()) {
		return types.FailureCauseTestDefinition
	}

	lp := strings.ToLower(topFrame.File.String())
	for _, ind := range testPathIndicators {
		if strings.Contains(lp, ind) || strings.HasSuffix(lp, ind) {
			return types.FailureCauseTestDefinition
		}
	}
	for _, ind := range frameworkPathIndicators {
		if strings.Contains(lp, ind) {
			return types.FailureCauseTestDefinition
		}
	}

	return types.FailureCauseProductionCode
}

// parseStackFrame parses a backtrace frame string into a StackFrame.
func parseStackFrame(frameStr string) types.StackFrame {
	// Common formats:
	// "app/models/user.rb:42:in `create_user'"
	// "app/models/user.rb:42"
	// "File \"app/models/user.rb\", line 42, in create_user"

	// Handle Python format first
	if strings.HasPrefix(frameStr, "File \"") {
		absPath, _ := types.NewAbsPath(frameStr)
		return types.StackFrame{
			File:     absPath,
			Line:     0,
			Function: "",
		}
	}

	parts := strings.Split(frameStr, ":")
	if len(parts) < 2 {
		absPath, _ := types.NewAbsPath(frameStr)
		return types.StackFrame{File: absPath}
	}

	file := parts[0]

	// Try to extract line number
	var line int
	var function string

	if len(parts) >= 2 {
		// Parse line number
		if _, err := fmt.Sscanf(parts[1], "%d", &line); err != nil {
			absPath, _ := types.NewAbsPath(file)
			return types.StackFrame{File: absPath}
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

	absPath, _ := types.NewAbsPath(file)
	return types.StackFrame{
		File:     absPath,
		Line:     line,
		Function: function,
	}
}

// extractString safely extracts a string value from a map.
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

// extractInt safely extracts an int value from a map (handles both int and string types).
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

// extractStringSlice safely extracts a []string value from a map.
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

// convertMapToTestResult converts a summary map into a TestResult.
func convertMapToTestResult(id int, summary map[string]interface{}, ctx *parseContext) types.TestResult {
	// Extract basic fields
	groupName := extractString(summary, "test_group_name")
	testCaseName := extractString(summary, "test_case_name")
	testFilePath := extractString(summary, "test_file_path")
	testLineNumber := extractInt(summary, "test_line_number")

	// Parse status
	statusStr := extractString(summary, "test_status")
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
	durationStr := extractString(summary, "duration")
	duration := 0.0
	if durationStr != "" {
		if parsed, err := strconv.ParseFloat(durationStr, 64); err == nil {
			duration = parsed
		}
	}

	// Extract unified failure fields (new schema)
	failureDetails := extractString(summary, "failure_details")
	failureFilePath := extractString(summary, "failure_file_path")
	failureLineNumber := extractInt(summary, "failure_line_number")


	// Parse backtrace
	backtraceStrings := extractStringSlice(summary, "full_backtrace")
	var fullBacktrace []types.StackFrame
	for _, frameStr := range backtraceStrings {
		frame := parseStackFrame(frameStr)
		// Only add frames that have a valid file:line format (check for colon separator)
		if frame.File.String() != "" && strings.Contains(frameStr, ":") {
			fullBacktrace = append(fullBacktrace, frame)
		}
	}

	// Cap at 50 frames
	if len(fullBacktrace) > 50 {
		fullBacktrace = fullBacktrace[:50]
	}

	var failureCause types.FailureCause
	if status == types.StatusFail {
		topFrame := firstFrameWithFile(fullBacktrace)
		if topFrame == nil && failureFilePath != "" {
			absPath, err := types.NewAbsPath(failureFilePath)
			if err == nil {
				topFrame = &types.StackFrame{
					File: absPath,
					Line: failureLineNumber,
				}
			}
		}
		failureCause = classifyFailure(failureDetails, topFrame, ctx)
	}

	// Convert file paths to AbsPath
	var failureFilePathAbs types.AbsPath
	if failureFilePath != "" {
		if abs, err := types.NewAbsPath(failureFilePath); err == nil {
			failureFilePathAbs = abs
		}
	}

	var testFilePathAbs types.AbsPath
	if testFilePath != "" {
		if abs, err := types.NewAbsPath(testFilePath); err == nil {
			testFilePathAbs = abs
		}
	}

	return types.TestResult{
		Id:                id,
		GroupName:         groupName,
		TestCaseName:      testCaseName,
		Status:            status,
		FailureCause:      failureCause,
		FailureDetails:    failureDetails,
		FailureFilePath:   failureFilePathAbs,
		FailureLineNumber: failureLineNumber,
		TestFilePath:      testFilePathAbs,
		TestLineNumber:    testLineNumber,
		FullBacktrace:     fullBacktrace,
		FilteredBacktrace: make([]types.StackFrame, 0),
		Duration:          duration,
	}
}
