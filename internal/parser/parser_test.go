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
			name:     "valid RSpec XML file",
			filename: "../../testdata/fixtures/rspec_failures.xml",
			wantErr:  false,
		},
		{
			name:     "valid Minitest XML file",
			filename: "../../testdata/fixtures/minitest_failures.xml",
			wantErr:  false,
		},
		{
			name:     "non-existent file",
			filename: "testdata/fixtures/nonexistent.xml",
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

func TestParseXML_RSpec(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
  <testsuite name="RSpec" tests="1" failures="1" skipped="0" time="0.123">
    <testcase classname="User" name="should be valid" time="0.123">
      <failure message="Expected User to be valid">
        app/models/user.rb:42:in 'create_user'
        spec/models/user_spec.rb:15:in 'block (2 levels) in &lt;top (required)&gt;'
      </failure>
    </testcase>
  </testsuite>
</testsuites>`

	result, err := ParseXML([]byte(xmlData))
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Len(t, result.Tests, 1)

	test := result.Tests[0]
	assert.Equal(t, "User should be valid", test.Name)
	assert.Equal(t, types.StatusFail, test.Status)
	assert.Equal(t, "Expected User to be valid", test.ErrorMessage)
    assert.Len(t, test.FullBacktrace, 2)
    // Expect assertion failure classification due to message
    assert.Equal(t, types.FailureCauseAssertion, test.FailureCause)

	// Check first frame
	frame := test.FullBacktrace[0]
	assert.Equal(t, "app/models/user.rb", frame.File)
	assert.Equal(t, 42, frame.Line)
	assert.Equal(t, "create_user", frame.Function)
}

func TestParseXML_Minitest(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
  <testsuite name="Minitest" tests="1" failures="1" skipped="0" time="0.156">
    <testcase classname="UserTest" name="test_user_creation" time="0.156">
      <failure message="Expected User to be valid">
        app/models/user.rb:42:in 'create_user'
        test/models/user_test.rb:15:in 'test_user_creation'
      </failure>
    </testcase>
  </testsuite>
</testsuites>`

	result, err := ParseXML([]byte(xmlData))
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 1, result.Summary.Total)
	assert.Equal(t, 1, result.Summary.Failed)
	assert.Len(t, result.Tests, 1)

	test := result.Tests[0]
	assert.Equal(t, "UserTest test_user_creation", test.Name)
	assert.Equal(t, types.StatusFail, test.Status)
    // Expect assertion failure classification due to message
    assert.Equal(t, types.FailureCauseAssertion, test.FailureCause)
}

func TestClassifyFailure_Heuristics(t *testing.T) {
    cases := []struct {
        name    string
        message string
        frames  []types.StackFrame
        want    types.FailureCause
    }{
        {
            name:    "assertion by message",
            message: "Expected 2 to equal 3",
            frames:  nil,
            want:    types.FailureCauseAssertion,
        },
        {
            name:    "test definition by spec path",
            message: "NoMethodError: undefined method",
            frames: []types.StackFrame{
                {File: "spec/models/user_spec.rb", Line: 10},
                {File: "app/models/user.rb", Line: 20},
            },
            want: types.FailureCauseTestDefinition,
        },
        {
            name:    "production by app path",
            message: "NoMethodError: undefined method",
            frames: []types.StackFrame{
                {File: "app/models/user.rb", Line: 20},
                {File: "gems/rspec-core/example.rb", Line: 100},
            },
            want: types.FailureCauseProductionCode,
        },
        {
            name:    "no frames -> test definition",
            message: "",
            frames:  nil,
            want:    types.FailureCauseTestDefinition,
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            got := classifyFailure(tc.message, tc.frames)
            assert.Equal(t, tc.want, got)
        })
    }
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


func TestInvalidXML(t *testing.T) {
	invalidXML := `<?xml version="1.0"?><invalid>xml`

	result, err := ParseXML([]byte(invalidXML))
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestParseYAMLFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "valid YAML summary file",
			filename: "../../testdata/fixtures/minitest_summary.yml",
			wantErr:  false,
		},
		{
			name:     "non-existent file",
			filename: "testdata/fixtures/nonexistent.yml",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseYAMLFile(tt.filename)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Greater(t, len(result.Tests), 0)
			}
		})
	}
}

