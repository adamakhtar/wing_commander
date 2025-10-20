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
	assert.Equal(t, "User.should be valid", test.Name)
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
	assert.Equal(t, "UserTest.test_user_creation", test.Name)
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
