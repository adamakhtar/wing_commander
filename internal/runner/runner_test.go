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
			errorMsg:    "test_command not specified in config",
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
			{Name: "Test 1", Status: types.StatusPass},
			{Name: "Test 2", Status: types.StatusFail},
			{Name: "Test 3", Status: types.StatusSkip},
			{Name: "Test 4", Status: types.StatusPass},
			{Name: "Test 5", Status: types.StatusFail},
		},
		FailureGroups: []types.FailureGroup{
			{Hash: "group1", Count: 2},
		},
		ExecutionTime: time.Now(),
	}

	summary := result.GetSummary()

	assert.Equal(t, 5, summary.TotalTests)
	assert.Equal(t, 2, summary.FailedTests)
	assert.Equal(t, 2, summary.PassedTests)
	assert.Equal(t, 1, summary.SkippedTests)
	assert.Equal(t, 1, summary.FailureGroups)
}

func TestTestExecutionResult_EmptyResults(t *testing.T) {
	result := &TestExecutionResult{
		TestResults:   []types.TestResult{},
		FailureGroups: []types.FailureGroup{},
		ExecutionTime: time.Now(),
	}

	summary := result.GetSummary()

	assert.Equal(t, 0, summary.TotalTests)
	assert.Equal(t, 0, summary.FailedTests)
	assert.Equal(t, 0, summary.PassedTests)
	assert.Equal(t, 0, summary.SkippedTests)
	assert.Equal(t, 0, summary.FailureGroups)
}

func TestTestRunner_ParseTestOutput_RSpec(t *testing.T) {
	runner := NewTestRunner(&config.Config{
		TestFramework: config.FrameworkRSpec,
	})

	// Mock RSpec JSON output (simplified format that parser expects)
	rspecJSON := `[
		{
			"name": "should pass",
			"status": "passed",
			"message": "",
			"backtrace": [],
			"duration": 0.1,
			"file": "spec/test_spec.rb",
			"line": 5
		},
		{
			"name": "should fail",
			"status": "failed",
			"message": "Expected true to be false",
			"backtrace": [
				"spec/test_spec.rb:10:in 'block (2 levels) in <top (required)>'"
			],
			"duration": 0.2,
			"file": "spec/test_spec.rb",
			"line": 10
		}
	]`

	results, err := runner.parseTestOutput(rspecJSON)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "should pass", results[0].Name)
	assert.Equal(t, types.StatusPass, results[0].Status)
	assert.Equal(t, "should fail", results[1].Name)
	assert.Equal(t, types.StatusFail, results[1].Status)
}

func TestTestRunner_ParseTestOutput_Minitest(t *testing.T) {
	runner := NewTestRunner(&config.Config{
		TestFramework: config.FrameworkMinitest,
	})

	// Mock Minitest JSON output (simplified format that parser expects)
	minitestJSON := `[
		{
			"name": "test_pass",
			"status": "passed",
			"message": "",
			"backtrace": [],
			"duration": 0.1,
			"file": "test/test_class.rb",
			"line": 5
		},
		{
			"name": "test_fail",
			"status": "failed",
			"message": "Expected true to be false",
			"backtrace": [
				"test/test_class.rb:10:in 'test_fail'"
			],
			"duration": 0.2,
			"file": "test/test_class.rb",
			"line": 10
		}
	]`

	results, err := runner.parseTestOutput(minitestJSON)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "test_pass", results[0].Name)
	assert.Equal(t, types.StatusPass, results[0].Status)
	assert.Equal(t, "test_fail", results[1].Name)
	assert.Equal(t, types.StatusFail, results[1].Status)
}

func TestTestRunner_ParseTestOutput_InvalidJSON(t *testing.T) {
	runner := NewTestRunner(&config.Config{
		TestFramework: config.FrameworkRSpec,
	})

	invalidJSON := `{ invalid json }`

	_, err := runner.parseTestOutput(invalidJSON)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON output")
}
