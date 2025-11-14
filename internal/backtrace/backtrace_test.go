package backtrace

import (
	"testing"

	"github.com/adamakhtar/wing_commander/internal/projectfs"
	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestNewBacktrace(t *testing.T) {
	bt := NewBacktrace()
	assert.NotNil(t, bt)
	assert.Empty(t, bt.AllStackFrames())
}

func TestAppend_RelativePath(t *testing.T) {
	// Setup ProjectFS singleton for tests
	rootPath, _ := types.NewAbsPath("/path/to/project")
	projectfs.InitProjectFS(rootPath)

	bt := NewBacktrace()
	bt.Append("app/models/user.rb:42:in 'create_user'")

	frames := bt.AllStackFrames()
	assert.Len(t, frames, 1)
	assert.Equal(t, 42, frames[0].Line)
	assert.Equal(t, "create_user", frames[0].Function)
	// Path should be converted to absolute using ProjectFS
	assert.Equal(t, types.AbsPath("/path/to/project/app/models/user.rb"), frames[0].FilePath)
}

func TestAppend_AbsolutePath(t *testing.T) {
	// Setup ProjectFS singleton for tests
	rootPath, _ := types.NewAbsPath("/path/to/project")
	projectfs.InitProjectFS(rootPath)

	bt := NewBacktrace()
	bt.Append("/absolute/path/file.rb:10")

	frames := bt.AllStackFrames()
	assert.Len(t, frames, 1)
	assert.Equal(t, 10, frames[0].Line)
	assert.Equal(t, types.AbsPath("/absolute/path/file.rb"), frames[0].FilePath)
}

func TestAppend_RubyFormat(t *testing.T) {
	// Setup ProjectFS singleton for tests
	rootPath, _ := types.NewAbsPath("/path/to/project")
	projectfs.InitProjectFS(rootPath)

	bt := NewBacktrace()
	bt.Append("app/models/user.rb:42:in 'create_user'")

	frames := bt.AllStackFrames()
	assert.Len(t, frames, 1)
	// Path should be converted to absolute using ProjectFS
	assert.Equal(t, types.AbsPath("/path/to/project/app/models/user.rb"), frames[0].FilePath)
	assert.Equal(t, 42, frames[0].Line)
	assert.Equal(t, "create_user", frames[0].Function)
}

func TestAppend_RubyFormatWithoutMethod(t *testing.T) {
	// Setup ProjectFS singleton for tests
	rootPath, _ := types.NewAbsPath("/path/to/project")
	projectfs.InitProjectFS(rootPath)

	bt := NewBacktrace()
	bt.Append("app/models/user.rb:42")

	frames := bt.AllStackFrames()
	assert.Len(t, frames, 1)
	assert.Equal(t, 42, frames[0].Line)
	assert.Empty(t, frames[0].Function)
}

func TestAppend_PythonFormat(t *testing.T) {
	// Setup ProjectFS singleton for tests
	rootPath, _ := types.NewAbsPath("/path/to/project")
	projectfs.InitProjectFS(rootPath)

	bt := NewBacktrace()
	bt.Append("File \"app/models/user.py\", line 42, in create_user")

	frames := bt.AllStackFrames()
	assert.Len(t, frames, 1)
	// Python format is not fully parsed yet, but should create a frame
	assert.NotEmpty(t, frames[0].FilePath)
}

func TestAppend_EmptyString(t *testing.T) {
	// Setup ProjectFS singleton for tests
	rootPath, _ := types.NewAbsPath("/path/to/project")
	projectfs.InitProjectFS(rootPath)

	bt := NewBacktrace()
	bt.Append("")

	frames := bt.AllStackFrames()
	assert.Len(t, frames, 1)
	// Should create minimal StackFrame
	assert.Empty(t, frames[0].FilePath)
}

func TestAppend_InvalidFormat(t *testing.T) {
	// Setup ProjectFS singleton for tests
	rootPath, _ := types.NewAbsPath("/path/to/project")
	projectfs.InitProjectFS(rootPath)

	bt := NewBacktrace()
	bt.Append("invalid_frame")

	frames := bt.AllStackFrames()
	assert.Len(t, frames, 1)
	// Should create minimal StackFrame with path
	assert.NotEmpty(t, frames[0].FilePath)
}

func TestFilterProjectStackFramesOnly(t *testing.T) {
	// Setup ProjectFS singleton for tests
	rootPath, _ := types.NewAbsPath("/path/to/project")
	projectfs.InitProjectFS(rootPath)

	bt := NewBacktrace()
	bt.Append("/path/to/project/app/test.rb:10")
	bt.Append("/gems/rspec.rb:20")
	bt.Append("/path/to/project/app/another.rb:30")

	filtered := bt.FilterProjectStackFramesOnly()
	frames := filtered.AllStackFrames()

	assert.Len(t, frames, 2)
	assert.Equal(t, types.AbsPath("/path/to/project/app/test.rb"), frames[0].FilePath)
	assert.Equal(t, types.AbsPath("/path/to/project/app/another.rb"), frames[1].FilePath)
}

func TestFilterProjectStackFramesOnly_EmptyBacktrace(t *testing.T) {
	// Setup ProjectFS singleton for tests
	rootPath, _ := types.NewAbsPath("/path/to/project")
	projectfs.InitProjectFS(rootPath)

	bt := NewBacktrace()
	filtered := bt.FilterProjectStackFramesOnly()

	assert.Empty(t, filtered.AllStackFrames())
}

func TestFilterProjectStackFramesOnly_NoProjectFrames(t *testing.T) {
	// Setup ProjectFS singleton for tests
	rootPath, _ := types.NewAbsPath("/path/to/project")
	projectfs.InitProjectFS(rootPath)

	bt := NewBacktrace()
	bt.Append("/gems/rspec.rb:20")
	bt.Append("/usr/lib/ruby.rb:10")

	filtered := bt.FilterProjectStackFramesOnly()
	frames := filtered.AllStackFrames()

	assert.Empty(t, frames)
}

func TestFilterProjectStackFramesOnly_AllProjectFrames(t *testing.T) {
	// Setup ProjectFS singleton for tests
	rootPath, _ := types.NewAbsPath("/path/to/project")
	projectfs.InitProjectFS(rootPath)

	bt := NewBacktrace()
	bt.Append("/path/to/project/app/test.rb:10")
	bt.Append("/path/to/project/lib/helper.rb:20")

	filtered := bt.FilterProjectStackFramesOnly()
	frames := filtered.AllStackFrames()

	assert.Len(t, frames, 2)
}
