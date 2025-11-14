package testresult

import (
	"testing"

	"github.com/adamakhtar/wing_commander/internal/backtrace"
	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestNewTestResult(t *testing.T) {
	test := NewTestResult("UserTest", "test_user_creation", StatusFail)

	assert.Equal(t, "UserTest", test.GroupName)
	assert.Equal(t, "test_user_creation", test.TestCaseName)
	assert.Equal(t, StatusFail, test.Status)
	assert.Empty(t, test.FailureDetails)
	assert.Empty(t, test.FullBacktrace.Frames)
	assert.Empty(t, test.FilteredBacktrace.Frames)
}

func TestTestStatusConstants(t *testing.T) {
	assert.Equal(t, TestStatus("pass"), StatusPass)
	assert.Equal(t, TestStatus("fail"), StatusFail)
	assert.Equal(t, TestStatus("skip"), StatusSkip)
}

func TestTestResultFields(t *testing.T) {
	absPath, _ := types.NewAbsPath("/test.rb")
	frame := types.NewStackFrame(absPath, 1, "test")
	test := TestResult{
		GroupName:         "TestClass",
		TestCaseName:      "test_method",
		Status:            StatusPass,
		FailureDetails:    "Error message",
		FullBacktrace:     backtrace.Backtrace{Frames: []types.StackFrame{frame}},
		FilteredBacktrace: backtrace.Backtrace{Frames: []types.StackFrame{frame}},
	}

	assert.Equal(t, "TestClass", test.GroupName)
	assert.Equal(t, "test_method", test.TestCaseName)
	assert.Equal(t, StatusPass, test.Status)
	assert.Equal(t, "Error message", test.FailureDetails)
	assert.Len(t, test.FullBacktrace.Frames, 1)
	assert.Len(t, test.FilteredBacktrace.Frames, 1)
}
