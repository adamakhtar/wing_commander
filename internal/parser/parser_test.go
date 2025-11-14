package parser

import (
	"testing"

	"github.com/adamakhtar/wing_commander/internal/backtrace"
	"github.com/adamakhtar/wing_commander/internal/projectfs"
	"github.com/adamakhtar/wing_commander/internal/testresult"
	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClassifyFailure_Heuristics(t *testing.T) {
	// Setup ProjectFS singleton for tests
	rootPath, _ := types.NewAbsPath("/path/to")
	err := projectfs.InitProjectFS(rootPath, "")
	require.NoError(t, err)

	cases := []struct {
		name     string
		message  string
		topFrame *types.StackFrame
		want     testresult.FailureCause
	}{
		{
			name:    "assertion by message",
			message: "Expected 2 to equal 3",
			want:    testresult.FailureCauseAssertion,
		},
		{
			name:    "test definition by spec path",
			message: "NoMethodError: undefined method",
			topFrame: &types.StackFrame{
				FilePath: types.AbsPath("/path/to/spec/models/user_spec.rb"),
				Line:     10,
			},
			want: testresult.FailureCauseTestDefinition,
		},
		{
			name:    "production by app path",
			message: "NoMethodError: undefined method",
			topFrame: &types.StackFrame{
				FilePath: types.AbsPath("/path/to/app/models/user.rb"),
				Line:     20,
			},
			want: testresult.FailureCauseProductionCode,
		},
		{
			name:    "no frames -> test definition",
			message: "",
			want:    testresult.FailureCauseTestDefinition,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, err := newParseContext(&ParseOptions{})
			require.NoError(t, err)
			got := classifyFailure(tc.message, tc.topFrame, ctx)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestClassifyFailure_UsesTestFilePattern(t *testing.T) {
	// Setup ProjectFS singleton for tests
	rootPath, _ := types.NewAbsPath("/abs/project")
	err := projectfs.InitProjectFS(rootPath, "test/**/*.rb")
	require.NoError(t, err)

	ctx, err := newParseContext(&ParseOptions{})
	require.NoError(t, err)

	cases := []struct {
		name     string
		topFrame *types.StackFrame
		want     testresult.FailureCause
	}{
		{
			name: "absolute path match",
			topFrame: &types.StackFrame{
				FilePath: types.AbsPath("/abs/project/test/models/user_test.rb"),
				Line:     12,
			},
			want: testresult.FailureCauseTestDefinition,
		},
		{
			name: "relative path match",
			topFrame: &types.StackFrame{
				FilePath: types.AbsPath("/abs/project/custom/subdir/test/models/user_test.rb"),
				Line:     8,
			},
			want: testresult.FailureCauseTestDefinition,
		},
		{
			name: "non matching path",
			topFrame: &types.StackFrame{
				FilePath: types.AbsPath("/abs/project/app/models/user.rb"),
				Line:     33,
			},
			want: testresult.FailureCauseProductionCode,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := classifyFailure("boom", tc.topFrame, ctx)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestParseStackFrame(t *testing.T) {
	// Setup ProjectFS for relative path conversion
	rootPath, _ := types.NewAbsPath("/path/to/project")
	err := projectfs.InitProjectFS(rootPath, "")
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    string
		expected types.StackFrame
	}{
		{
			name:  "Ruby with method",
			input: "app/models/user.rb:42:in 'create_user'",
			expected: types.StackFrame{
				FilePath: types.AbsPath("/path/to/project/app/models/user.rb"),
				Line:     42,
				Function: "create_user",
			},
		},
		{
			name:  "Ruby without method",
			input: "app/models/user.rb:42",
			expected: types.StackFrame{
				FilePath: types.AbsPath("/path/to/project/app/models/user.rb"),
				Line:     42,
				Function: "",
			},
		},
		{
			name:  "Python format",
			input: "File \"app/models/user.py\", line 42, in create_user",
			expected: types.StackFrame{
				// Python format is not fully parsed - whole string becomes path
				FilePath: types.AbsPath(""), // Will be set to converted path
				Line:     0,
				Function: "",
			},
		},
		{
			name:  "Invalid format",
			input: "invalid_frame",
			expected: types.StackFrame{
				FilePath: types.AbsPath("/path/to/project/invalid_frame"),
				Line:     0,
				Function: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bt := backtrace.NewBacktrace()
			bt.Append(tt.input)
			frames := bt.AllStackFrames()
			require.Len(t, frames, 1)
			result := frames[0]
			if tt.name == "Python format" {
				// Python format creates a path from the whole string
				assert.NotEmpty(t, result.FilePath)
			} else {
				assert.Equal(t, tt.expected.FilePath, result.FilePath)
			}
			assert.Equal(t, tt.expected.Line, result.Line)
			assert.Equal(t, tt.expected.Function, result.Function)
		})
	}
}

func TestParseFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "valid summary file",
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
			result, err := ParseFile(tt.filename, nil)
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

func TestParse_BasicFields(t *testing.T) {
	yamlData := `---
- test_group_name: TestClass
  test_status: passed
  duration: '1.23'
  test_file_path: "/path/to/test.rb"
  test_line_number: 42
  failure_details:
  failure_file_path:
  failure_line_number: 0
  full_backtrace: []
`

	result, err := Parse([]byte(yamlData), nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Tests, 1)

	test := result.Tests[0]
	assert.Equal(t, 1, test.Id)
	assert.Equal(t, "TestClass", test.GroupName)
	assert.Equal(t, "", test.TestCaseName)
	assert.Equal(t, testresult.StatusPass, test.Status)
	assert.Equal(t, 1.23, test.Duration)
	assert.Equal(t, types.AbsPath("/path/to/test.rb"), test.TestFilePath)
	assert.Equal(t, 42, test.TestLineNumber)
	assert.Equal(t, testresult.FailureCause(""), test.FailureCause)
	assert.Empty(t, test.FailureDetails)
	assert.Empty(t, test.FailureFilePath)
	assert.Empty(t, test.FullBacktrace.Frames)
}

func TestParse_StatusEnum(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected testresult.TestStatus
	}{
		{"passed", "passed", testresult.StatusPass},
		{"failed", "failed", testresult.StatusFail},
		{"skipped", "skipped", testresult.StatusSkip},
		{"unknown", "unknown", testresult.StatusFail}, // Default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yamlData := `---
- test_group_name: Test
  test_status: ` + tt.status + `
  duration: '0.00'
  test_file_path: ""
  test_line_number: 0
  full_backtrace: []
`
			result, err := Parse([]byte(yamlData), nil)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result.Tests[0].Status)
		})
	}
}

func TestParse_FailureCauseClassification(t *testing.T) {
	tests := []struct {
		name     string
		yamlData string
		expected testresult.FailureCause
	}{
		{
			name: "assertion by message",
			yamlData: `---
- test_group_name: Test
  test_status: failed
  duration: '0.00'
  test_file_path: ""
  test_line_number: 0
  failure_details: "Expected: foo\n  Actual: bar"
  failure_file_path: "/path/to/test.rb"
  failure_line_number: 21
  full_backtrace:
    - "/path/to/test.rb:21:in 'test_method'"
`,
			expected: testresult.FailureCauseAssertion,
		},
		{
			name: "production code by app frame",
			yamlData: `---
- test_group_name: Test
  test_status: failed
  duration: '0.00'
  test_file_path: ""
  test_line_number: 0
  failure_details: "RuntimeError: boom"
  failure_file_path: "/app/models/user.rb"
  failure_line_number: 14
  full_backtrace:
    - "/app/models/user.rb:14:in 'explode'"
    - "/Users/me/.asdf/installs/ruby/3.3.0/lib/ruby/gems/3.3.0/gems/minitest/test.rb:90:in 'run'"
`,
			expected: testresult.FailureCauseProductionCode,
		},
		{
			name: "test definition by test frame",
			yamlData: `---
- test_group_name: Test
  test_status: failed
  duration: '0.00'
  test_file_path: ""
  test_line_number: 0
  failure_details: "NoMethodError: undefined method"
  failure_file_path: "/test/models/user_test.rb"
  failure_line_number: 10
  full_backtrace:
    - "/test/models/user_test.rb:10:in 'block in <class:UserTest>'"
    - "/Users/me/.asdf/installs/ruby/3.3.0/lib/ruby/gems/3.3.0/gems/minitest/test.rb:90:in 'run'"
`,
			expected: testresult.FailureCauseTestDefinition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse([]byte(tt.yamlData), nil)
			require.NoError(t, err)
			require.Len(t, result.Tests, 1)
			assert.Equal(t, tt.expected, result.Tests[0].FailureCause)
		})
	}
}

func TestParse_FailureFields(t *testing.T) {
	yamlData := `---
- test_group_name: TestClass
  test_status: failed
  duration: '0.50'
  test_file_path: "/path/to/test.rb"
  test_line_number: 10
  failure_details: "NameError: undefined method"
  failure_file_path: "/path/to/error.rb"
  failure_line_number: 15
  full_backtrace:
    - "/path/to/error.rb:15:in 'method'"
`

	result, err := Parse([]byte(yamlData), nil)
	require.NoError(t, err)
	require.Len(t, result.Tests, 1)

	test := result.Tests[0]
	assert.Equal(t, testresult.StatusFail, test.Status)
	assert.Equal(t, testresult.FailureCauseProductionCode, test.FailureCause)
	assert.Equal(t, "NameError: undefined method", test.FailureDetails)
	assert.Equal(t, types.AbsPath("/path/to/error.rb"), test.FailureFilePath)
	assert.Equal(t, 15, test.FailureLineNumber)
	assert.Len(t, test.FullBacktrace.Frames, 1)
	assert.Equal(t, types.AbsPath("/path/to/error.rb"), test.FullBacktrace.Frames[0].FilePath)
	assert.Equal(t, 15, test.FullBacktrace.Frames[0].Line)
}

func TestParse_AssertionFields(t *testing.T) {
	yamlData := `---
- test_group_name: TestClass
  test_status: failed
  duration: '0.30'
  test_file_path: "/path/to/test.rb"
  test_line_number: 20
  failure_details: "Expected: foo\n  Actual: bar"
  failure_file_path: "/path/to/test.rb"
  failure_line_number: 21
  full_backtrace:
    - "/path/to/test.rb:21:in 'test_method'"
`

	result, err := Parse([]byte(yamlData), nil)
	require.NoError(t, err)
	require.Len(t, result.Tests, 1)

	test := result.Tests[0]
	assert.Equal(t, testresult.StatusFail, test.Status)
	assert.Equal(t, testresult.FailureCauseAssertion, test.FailureCause)
	assert.Equal(t, "Expected: foo\n  Actual: bar", test.FailureDetails)
	assert.Equal(t, types.AbsPath("/path/to/test.rb"), test.FailureFilePath)
	assert.Equal(t, 21, test.FailureLineNumber)
}

func TestParse_FailureCause_UsesPattern(t *testing.T) {
	// Setup ProjectFS singleton for tests
	rootPath, _ := types.NewAbsPath("/abs/project")
	err := projectfs.InitProjectFS(rootPath, "custom_specs/**/*.rb")
	require.NoError(t, err)

	yamlData := `---
- test_group_name: CustomSuite
  test_status: failed
  duration: '0.10'
  test_file_path: "/abs/project/custom_specs/models/user_case.rb"
  test_line_number: 5
  failure_details: "Boom"
  failure_file_path: "/abs/project/custom_specs/models/user_case.rb"
  failure_line_number: 7
  full_backtrace:
    - "/abs/project/custom_specs/models/user_case.rb:7:in 'test_case'"
`

	opts := &ParseOptions{}

	result, err := Parse([]byte(yamlData), opts)
	require.NoError(t, err)
	require.Len(t, result.Tests, 1)

	test := result.Tests[0]
	assert.Equal(t, testresult.FailureCauseTestDefinition, test.FailureCause)
}

func TestParse_BacktraceParsing(t *testing.T) {
	yamlData := `---
- test_group_name: TestClass
  test_status: failed
  duration: '0.00'
  test_file_path: ""
  test_line_number: 0
  full_backtrace:
    - "/path/to/file.rb:42:in 'method_name'"
    - "/path/to/other.rb:10"
    - "invalid_frame"
    - "/valid/path.rb:5:in 'valid_method'"
`

	result, err := Parse([]byte(yamlData), nil)
	require.NoError(t, err)
	require.Len(t, result.Tests, 1)

	test := result.Tests[0]
	assert.Len(t, test.FullBacktrace.Frames, 3) // invalid_frame should be skipped

	// Check first frame
	assert.Equal(t, types.AbsPath("/path/to/file.rb"), test.FullBacktrace.Frames[0].FilePath)
	assert.Equal(t, 42, test.FullBacktrace.Frames[0].Line)
	assert.Equal(t, "method_name", test.FullBacktrace.Frames[0].Function)

	// Check second frame
	assert.Equal(t, types.AbsPath("/path/to/other.rb"), test.FullBacktrace.Frames[1].FilePath)
	assert.Equal(t, 10, test.FullBacktrace.Frames[1].Line)

	// Check third frame
	assert.Equal(t, types.AbsPath("/valid/path.rb"), test.FullBacktrace.Frames[2].FilePath)
	assert.Equal(t, 5, test.FullBacktrace.Frames[2].Line)
}

func TestParse_EmptyArray(t *testing.T) {
	yamlData := `--- []`

	result, err := Parse([]byte(yamlData), nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Tests, 0)
	assert.Equal(t, 0, result.Summary.Total)
	assert.Equal(t, 0, result.Summary.Passed)
	assert.Equal(t, 0, result.Summary.Failed)
	assert.Equal(t, 0, result.Summary.Skipped)
}

func TestParse_MissingFields(t *testing.T) {
	yamlData := `---
- test_group_name: Test
  test_status: passed
`

	result, err := Parse([]byte(yamlData), nil)
	require.NoError(t, err)
	require.Len(t, result.Tests, 1)

	test := result.Tests[0]
	assert.Equal(t, "Test", test.GroupName)
	assert.Equal(t, "", test.TestCaseName)
	assert.Equal(t, testresult.StatusPass, test.Status)
	assert.Equal(t, 0.0, test.Duration) // Default when missing
	assert.Empty(t, test.TestFilePath)
	assert.Equal(t, 0, test.TestLineNumber)
	assert.Empty(t, test.FailureDetails)
}

func TestParse_InvalidDuration(t *testing.T) {
	yamlData := `---
- test_group_name: Test
  test_status: passed
  duration: "invalid"
`

	result, err := Parse([]byte(yamlData), nil)
	require.NoError(t, err)
	require.Len(t, result.Tests, 1)

	test := result.Tests[0]
	assert.Equal(t, 0.0, test.Duration) // Default on parse error
}

func TestParse_MultipleTests(t *testing.T) {
	yamlData := `---
- test_group_name: Test1
  test_status: passed
  duration: '1.0'
  test_file_path: ""
  test_line_number: 0
  full_backtrace: []
- test_group_name: Test2
  test_status: failed
  duration: '2.0'
  test_file_path: ""
  test_line_number: 0
  full_backtrace: []
- test_group_name: Test3
  test_status: skipped
  duration: '0.5'
  test_file_path: ""
  test_line_number: 0
  full_backtrace: []
`

	result, err := Parse([]byte(yamlData), nil)
	require.NoError(t, err)
	require.Len(t, result.Tests, 3)

	assert.Equal(t, 3, result.Summary.Total)
	assert.Equal(t, 1, result.Summary.Passed)
	assert.Equal(t, 1, result.Summary.Failed)
	assert.Equal(t, 1, result.Summary.Skipped)

	assert.Equal(t, "Test1", result.Tests[0].GroupName)
	assert.Equal(t, "", result.Tests[0].TestCaseName)
	assert.Equal(t, testresult.StatusPass, result.Tests[0].Status)
	assert.Equal(t, 1.0, result.Tests[0].Duration)

	assert.Equal(t, "Test2", result.Tests[1].GroupName)
	assert.Equal(t, "", result.Tests[1].TestCaseName)
	assert.Equal(t, testresult.StatusFail, result.Tests[1].Status)
	assert.Equal(t, 2.0, result.Tests[1].Duration)

	assert.Equal(t, "Test3", result.Tests[2].GroupName)
	assert.Equal(t, "", result.Tests[2].TestCaseName)
	assert.Equal(t, testresult.StatusSkip, result.Tests[2].Status)
	assert.Equal(t, 0.5, result.Tests[2].Duration)
}

func TestParse_IntAsString(t *testing.T) {
	yamlData := `---
- test_group_name: Test
  test_status: passed
  duration: '0.00'
  test_file_path: ""
  test_line_number: "42"
  full_backtrace: []
`

	result, err := Parse([]byte(yamlData), nil)
	require.NoError(t, err)
	require.Len(t, result.Tests, 1)

	test := result.Tests[0]
	assert.Equal(t, 42, test.TestLineNumber) // Should parse string to int
}

func TestParse_InvalidData(t *testing.T) {
	invalidYAML := `--- invalid: yaml: content`

	result, err := Parse([]byte(invalidYAML), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
}
