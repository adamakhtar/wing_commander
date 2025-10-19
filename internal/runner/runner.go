package runner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/adamakhtar/wing_commander/internal/config"
	"github.com/adamakhtar/wing_commander/internal/grouper"
	"github.com/adamakhtar/wing_commander/internal/parser"
	"github.com/adamakhtar/wing_commander/internal/types"
)

// TestRunner handles execution of test commands and parsing of results
type TestRunner struct {
	config *config.Config
}

// NewTestRunner creates a new TestRunner with the given configuration
func NewTestRunner(cfg *config.Config) *TestRunner {
	return &TestRunner{
		config: cfg,
	}
}

// ExecuteTests runs the configured test command and returns parsed results
func (r *TestRunner) ExecuteTests() (*TestExecutionResult, error) {
	// Execute the test command
	output, err := r.executeTestCommand()
	if err != nil {
		return nil, fmt.Errorf("failed to execute test command: %w", err)
	}

	// Parse the XML output
	testResults, err := r.parseTestOutput(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse test output: %w", err)
	}

	// Normalize backtraces
	normalizer := grouper.NewNormalizer(r.config)
	normalizedResults := normalizer.NormalizeTestResults(testResults)

	// Group failures with change detection
	strategy := grouper.NewErrorLocationStrategy()
	grouperInstance := grouper.NewGrouper(strategy)
	failureGroups := grouperInstance.GroupFailures(normalizedResults)

	return &TestExecutionResult{
		TestResults:    normalizedResults,
		FailureGroups:  failureGroups,
		ExecutionTime:  time.Now(),
		CommandOutput:  output,
	}, nil
}

// executeTestCommand runs the configured test command and returns the output
func (r *TestRunner) executeTestCommand() (string, error) {
	// Split the command into parts
	parts := strings.Fields(r.config.TestCommand)
	if len(parts) == 0 {
		return "", fmt.Errorf("no test command configured")
	}

	// Create command
	cmd := exec.Command(parts[0], parts[1:]...)

	// Set working directory to current directory
	cmd.Dir = "."

	// Execute command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("test command failed: %w\nOutput: %s", err, string(output))
	}

	return string(output), nil
}

// parseTestOutput parses the XML output from the test command
func (r *TestRunner) parseTestOutput(output string) ([]types.TestResult, error) {
	// Parse the output using the existing parser
	result, err := parser.ParseXML([]byte(output))
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML output: %w", err)
	}

	return result.Tests, nil
}

// TestExecutionResult represents the complete result of a test execution
type TestExecutionResult struct {
	TestResults    []types.TestResult  // All test results (normalized)
	FailureGroups  []types.FailureGroup // Grouped failures with change detection
	ExecutionTime  time.Time           // When the tests were executed
	CommandOutput  string              // Raw output from test command
}

// GetSummary returns a summary of the test execution
func (r *TestExecutionResult) GetSummary() TestSummary {
	totalTests := len(r.TestResults)
	failedTests := 0
	passedTests := 0
	skippedTests := 0

	for _, result := range r.TestResults {
		switch result.Status {
		case types.StatusFail:
			failedTests++
		case types.StatusPass:
			passedTests++
		case types.StatusSkip:
			skippedTests++
		}
	}

	return TestSummary{
		TotalTests:   totalTests,
		FailedTests:  failedTests,
		PassedTests:  passedTests,
		SkippedTests: skippedTests,
		FailureGroups: len(r.FailureGroups),
	}
}

// TestSummary provides a high-level summary of test results
type TestSummary struct {
	TotalTests     int
	FailedTests    int
	PassedTests    int
	SkippedTests   int
	FailureGroups  int
}

// ValidateConfig checks if the test configuration is valid
func (r *TestRunner) ValidateConfig() error {
	if r.config.TestFramework == "" {
		return fmt.Errorf("test_framework not specified in config")
	}

	if r.config.TestCommand == "" {
		return fmt.Errorf("test_command not specified in config")
	}

	// Check if test framework is supported
	supportedFrameworks := []config.TestFramework{
		config.FrameworkRSpec,
		config.FrameworkMinitest,
		config.FrameworkPytest,
		config.FrameworkJest,
	}
	found := false
	for _, framework := range supportedFrameworks {
		if r.config.TestFramework == framework {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("unsupported test framework: %s (supported: %s)",
			r.config.TestFramework, strings.Join([]string{"rspec", "minitest", "pytest", "jest"}, ", "))
	}

	return nil
}

// GetWorkingDirectory returns the current working directory
func (r *TestRunner) GetWorkingDirectory() (string, error) {
	return os.Getwd()
}

// CheckTestCommandExists verifies that the test command can be found
func (r *TestRunner) CheckTestCommandExists() error {
	parts := strings.Fields(r.config.TestCommand)
	if len(parts) == 0 {
		return fmt.Errorf("no test command configured")
	}

	// Check if the command exists in PATH
	_, err := exec.LookPath(parts[0])
	if err != nil {
		return fmt.Errorf("test command not found in PATH: %s", parts[0])
	}

	return nil
}
