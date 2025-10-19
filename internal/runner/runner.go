package runner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
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
	// Parse the test command template
	cmdTemplate, err := template.New("testCommand").Parse(r.config.TestCommand)
	if err != nil {
		return "", fmt.Errorf("failed to parse test command template: %w", err)
	}

	// Prepare template data for interpolation
	templateData := struct {
		Paths string // For now, empty - will be used for specific test paths in future
	}{
		Paths: "", // Empty by default to run all tests
	}

	// Execute template to get the final command
	var cmdBuilder strings.Builder
	if err := cmdTemplate.Execute(&cmdBuilder, templateData); err != nil {
		return "", fmt.Errorf("failed to execute test command template: %w", err)
	}

	finalCommand := cmdBuilder.String()

	// Split the command into parts
	parts := strings.Fields(finalCommand)
	if len(parts) == 0 {
		return "", fmt.Errorf("no test command configured")
	}

	// Create command
	cmd := exec.Command(parts[0], parts[1:]...)

	// Set working directory to project path
	if r.config.ProjectPath != "" {
		cmd.Dir = r.config.ProjectPath
	} else {
		// Fallback to current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current working directory: %w", err)
		}
		cmd.Dir = cwd
	}

	// Execute command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if it's a command not found error
		if execErr, ok := err.(*exec.Error); ok && execErr.Err == exec.ErrNotFound {
			return "", fmt.Errorf("test command not found: %s (make sure it's installed and in PATH)", parts[0])
		}
		
		// Check if it's a permission error
		if execErr, ok := err.(*exec.Error); ok && execErr.Err == os.ErrPermission {
			return "", fmt.Errorf("permission denied running test command: %s", finalCommand)
		}
		
		// Generic test command failure with output
		return "", fmt.Errorf("test command failed (exit code %d): %w\nCommand: %s\nOutput: %s", 
			cmd.ProcessState.ExitCode(), err, finalCommand, string(output))
	}

	return string(output), nil
}

// parseTestOutput parses the XML output from the test command
func (r *TestRunner) parseTestOutput(output string) ([]types.TestResult, error) {
	// Check if output is empty
	if strings.TrimSpace(output) == "" {
		return nil, fmt.Errorf("test command produced no output - check if tests are configured correctly")
	}

	// Parse the output using the existing parser
	result, err := parser.ParseXML([]byte(output))
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML output: %w\nOutput preview: %s", err, truncateString(output, 200))
	}

	// Check if we got any test results
	if len(result.Tests) == 0 {
		return nil, fmt.Errorf("no test results found in XML output - check if test framework is generating JUnit XML correctly")
	}

	return result.Tests, nil
}

// truncateString truncates a string to the specified length and adds ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
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
		return fmt.Errorf("test_command must be specified either via CLI option --test-command or in config file")
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
