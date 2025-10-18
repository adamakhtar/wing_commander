package parser

import (
	"testing"

	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "valid RSpec JSON file",
			filename: "../../testdata/fixtures/rspec_failures.json",
			wantErr:  false,
		},
		{
			name:     "valid Minitest JSON file",
			filename: "../../testdata/fixtures/minitest_failures.json",
			wantErr:  false,
		},
		{
			name:     "non-existent file",
			filename: "testdata/fixtures/nonexistent.json",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseFile(tt.filename)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestParseJSON_RSpec(t *testing.T) {
	jsonData := `[
		{
			"name": "User should be valid",
			"status": "failed",
			"message": "Expected User to be valid",
			"backtrace": [
				"app/models/user.rb:42:in 'create_user'",
				"spec/models/user_spec.rb:15:in 'block (2 levels) in <top (required)>'"
			]
		}
	]`

	result, err := ParseJSON([]byte(jsonData))
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Len(t, result.Tests, 1)

	test := result.Tests[0]
	assert.Equal(t, "User should be valid", test.Name)
	assert.Equal(t, types.StatusFail, test.Status)
	assert.Equal(t, "Expected User to be valid", test.ErrorMessage)
	assert.Len(t, test.FullBacktrace, 2)

	// Check first frame
	frame := test.FullBacktrace[0]
	assert.Equal(t, "app/models/user.rb", frame.File)
	assert.Equal(t, 42, frame.Line)
	assert.Equal(t, "create_user", frame.Function)
}

func TestParseJSON_Minitest(t *testing.T) {
	jsonData := `{
		"tests": [
			{
				"name": "test_user_creation",
				"status": "failed",
				"message": "Expected User to be valid",
				"backtrace": [
				"app/models/user.rb:42:in 'create_user'",
				"test/models/user_test.rb:15:in 'test_user_creation'"
				]
			}
		],
		"summary": {
			"total": 1,
			"passed": 0,
			"failed": 1,
			"skipped": 0
		}
	}`

	result, err := ParseJSON([]byte(jsonData))
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 1, result.Summary.Total)
	assert.Equal(t, 1, result.Summary.Failed)
	assert.Len(t, result.Tests, 1)

	test := result.Tests[0]
	assert.Equal(t, "test_user_creation", test.Name)
	assert.Equal(t, types.StatusFail, test.Status)
}

func TestParseStackFrame(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected types.StackFrame
	}{
		{
			name:  "Ruby with method",
			input: "app/models/user.rb:42:in 'create_user'",
			expected: types.StackFrame{
				File:     "app/models/user.rb",
				Line:     42,
				Function: "create_user",
			},
		},
		{
			name:  "Ruby without method",
			input: "app/models/user.rb:42",
			expected: types.StackFrame{
				File:     "app/models/user.rb",
				Line:     42,
				Function: "",
			},
		},
		{
			name:  "Python format",
			input: "File \"app/models/user.py\", line 42, in create_user",
			expected: types.StackFrame{
				File:     "File \"app/models/user.py\", line 42, in create_user",
				Line:     0,
				Function: "",
			},
		},
		{
			name:  "Invalid format",
			input: "invalid_frame",
			expected: types.StackFrame{
				File:     "invalid_frame",
				Line:     0,
				Function: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseStackFrame(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}


func TestConvertToTestResult(t *testing.T) {
	tests := []struct {
		name     string
		input    InputTestResult
		expected types.TestResult
	}{
		{
			name: "Failed test",
			input: InputTestResult{
				Name:      "Test name",
				Status:    "failed",
				Message:   "Error message",
				Backtrace: []string{"app/test.rb:10:in 'test_method'"},
			},
			expected: types.TestResult{
				Name:              "Test name",
				Status:            types.StatusFail,
				ErrorMessage:      "Error message",
				FullBacktrace:     []types.StackFrame{{File: "app/test.rb", Line: 10, Function: "test_method"}},
				FilteredBacktrace: []types.StackFrame{},
			},
		},
		{
			name: "Passed test",
			input: InputTestResult{
				Name:      "Test name",
				Status:    "passed",
				Message:   "",
				Backtrace: []string{},
			},
			expected: types.TestResult{
				Name:              "Test name",
				Status:            types.StatusPass,
				ErrorMessage:      "",
				FullBacktrace:     make([]types.StackFrame, 0),
				FilteredBacktrace: make([]types.StackFrame, 0),
			},
		},
		{
			name: "Skipped test",
			input: InputTestResult{
				Name:      "Test name",
				Status:    "skipped",
				Message:   "",
				Backtrace: []string{},
			},
			expected: types.TestResult{
				Name:              "Test name",
				Status:            types.StatusSkip,
				ErrorMessage:      "",
				FullBacktrace:     make([]types.StackFrame, 0),
				FilteredBacktrace: make([]types.StackFrame, 0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertToTestResult(tt.input)
		assert.Equal(t, tt.expected.Name, result.Name)
		assert.Equal(t, tt.expected.Status, result.Status)
		assert.Equal(t, tt.expected.ErrorMessage, result.ErrorMessage)
		assert.Len(t, result.FullBacktrace, len(tt.expected.FullBacktrace))
		assert.Len(t, result.FilteredBacktrace, 0) // Check length instead of exact comparison
		})
	}
}

func TestBacktraceFrameLimit(t *testing.T) {
	// Create a test with more than 50 backtrace frames
	backtrace := make([]string, 60)
	for i := 0; i < 60; i++ {
		backtrace[i] = "app/test.rb:10:in 'test_method'"
	}

	input := InputTestResult{
		Name:      "Test with many frames",
		Status:    "failed",
		Message:   "Error",
		Backtrace: backtrace,
	}

	result := convertToTestResult(input)
	assert.Len(t, result.FullBacktrace, 50, "Backtrace should be capped at 50 frames")
}

func TestInvalidJSON(t *testing.T) {
	invalidJSON := `{"invalid": "json"`

	result, err := ParseJSON([]byte(invalidJSON))
	assert.Error(t, err)
	assert.Nil(t, result)
}
