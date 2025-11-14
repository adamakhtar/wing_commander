package types

// StackFrame represents a single frame in a backtrace
type StackFrame struct {
	File            AbsPath
	Line            int
	Function        string
	ChangeIntensity int
	ChangeReason    string
}

// TestStatus represents the status of a test
type TestStatus string

const (
	StatusPass TestStatus = "pass"
	StatusFail TestStatus = "fail"
	StatusSkip TestStatus = "skip"
)

func (ts TestStatus) Abbreviated() string {
	switch ts {
	case StatusPass:
		return "P"
	case StatusFail:
		return "F"
	case StatusSkip:
		return "S"
	default:
		return "?"
	}
}

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
	Id                int
	GroupName         string
	TestCaseName      string
	Status            TestStatus
	FailureCause      FailureCause
	FailureDetails    string
	FailureFilePath   AbsPath
	FailureLineNumber int
	TestFilePath      AbsPath
	TestLineNumber    int
	FullBacktrace     []StackFrame
	FilteredBacktrace []StackFrame
	Duration          float64
}

func (tr *TestResult) AbbreviatedResult() string {
	if tr.IsFailed() {
		return tr.FailureCause.Abbreviated()
	} else {
		return tr.Status.Abbreviated()
	}
}

// IsFailed reports whether the test result represents a failure.
func (tr *TestResult) IsFailed() bool {
	return tr.Status == StatusFail
}

// IsPassed reports whether the test result represents a passing test.
func (tr *TestResult) IsPassed() bool {
	return tr.Status == StatusPass
}

// IsSkipped reports whether the test result represents a skipped test.
func (tr *TestResult) IsSkipped() bool {
	return tr.Status == StatusSkip
}

// NewStackFrame creates a new StackFrame
func NewStackFrame(file AbsPath, line int, function string) StackFrame {
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
		GroupName:         groupName,
		TestCaseName:      testCaseName,
		Status:            status,
		FullBacktrace:     []StackFrame{},
		FilteredBacktrace: []StackFrame{},
	}
}
