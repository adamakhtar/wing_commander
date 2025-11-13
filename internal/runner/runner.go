package runner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/adamakhtar/wing_commander/internal/backtrace"
	"github.com/adamakhtar/wing_commander/internal/config"
	"github.com/adamakhtar/wing_commander/internal/parser"
	"github.com/adamakhtar/wing_commander/internal/testrun"
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
func (r *TestRunner) ExecuteTests(testRun testrun.TestRun) (*TestExecutionResult, error) {
	// Execute the test command
	output, err := r.executeTestCommand(testRun.Filepaths)
	if err != nil {
		return nil, fmt.Errorf("failed to execute test command: %w", err)
	}

	var testResults []types.TestResult

	// Parse YAML summary file
	parseOpts := &parser.ParseOptions{
		ProjectPath:     r.config.ProjectPath,
		TestFilePattern: r.config.TestFilePattern,
	}

	parsed, err := parser.ParseFile(r.config.TestResultsPath, parseOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to parse test output summary: %w", err)
	}
	testResults = parsed.Tests

	// Normalize backtraces
	normalizer := backtrace.NewNormalizer(r.config)
	normalizedResults := normalizer.NormalizeTestResults(testResults)

	// Partition results by status (no backtrace grouping)
	var passedTests []types.TestResult
	var failedTests []types.TestResult
	var skippedTests []types.TestResult

	for _, tr := range normalizedResults {
		switch tr.Status {
		case types.StatusPass:
			passedTests = append(passedTests, tr)
		case types.StatusFail:
			failedTests = append(failedTests, tr)
		case types.StatusSkip:
			skippedTests = append(skippedTests, tr)
		}
	}

	metrics := calculateMetrics(normalizedResults)

	return &TestExecutionResult{
		TestRunId:     testRun.Id,
		TestResults:   normalizedResults,
		Metrics:       metrics,
		PassedTests:   passedTests,
		FailedTests:   failedTests,
		SkippedTests:  skippedTests,
		ExecutionTime: time.Now(),
		CommandOutput: output,
	}, nil
}

// executeTestCommand runs the configured test command and returns the output
func (r *TestRunner) executeTestCommand(filepaths []string) (string, error) {
	// Build the full command string
	commandStr := r.config.TestCommand
	if len(filepaths) > 0 {
		commandStr = commandStr + " " + strings.Join(filepaths, " ")
	}

	// Execute via shell to handle multi-word commands like "bundle exec rake test"
	cmd := exec.Command("sh", "-c", commandStr)

	// Set working directory to project path
	cmd.Dir = r.config.ProjectPath
	// Execute command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if it's a command not found error (when sh itself is not found)
		if execErr, ok := err.(*exec.Error); ok && execErr.Err == exec.ErrNotFound {
			return "", fmt.Errorf("test command not found: %s (make sure it's installed and in PATH)", commandStr)
		}

		// Check if it's a permission error
		if execErr, ok := err.(*exec.Error); ok && execErr.Err == os.ErrPermission {
			return "", fmt.Errorf("permission denied running test command: %s", commandStr)
		}

		// Check if it's an exit error (command ran but returned non-zero)
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode := -1
			if exitErr.ProcessState != nil {
				exitCode = exitErr.ExitCode()
			}
			// Exit code 1 typically means tests failed, which is normal - we still have valid output
			// Only return error if exit code suggests a real problem (not 1)
			if exitCode != 1 {
				return "", fmt.Errorf("test command failed (exit code %d): %w\nCommand: %s\nOutput: %s",
					exitCode, err, commandStr, string(output))
			}
			// Exit code 1 (test failures) - return the output so it can be parsed
			return string(output), nil
		}

		// Unknown error type
		return "", fmt.Errorf("unexpected error executing test command: %w\nCommand: %s\nOutput: %s",
			err, commandStr, string(output))
	}

	return string(output), nil
}

// TestExecutionResult represents the complete result of a test execution
type TestExecutionResult struct {
	TestRunId     int                // The ID of the test run that was executed
	TestResults   []types.TestResult // All test results (normalized)
	PassedTests   []types.TestResult // Passed tests
	FailedTests   []types.TestResult // Failed tests
	SkippedTests  []types.TestResult // Skipped tests
	ExecutionTime time.Time          // When the tests were executed
	Metrics       Metrics            // Metrics of the test execution
	CommandOutput string             // Raw output from test command
}

// GetSummary returns a summary of the test execution
func calculateMetrics(testResults []types.TestResult) Metrics {
	totalTests := len(testResults)
	failedTests := 0
	passedTests := 0
	skippedTests := 0

	for _, result := range testResults {
		switch result.Status {
		case types.StatusFail:
			failedTests++
		case types.StatusPass:
			passedTests++
		case types.StatusSkip:
			skippedTests++
		}
	}

	return Metrics{
		TotalTests:   totalTests,
		FailedTests:  failedTests,
		PassedTests:  passedTests,
		SkippedTests: skippedTests,
	}
}

// TestSummary provides a high-level summary of test results
type Metrics struct {
	TotalTests   int
	FailedTests  int
	PassedTests  int
	SkippedTests int
}

// ValidateConfig checks if the test configuration is valid
func (r *TestRunner) ValidateConfig() error {
	if r.config.TestFramework == "" {
		return fmt.Errorf("test_framework not specified in config")
	}

	if r.config.TestCommand == "" {
		return fmt.Errorf("test_command must be specified either via CLI option --test-command or in config file")
	}

	if r.config.TestFramework != config.FrameworkMinitest {
		return fmt.Errorf("unsupported test framework: %s (Wing Commander currently requires WingCommanderReporter YAML output)", r.config.TestFramework)
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
