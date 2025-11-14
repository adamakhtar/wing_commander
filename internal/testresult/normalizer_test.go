package testresult

import (
	"testing"

	"github.com/adamakhtar/wing_commander/internal/backtrace"
	"github.com/adamakhtar/wing_commander/internal/projectfs"
	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestNewNormalizer(t *testing.T) {
	// Setup ProjectFS singleton for tests
	rootPath, _ := types.NewAbsPath("/path/to/project")
	err := projectfs.InitProjectFS(rootPath, "")
	if err != nil {
		t.Fatalf("failed to initialize ProjectFS: %v", err)
	}

	normalizer := NewNormalizer()

	assert.NotNil(t, normalizer)
}

func TestNormalizeTestResults(t *testing.T) {
	// Setup ProjectFS singleton for tests
	rootPath, _ := types.NewAbsPath("/path/to/project")
	err := projectfs.InitProjectFS(rootPath, "")
	if err != nil {
		t.Fatalf("failed to initialize ProjectFS: %v", err)
	}

	normalizer := NewNormalizer()

	results := []TestResult{
		{
			GroupName: "Test 1",
			Status:    StatusFail,
			FullBacktrace: backtrace.Backtrace{
				Frames: []types.StackFrame{
					{FilePath: types.AbsPath("/path/to/project/app/test.rb"), Line: 10},
					{FilePath: types.AbsPath("/gems/rspec.rb"), Line: 20},
				},
			},
		},
		{
			GroupName: "Test 2",
			Status:    StatusFail,
			FullBacktrace: backtrace.Backtrace{
				Frames: []types.StackFrame{
					{FilePath: types.AbsPath("/path/to/project/app/another.rb"), Line: 30},
				},
			},
		},
	}

	normalized := normalizer.NormalizeTestResults(results)

	assert.Len(t, normalized, 2)
	assert.Len(t, normalized[0].FilteredBacktrace.Frames, 1)
	assert.Equal(t, types.AbsPath("/path/to/project/app/test.rb"), normalized[0].FilteredBacktrace.Frames[0].FilePath)
	assert.Len(t, normalized[1].FilteredBacktrace.Frames, 1)
	assert.Equal(t, types.AbsPath("/path/to/project/app/another.rb"), normalized[1].FilteredBacktrace.Frames[0].FilePath)
}
