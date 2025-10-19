package grouper

import (
	"testing"

	"github.com/adamakhtar/wing_commander/internal/config"
	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestNewNormalizer(t *testing.T) {
	cfg := &config.Config{
		ExcludePatterns: []string{"/gems/", "/lib/ruby/"},
	}

	normalizer := NewNormalizer(cfg)

	assert.NotNil(t, normalizer)
	assert.Equal(t, cfg.ExcludePatterns, normalizer.excludePatterns)
}


func TestNormalizeTestResults(t *testing.T) {
	normalizer := &Normalizer{
		excludePatterns: []string{"/gems/"},
	}

	results := []types.TestResult{
		{
			Name:   "Test 1",
			Status: types.StatusFail,
			FullBacktrace: []types.StackFrame{
				{File: "app/test.rb", Line: 10},
				{File: "/gems/rspec.rb", Line: 20},
			},
		},
		{
			Name:   "Test 2",
			Status: types.StatusFail,
			FullBacktrace: []types.StackFrame{
				{File: "app/another.rb", Line: 30},
			},
		},
	}

	normalized := normalizer.NormalizeTestResults(results)

	assert.Len(t, normalized, 2)
	assert.Len(t, normalized[0].FilteredBacktrace, 1)
	assert.Len(t, normalized[1].FilteredBacktrace, 1)
}

func TestNormalizeFrameForGrouping(t *testing.T) {
	tests := []struct {
		name     string
		frame    types.StackFrame
		expected string
	}{
		{
			name:     "Frame with function",
			frame:    types.StackFrame{File: "app/models/user.rb", Line: 42, Function: "create_user"},
			expected: "app/models/user.rb::create_user",
		},
		{
			name:     "Frame without function",
			frame:    types.StackFrame{File: "app/models/user.rb", Line: 42},
			expected: "app/models/user.rb",
		},
		{
			name:     "Frame with different line number",
			frame:    types.StackFrame{File: "app/models/user.rb", Line: 100, Function: "create_user"},
			expected: "app/models/user.rb::create_user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeFrameForGrouping(tt.frame)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizeBacktraceForGrouping(t *testing.T) {
	tests := []struct {
		name     string
		frames   []types.StackFrame
		expected string
	}{
		{
			name:     "Empty backtrace",
			frames:   []types.StackFrame{},
			expected: "",
		},
		{
			name: "Single frame",
			frames: []types.StackFrame{
				{File: "app/user.rb", Line: 10, Function: "test"},
			},
			expected: "app/user.rb::test",
		},
		{
			name: "Multiple frames",
			frames: []types.StackFrame{
				{File: "app/user.rb", Line: 10, Function: "create"},
				{File: "app/service.rb", Line: 20, Function: "process"},
				{File: "app/controller.rb", Line: 30, Function: "index"},
			},
			expected: "app/user.rb::create|app/service.rb::process|app/controller.rb::index",
		},
		{
			name: "Same file different line numbers",
			frames: []types.StackFrame{
				{File: "app/user.rb", Line: 10, Function: "create"},
				{File: "app/user.rb", Line: 50, Function: "create"},
			},
			expected: "app/user.rb::create|app/user.rb::create",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeBacktraceForGrouping(tt.frames)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetProjectFrames(t *testing.T) {
	t.Run("Returns filtered backtrace when available", func(t *testing.T) {
		result := types.TestResult{
			FullBacktrace: []types.StackFrame{
				{File: "app/test.rb", Line: 10},
				{File: "/gems/rspec.rb", Line: 20},
			},
			FilteredBacktrace: []types.StackFrame{
				{File: "app/test.rb", Line: 10},
			},
		}

		frames := GetProjectFrames(result)
		assert.Len(t, frames, 1)
		assert.Equal(t, "app/test.rb", frames[0].File)
	})

	t.Run("Returns full backtrace when filtered is empty", func(t *testing.T) {
		result := types.TestResult{
			FullBacktrace: []types.StackFrame{
				{File: "app/test.rb", Line: 10},
			},
			FilteredBacktrace: []types.StackFrame{},
		}

		frames := GetProjectFrames(result)
		assert.Len(t, frames, 1)
		assert.Equal(t, "app/test.rb", frames[0].File)
	})
}

func TestCountFilteredFrames(t *testing.T) {
	results := []types.TestResult{
		{
			FullBacktrace: []types.StackFrame{
				{File: "app/test.rb", Line: 10},
				{File: "/gems/rspec.rb", Line: 20},
				{File: "/gems/another.rb", Line: 30},
			},
			FilteredBacktrace: []types.StackFrame{
				{File: "app/test.rb", Line: 10},
			},
		},
		{
			FullBacktrace: []types.StackFrame{
				{File: "app/another.rb", Line: 40},
				{File: "/gems/something.rb", Line: 50},
			},
			FilteredBacktrace: []types.StackFrame{
				{File: "app/another.rb", Line: 40},
			},
		},
	}

	total, filtered := CountFilteredFrames(results)

	assert.Equal(t, 5, total, "Should count all frames from full backtraces")
	assert.Equal(t, 2, filtered, "Should count all frames from filtered backtraces")
}
