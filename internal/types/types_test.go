package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStackFrame(t *testing.T) {
	frame := NewStackFrame("app/models/user.rb", 42, "create_user")

	assert.Equal(t, "app/models/user.rb", frame.File)
	assert.Equal(t, 42, frame.Line)
	assert.Equal(t, "create_user", frame.Function)
}

func TestNewTestResult(t *testing.T) {
	test := NewTestResult("UserTest", "test_user_creation", StatusFail)

	assert.Equal(t, "UserTest", test.GroupName)
	assert.Equal(t, "test_user_creation", test.TestCaseName)
	assert.Equal(t, StatusFail, test.Status)
	assert.Empty(t, test.FailureDetails)
	assert.Empty(t, test.FullBacktrace)
	assert.Empty(t, test.FilteredBacktrace)
}

func TestTestStatusConstants(t *testing.T) {
	assert.Equal(t, TestStatus("pass"), StatusPass)
	assert.Equal(t, TestStatus("fail"), StatusFail)
	assert.Equal(t, TestStatus("skip"), StatusSkip)
}

func TestStackFrameFields(t *testing.T) {
	frame := StackFrame{
		File:     "test.rb",
		Line:     10,
		Function: "test_method",
	}

	assert.Equal(t, "test.rb", frame.File)
	assert.Equal(t, 10, frame.Line)
	assert.Equal(t, "test_method", frame.Function)
}

func TestTestResultFields(t *testing.T) {
	frame := NewStackFrame("test.rb", 1, "test")
	test := TestResult{
		GroupName:         "TestClass",
		TestCaseName:      "test_method",
		Status:            StatusPass,
		FailureDetails:    "Error message",
		FullBacktrace:     []StackFrame{frame},
		FilteredBacktrace: []StackFrame{frame},
	}

	assert.Equal(t, "TestClass", test.GroupName)
	assert.Equal(t, "test_method", test.TestCaseName)
	assert.Equal(t, StatusPass, test.Status)
	assert.Equal(t, "Error message", test.FailureDetails)
	assert.Len(t, test.FullBacktrace, 1)
	assert.Len(t, test.FilteredBacktrace, 1)
}
