package runner

import (
	"testing"
	"time"

	"github.com/adamakhtar/wing_commander/internal/config"
	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestNewTestRunner(t *testing.T) {
	cfg := &config.Config{
		TestFramework: config.FrameworkRSpec,
		TestCommand:   "bundle exec rspec --format json",
	}

	runner := NewTestRunner(cfg)

	assert.NotNil(t, runner)
	assert.Equal(t, cfg, runner.config)
}

func TestTestRunner_ValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid RSpec config",
			config: &config.Config{
				TestFramework: config.FrameworkRSpec,
				TestCommand:   "bundle exec rspec --format json",
			},
			expectError: false,
		},
		{
			name: "Valid Minitest config",
			config: &config.Config{
				TestFramework: config.FrameworkMinitest,
				TestCommand:   "bundle exec ruby -Itest test/**/*_test.rb",
			},
			expectError: false,
		},
		{
			name: "Missing test framework",
			config: &config.Config{
				TestCommand: "bundle exec rspec --format json",
			},
			expectError: true,
			errorMsg:    "test_framework not specified in config",
		},
		{
			name: "Missing test command",
			config: &config.Config{
				TestFramework: "rspec",
			},
			expectError: true,
			errorMsg:    "test_command must be specified either via CLI option --test-command or in config file",
		},
		{
			name: "Unsupported test framework",
			config: &config.Config{
				TestFramework: config.FrameworkUnknown,
				TestCommand:   "some command",
			},
			expectError: true,
			errorMsg:    "unsupported test framework: unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := NewTestRunner(tt.config)
			err := runner.ValidateConfig()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTestRunner_CheckTestCommandExists(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		expectError bool
	}{
		{
			name: "Valid command (ls should exist)",
			config: &config.Config{
				TestCommand: "ls",
			},
			expectError: false,
		},
		{
			name: "Invalid command",
			config: &config.Config{
				TestCommand: "nonexistentcommand12345",
			},
			expectError: true,
		},
		{
			name: "Empty command",
			config: &config.Config{
				TestCommand: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := NewTestRunner(tt.config)
			err := runner.CheckTestCommandExists()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTestRunner_GetWorkingDirectory(t *testing.T) {
	runner := NewTestRunner(&config.Config{})

	wd, err := runner.GetWorkingDirectory()

	assert.NoError(t, err)
	assert.NotEmpty(t, wd)
}

func TestTestExecutionResult_GetSummary(t *testing.T) {
	result := &TestExecutionResult{
		TestResults: []types.TestResult{
			{GroupName: "Test 1", TestCaseName: "", Status: types.StatusPass},
			{GroupName: "Test 2", TestCaseName: "", Status: types.StatusFail},
			{GroupName: "Test 3", TestCaseName: "", Status: types.StatusSkip},
			{GroupName: "Test 4", TestCaseName: "", Status: types.StatusPass},
			{GroupName: "Test 5", TestCaseName: "", Status: types.StatusFail},
		},
		ExecutionTime: time.Now(),
	}

    summary := calculateMetrics(result.TestResults)

	assert.Equal(t, 5, summary.TotalTests)
	assert.Equal(t, 2, summary.FailedTests)
	assert.Equal(t, 2, summary.PassedTests)
	assert.Equal(t, 1, summary.SkippedTests)
}

func TestTestExecutionResult_EmptyResults(t *testing.T) {
	result := &TestExecutionResult{
		TestResults:   []types.TestResult{},
		ExecutionTime: time.Now(),
	}

    summary := calculateMetrics(result.TestResults)

	assert.Equal(t, 0, summary.TotalTests)
	assert.Equal(t, 0, summary.FailedTests)
	assert.Equal(t, 0, summary.PassedTests)
	assert.Equal(t, 0, summary.SkippedTests)
}

func TestTestRunner_ParseTestOutput_RSpec(t *testing.T) {
	runner := NewTestRunner(&config.Config{
		TestFramework: config.FrameworkRSpec,
	})

	// Mock RSpec XML output
	rspecXML := `<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
  <testsuite name="RSpec" tests="2" failures="1" skipped="0" time="0.3">
    <testcase classname="TestSpec" name="should pass" time="0.1">
    </testcase>
    <testcase classname="TestSpec" name="should fail" time="0.2">
      <failure message="Expected true to be false">
        spec/test_spec.rb:10:in 'block (2 levels) in &lt;top (required)&gt;'
      </failure>
    </testcase>
  </testsuite>
</testsuites>`

	results, err := runner.parseTestOutput(rspecXML)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "TestSpec", results[0].GroupName)
	assert.Equal(t, "should pass", results[0].TestCaseName)
	assert.Equal(t, types.StatusPass, results[0].Status)
	assert.Equal(t, "TestSpec", results[1].GroupName)
	assert.Equal(t, "should fail", results[1].TestCaseName)
	assert.Equal(t, types.StatusFail, results[1].Status)
}

func TestTestRunner_ParseTestOutput_Minitest(t *testing.T) {
	runner := NewTestRunner(&config.Config{
		TestFramework: config.FrameworkMinitest,
	})

	// Mock Minitest XML output
	minitestXML := `<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
  <testsuite name="Minitest" tests="2" failures="1" skipped="0" time="0.3">
    <testcase classname="TestClass" name="test_pass" time="0.1">
    </testcase>
    <testcase classname="TestClass" name="test_fail" time="0.2">
      <failure message="Expected true to be false">
        test/test_class.rb:10:in 'test_fail'
      </failure>
    </testcase>
  </testsuite>
</testsuites>`

	results, err := runner.parseTestOutput(minitestXML)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "TestClass", results[0].GroupName)
	assert.Equal(t, "test_pass", results[0].TestCaseName)
	assert.Equal(t, types.StatusPass, results[0].Status)
	assert.Equal(t, "TestClass", results[1].GroupName)
	assert.Equal(t, "test_fail", results[1].TestCaseName)
	assert.Equal(t, types.StatusFail, results[1].Status)
}

func TestTestRunner_ParseTestOutput_InvalidXML(t *testing.T) {
	runner := NewTestRunner(&config.Config{
		TestFramework: config.FrameworkRSpec,
	})

	invalidXML := `<?xml version="1.0"?><invalid>xml`

	_, err := runner.parseTestOutput(invalidXML)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse XML output")
}