func TestParseYAML_BasicFields(t *testing.T) {
	yamlData := `---
- test_case_name: TestClass
  test_status: passed
  duration: '1.23'
  test_file_path: "/path/to/test.rb"
  test_line_number: 42
  failure_cause:
  error_message:
  error_file_path:
  error_line_number: 0
  failed_assertion_details:
  assertion_file_path:
  assertion_line_number: 0
  full_backtrace: []
`

	result, err := ParseYAML([]byte(yamlData))
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Tests, 1)

	test := result.Tests[0]
	assert.Equal(t, 1, test.Id)
	assert.Equal(t, "TestClass", test.Name)
	assert.Equal(t, types.StatusPass, test.Status)
	assert.Equal(t, 1.23, test.Duration)
	assert.Equal(t, "/path/to/test.rb", test.TestFilePath)
	assert.Equal(t, 42, test.TestLineNumber)
	assert.Equal(t, types.FailureCause(""), test.FailureCause)
	assert.Empty(t, test.ErrorMessage)
	assert.Empty(t, test.FailedAssertionMessage)
	assert.Empty(t, test.FullBacktrace)
}

func TestParseYAML_StatusEnum(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected types.TestStatus
	}{
		{"passed", "passed", types.StatusPass},
		{"failed", "failed", types.StatusFail},
		{"skipped", "skipped", types.StatusSkip},
		{"unknown", "unknown", types.StatusFail}, // Default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yamlData := `---
- test_case_name: Test
  test_status: ` + tt.status + `
  duration: '0.00'
  test_file_path: ""
  test_line_number: 0
  full_backtrace: []
`
			result, err := ParseYAML([]byte(yamlData))
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result.Tests[0].Status)
		})
	}
}

func TestParseYAML_FailureCauseEnum(t *testing.T) {
	tests := []struct {
		name         string
		failureCause string
		expected     types.FailureCause
	}{
		{"error", "error", types.FailureCauseProductionCode},
		{"failed_assertion", "failed_assertion", types.FailureCauseAssertion},
		{"empty", "", types.FailureCause("")},
		{"unknown", "unknown", types.FailureCause("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yamlData := `---
- test_case_name: Test
  test_status: failed
  duration: '0.00'
  test_file_path: ""
  test_line_number: 0
  failure_cause: ` + tt.failureCause + `
  full_backtrace: []
`
			result, err := ParseYAML([]byte(yamlData))
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result.Tests[0].FailureCause)
		})
	}
}

func TestParseYAML_ErrorFields(t *testing.T) {
	yamlData := `---
- test_case_name: TestClass
  test_status: failed
  duration: '0.50'
  test_file_path: "/path/to/test.rb"
  test_line_number: 10
  failure_cause: error
  error_message: "NameError: undefined method"
  error_file_path: "/path/to/error.rb"
  error_line_number: 15
  failed_assertion_details:
  assertion_file_path:
  assertion_line_number: 0
  full_backtrace:
    - "/path/to/error.rb:15:in 'method'"
`

	result, err := ParseYAML([]byte(yamlData))
	require.NoError(t, err)
	require.Len(t, result.Tests, 1)

	test := result.Tests[0]
	assert.Equal(t, types.StatusFail, test.Status)
	assert.Equal(t, types.FailureCauseProductionCode, test.FailureCause)
	assert.Equal(t, "NameError: undefined method", test.ErrorMessage)
	assert.Equal(t, "/path/to/error.rb", test.ErrorFilePath)
	assert.Equal(t, 15, test.ErrorLineNumber)
	assert.Len(t, test.FullBacktrace, 1)
	assert.Equal(t, "/path/to/error.rb", test.FullBacktrace[0].File)
	assert.Equal(t, 15, test.FullBacktrace[0].Line)
}

func TestParseYAML_AssertionFields(t *testing.T) {
	yamlData := `---
- test_case_name: TestClass
  test_status: failed
  duration: '0.30'
  test_file_path: "/path/to/test.rb"
  test_line_number: 20
  failure_cause: failed_assertion
  error_message:
  error_file_path:
  error_line_number: 0
  failed_assertion_details: "Expected: foo\n  Actual: bar"
  assertion_file_path: "/path/to/test.rb"
  assertion_line_number: 21
  full_backtrace:
    - "/path/to/test.rb:21:in 'test_method'"
`

	result, err := ParseYAML([]byte(yamlData))
	require.NoError(t, err)
	require.Len(t, result.Tests, 1)

	test := result.Tests[0]
	assert.Equal(t, types.StatusFail, test.Status)
	assert.Equal(t, types.FailureCauseAssertion, test.FailureCause)
	assert.Equal(t, "Expected: foo\n  Actual: bar", test.FailedAssertionMessage)
	assert.Equal(t, "/path/to/test.rb", test.FailedAssertionFilePath)
	assert.Equal(t, 21, test.FailedAssertionLineNumber)
	assert.Empty(t, test.ErrorMessage)
}

func TestParseYAML_BacktraceParsing(t *testing.T) {
	yamlData := `---
- test_case_name: TestClass
  test_status: failed
  duration: '0.00'
  test_file_path: ""
  test_line_number: 0
  failure_cause: error
  full_backtrace:
    - "/path/to/file.rb:42:in 'method_name'"
    - "/path/to/other.rb:10"
    - "invalid_frame"
    - "/valid/path.rb:5:in 'valid_method'"
`

	result, err := ParseYAML([]byte(yamlData))
	require.NoError(t, err)
	require.Len(t, result.Tests, 1)

	test := result.Tests[0]
	assert.Len(t, test.FullBacktrace, 3) // invalid_frame should be skipped

	// Check first frame
	assert.Equal(t, "/path/to/file.rb", test.FullBacktrace[0].File)
	assert.Equal(t, 42, test.FullBacktrace[0].Line)
	assert.Equal(t, "method_name", test.FullBacktrace[0].Function)

	// Check second frame
	assert.Equal(t, "/path/to/other.rb", test.FullBacktrace[1].File)
	assert.Equal(t, 10, test.FullBacktrace[1].Line)

	// Check third frame
	assert.Equal(t, "/valid/path.rb", test.FullBacktrace[2].File)
	assert.Equal(t, 5, test.FullBacktrace[2].Line)
}

func TestParseYAML_EmptyArray(t *testing.T) {
	yamlData := `--- []`

	result, err := ParseYAML([]byte(yamlData))
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Tests, 0)
	assert.Equal(t, 0, result.Summary.Total)
	assert.Equal(t, 0, result.Summary.Passed)
	assert.Equal(t, 0, result.Summary.Failed)
	assert.Equal(t, 0, result.Summary.Skipped)
}

func TestParseYAML_MissingFields(t *testing.T) {
	yamlData := `---
- test_case_name: Test
  test_status: passed
`

	result, err := ParseYAML([]byte(yamlData))
	require.NoError(t, err)
	require.Len(t, result.Tests, 1)

	test := result.Tests[0]
	assert.Equal(t, "Test", test.Name)
	assert.Equal(t, types.StatusPass, test.Status)
	assert.Equal(t, 0.0, test.Duration) // Default when missing
	assert.Empty(t, test.TestFilePath)
	assert.Equal(t, 0, test.TestLineNumber)
	assert.Empty(t, test.ErrorMessage)
}

func TestParseYAML_InvalidDuration(t *testing.T) {
	yamlData := `---
- test_case_name: Test
  test_status: passed
  duration: "invalid"
`

	result, err := ParseYAML([]byte(yamlData))
	require.NoError(t, err)
	require.Len(t, result.Tests, 1)

	test := result.Tests[0]
	assert.Equal(t, 0.0, test.Duration) // Default on parse error
}

func TestParseYAML_MultipleTests(t *testing.T) {
	yamlData := `---
- test_case_name: Test1
  test_status: passed
  duration: '1.0'
  test_file_path: ""
  test_line_number: 0
  full_backtrace: []
- test_case_name: Test2
  test_status: failed
  duration: '2.0'
  test_file_path: ""
  test_line_number: 0
  failure_cause: error
  full_backtrace: []
- test_case_name: Test3
  test_status: skipped
  duration: '0.5'
  test_file_path: ""
  test_line_number: 0
  full_backtrace: []
`

	result, err := ParseYAML([]byte(yamlData))
	require.NoError(t, err)
	require.Len(t, result.Tests, 3)

	assert.Equal(t, 3, result.Summary.Total)
	assert.Equal(t, 1, result.Summary.Passed)
	assert.Equal(t, 1, result.Summary.Failed)
	assert.Equal(t, 1, result.Summary.Skipped)

	assert.Equal(t, "Test1", result.Tests[0].Name)
	assert.Equal(t, types.StatusPass, result.Tests[0].Status)
	assert.Equal(t, 1.0, result.Tests[0].Duration)

	assert.Equal(t, "Test2", result.Tests[1].Name)
	assert.Equal(t, types.StatusFail, result.Tests[1].Status)
	assert.Equal(t, 2.0, result.Tests[1].Duration)

	assert.Equal(t, "Test3", result.Tests[2].Name)
	assert.Equal(t, types.StatusSkip, result.Tests[2].Status)
	assert.Equal(t, 0.5, result.Tests[2].Duration)
}

func TestParseYAML_IntAsString(t *testing.T) {
	yamlData := `---
- test_case_name: Test
  test_status: passed
  duration: '0.00'
  test_file_path: ""
  test_line_number: "42"
  full_backtrace: []
`

	result, err := ParseYAML([]byte(yamlData))
	require.NoError(t, err)
	require.Len(t, result.Tests, 1)

	test := result.Tests[0]
	assert.Equal(t, 42, test.TestLineNumber) // Should parse string to int
}

func TestParseYAML_InvalidYAML(t *testing.T) {
	invalidYAML := `--- invalid: yaml: content`

	result, err := ParseYAML([]byte(invalidYAML))
	assert.Error(t, err)
	assert.Nil(t, result)
}
