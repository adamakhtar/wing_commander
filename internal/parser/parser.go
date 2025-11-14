package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/adamakhtar/wing_commander/internal/backtrace"
	"github.com/adamakhtar/wing_commander/internal/projectfs"
	"github.com/adamakhtar/wing_commander/internal/testresult"
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
	absPath, err := parseFilePath(path)
	if err == nil {
		return c.matchesTestFileFromAbs(absPath)
	}

	return false
}

func (c *parseContext) matchesTestFileFromAbs(absPath types.AbsPath) bool {
	if c == nil || c.testFileMatcher == nil || absPath == "" {
		return false
	}

	// Try matching absolute path first
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

	return false
}

// ParseResult contains parsed test results and metadata.
type ParseResult struct {
	Tests   []testresult.TestResult
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
		case testresult.StatusPass:
			result.Summary.Passed++
		case testresult.StatusFail:
			result.Summary.Failed++
		case testresult.StatusSkip:
			result.Summary.Skipped++
		}
	}

	return result, nil
}

func firstFrameWithFile(frames []types.StackFrame) *types.StackFrame {
	for _, frame := range frames {
		if frame.FilePath.String() != "" {
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
func classifyFailure(message string, topFrame *types.StackFrame, ctx *parseContext) testresult.FailureCause {
	m := strings.ToLower(message)
	if m != "" {
		if strings.Contains(m, "assertionerror") || strings.Contains(m, "expected ") || strings.HasPrefix(m, "expected:") || strings.Contains(m, "expected:") {
			return testresult.FailureCauseAssertion
		}
	}

	// If we have no frames, treat as test definition error (runner/setup/teardown/unmapped)
	if topFrame == nil || topFrame.FilePath.String() == "" {
		return testresult.FailureCauseTestDefinition
	}

	if ctx != nil && ctx.matchesTestFileFromAbs(topFrame.FilePath) {
		return testresult.FailureCauseTestDefinition
	}

	lp := strings.ToLower(topFrame.FilePath.String())
	for _, ind := range testPathIndicators {
		if strings.Contains(lp, ind) || strings.HasSuffix(lp, ind) {
			return testresult.FailureCauseTestDefinition
		}
	}
	for _, ind := range frameworkPathIndicators {
		if strings.Contains(lp, ind) {
			return testresult.FailureCauseTestDefinition
		}
	}

	return testresult.FailureCauseProductionCode
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

// parseFilePath converts a file path string to AbsPath.
// If the path is absolute, it uses NewAbsPath directly.
// If the path is relative, it assumes it's relative to ProjectFS root and converts it using ProjectFS.Abs().
func parseFilePath(path string) (types.AbsPath, error) {
	if path == "" {
		return types.AbsPath(""), fmt.Errorf("path cannot be empty")
	}

	cleaned := filepath.Clean(path)

	if filepath.IsAbs(cleaned) {
		return types.NewAbsPath(cleaned)
	}

	// Relative path - convert using ProjectFS
	fs := projectfs.GetProjectFS()
	relPath, err := types.NewRelPath(cleaned)
	if err != nil {
		return types.AbsPath(""), fmt.Errorf("failed to create RelPath: %w", err)
	}

	return fs.Abs(relPath), nil
}

// convertMapToTestResult converts a summary map into a TestResult.
func convertMapToTestResult(id int, summary map[string]interface{}, ctx *parseContext) testresult.TestResult {
	// Extract basic fields
	groupName := extractString(summary, "test_group_name")
	testCaseName := extractString(summary, "test_case_name")
	testFilePath := extractString(summary, "test_file_path")
	testLineNumber := extractInt(summary, "test_line_number")

	// Parse status
	statusStr := extractString(summary, "test_status")
	var status testresult.TestStatus
	switch statusStr {
	case "passed":
		status = testresult.StatusPass
	case "failed":
		status = testresult.StatusFail
	case "skipped":
		status = testresult.StatusSkip
	default:
		status = testresult.StatusFail
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
	fullBacktrace := backtrace.NewBacktrace()
	frameCount := 0
	for _, frameStr := range backtraceStrings {
		// Cap at 50 frames
		if frameCount >= 50 {
			break
		}
		// Only add frames that have a valid file:line format (check for colon separator)
		if strings.Contains(frameStr, ":") {
			fullBacktrace.Append(frameStr)
			frameCount++
		}
	}

	var failureCause testresult.FailureCause
	if status == testresult.StatusFail {
		topFrame := firstFrameWithFile(fullBacktrace.AllStackFrames())
		if topFrame == nil && failureFilePath != "" {
			absPath, err := parseFilePath(failureFilePath)
			if err == nil {
				topFrame = &types.StackFrame{
					FilePath: absPath,
					Line:     failureLineNumber,
				}
			}
		}
		failureCause = classifyFailure(failureDetails, topFrame, ctx)
	}

	// Convert file paths to AbsPath
	var failureFilePathAbs types.AbsPath
	if failureFilePath != "" {
		if abs, err := parseFilePath(failureFilePath); err == nil {
			failureFilePathAbs = abs
		}
	}

	var testFilePathAbs types.AbsPath
	if testFilePath != "" {
		if abs, err := parseFilePath(testFilePath); err == nil {
			testFilePathAbs = abs
		}
	}

	return testresult.TestResult{
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
		FilteredBacktrace: backtrace.NewBacktrace(),
		Duration:          duration,
	}
}
