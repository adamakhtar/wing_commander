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
	test := NewTestResult("User creation", StatusFail)

	assert.Equal(t, "User creation", test.Name)
	assert.Equal(t, StatusFail, test.Status)
	assert.Empty(t, test.ErrorMessage)
	assert.Empty(t, test.FullBacktrace)
	assert.Empty(t, test.FilteredBacktrace)
}

func TestNewFailureGroup(t *testing.T) {
	group := NewFailureGroup("abc123", "Validation failed")

	assert.Equal(t, "abc123", group.Hash)
	assert.Equal(t, "Validation failed", group.ErrorMessage)
	assert.Empty(t, group.NormalizedBacktrace)
	assert.Empty(t, group.Tests)
	assert.Equal(t, 0, group.Count)
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
		Name:              "Test name",
		Status:            StatusPass,
		ErrorMessage:      "Error message",
		FullBacktrace:     []StackFrame{frame},
		FilteredBacktrace: []StackFrame{frame},
	}

	assert.Equal(t, "Test name", test.Name)
	assert.Equal(t, StatusPass, test.Status)
	assert.Equal(t, "Error message", test.ErrorMessage)
	assert.Len(t, test.FullBacktrace, 1)
	assert.Len(t, test.FilteredBacktrace, 1)
}

func TestFailureGroupFields(t *testing.T) {
	test := NewTestResult("Test", StatusFail)
	group := FailureGroup{
		Hash:                "hash123",
		ErrorMessage:        "Error",
		NormalizedBacktrace: []StackFrame{},
		Tests:               []TestResult{test},
		Count:               1,
	}

	assert.Equal(t, "hash123", group.Hash)
	assert.Equal(t, "Error", group.ErrorMessage)
	assert.Empty(t, group.NormalizedBacktrace)
	assert.Len(t, group.Tests, 1)
	assert.Equal(t, 1, group.Count)
}
