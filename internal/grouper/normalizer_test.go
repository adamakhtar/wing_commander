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
			GroupName:   "Test 1",
			TestCaseName: "",
			Status: types.StatusFail,
			FullBacktrace: []types.StackFrame{
				{File: "app/test.rb", Line: 10},
				{File: "/gems/rspec.rb", Line: 20},
			},
		},
		{
			GroupName:   "Test 2",
			TestCaseName: "",
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
