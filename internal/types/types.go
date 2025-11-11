package types

import (
	"path/filepath"
	"time"
)

// StackFrame represents a single frame in a backtrace
type StackFrame struct {
	File            string // File path
	Line            int    // Line number
	Function        string // Function/method name (may be empty)
	ChangeIntensity int    // 0-3 (0 = no highlight, 3 = strongest)
	ChangeReason    string // "uncommitted" | "last_commit" | "previous_commit"
}

func (s StackFrame) RelativeFilePath(absolutePath string) string {
	if filepath.IsAbs(s.File) {
		rel, err := filepath.Rel(absolutePath, s.File)
		if err != nil {
			return s.File
		}
		return rel
	}
	return s.File
}


// TestStatus represents the status of a test
type TestStatus string

const (
	StatusPass TestStatus = "pass"
	StatusFail TestStatus = "fail"
	StatusSkip TestStatus = "skip"
)

// FailureCause represents the coarse-grained cause of a test failure
type FailureCause string

func (fc FailureCause) Abbreviated() string {
	switch fc {
	case FailureCauseTestDefinition:
		return "T"
	case FailureCauseProductionCode:
		return "C"
	case FailureCauseAssertion:
		return "A"
	default:
		return ""
	}
}

func (fc FailureCause) String() string {
	switch fc {
	case FailureCauseTestDefinition:
		return "Test Definition Error"
	case FailureCauseProductionCode:
		return "Production Code Error"
	case FailureCauseAssertion:
		return "Assertion Failure"
	default:
		return ""
	}
}

const (
    // FailureCauseTestDefinition indicates the failure originated from test code/framework
    FailureCauseTestDefinition FailureCause = "test_definition_error"
    // FailureCauseProductionCode indicates the failure was raised by the application under test
    FailureCauseProductionCode FailureCause = "production_code_error"
    // FailureCauseAssertion indicates the test completed but an expectation failed
    FailureCauseAssertion FailureCause = "assertion_failure"
)

// TestResult represents a single test execution result
type TestResult struct {
	Id               int           // Unique ID for the test result
	GroupName        string        // Test group name
	TestCaseName     string        // Test case name
	Status           TestStatus    // Test status
	FailureCause     FailureCause  // Cause of the failure, derived from failure details/backtrace
	FailureDetails   string        // Human-readable description of the failure
	FailureFilePath  string        // File path where the failure originated
	FailureLineNumber int          // Line number where the failure originated
	TestFilePath     string        // File path of the test
	TestLineNumber   int           // Line number of the test definition
	FullBacktrace    []StackFrame  // Complete backtrace (up to 50 frames)
	FilteredBacktrace []StackFrame // Filtered backtrace (project frames only)
	Duration         float64       // Duration of the test in seconds
}

// FailureGroup represents a group of tests that failed with similar backtraces
type FailureGroup struct {
	Hash              string        // Hash ID based on normalized backtrace
	ErrorMessage      string        // Representative error message
	NormalizedBacktrace []StackFrame // Normalized backtrace signature
	Tests             []TestResult  // All tests in this group
	Count             int           // Number of failed tests in group
}

// CacheData represents the structure saved to cache file
type CacheData struct {
	Timestamp      time.Time       `json:"timestamp"`
	FailureGroups  []FailureGroup  `json:"failure_groups"`
	RawTestResults []TestResult    `json:"raw_test_results"`
}

// NewStackFrame creates a new StackFrame
func NewStackFrame(file string, line int, function string) StackFrame {
	return StackFrame{
		File:            file,
		Line:            line,
		Function:        function,
		ChangeIntensity: 0,
		ChangeReason:    "",
	}
}

// NewTestResult creates a new TestResult
func NewTestResult(groupName string, testCaseName string, status TestStatus) TestResult {
	return TestResult{
		GroupName:          groupName,
		TestCaseName:       testCaseName,
		Status:             status,
		FullBacktrace:      []StackFrame{},
		FilteredBacktrace:  []StackFrame{},
	}
}

// NewFailureGroup creates a new FailureGroup
func NewFailureGroup(hash string, errorMessage string) FailureGroup {
	return FailureGroup{
		Hash:                hash,
		ErrorMessage:        errorMessage,
		NormalizedBacktrace: []StackFrame{},
		Tests:               []TestResult{},
		Count:               0,
	}
}
