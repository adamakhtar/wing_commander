package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStackFrame(t *testing.T) {
	absPath, _ := NewAbsPath("app/models/user.rb")
	frame := NewStackFrame(absPath, 42, "create_user")

	assert.Equal(t, absPath, frame.FilePath)
	assert.Equal(t, 42, frame.Line)
	assert.Equal(t, "create_user", frame.Function)
}

func TestStackFrameFields(t *testing.T) {
	absPath, _ := NewAbsPath("test.rb")
	frame := StackFrame{
		FilePath: absPath,
		Line:     10,
		Function: "test_method",
	}

	assert.Equal(t, absPath, frame.FilePath)
	assert.Equal(t, 10, frame.Line)
	assert.Equal(t, "test_method", frame.Function)
}
