package backtrace

import (
	"testing"

	"github.com/adamakhtar/wing_commander/internal/config"
	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestNewNormalizer(t *testing.T) {
	cfg := &config.Config{
		ProjectPath: "/path/to/project",
	}

	normalizer := NewNormalizer(cfg)

	assert.NotNil(t, normalizer)
	assert.Equal(t, "/path/to/project", normalizer.projectPath)
}

func TestNormalizeTestResults(t *testing.T) {
	normalizer := &Normalizer{
		projectPath: "/path/to/project",
	}

	results := []types.TestResult{
		{
			GroupName: "Test 1",
			Status:    types.StatusFail,
			FullBacktrace: []types.StackFrame{
				{File: "/path/to/project/app/test.rb", Line: 10},
				{File: "/gems/rspec.rb", Line: 20},
			},
		},
		{
			GroupName: "Test 2",
			Status:    types.StatusFail,
			FullBacktrace: []types.StackFrame{
				{File: "/path/to/project/app/another.rb", Line: 30},
			},
		},
	}

	normalized := normalizer.NormalizeTestResults(results)

	assert.Len(t, normalized, 2)
	assert.Len(t, normalized[0].FilteredBacktrace, 1)
	assert.Equal(t, "/path/to/project/app/test.rb", normalized[0].FilteredBacktrace[0].File)
	assert.Len(t, normalized[1].FilteredBacktrace, 1)
	assert.Equal(t, "/path/to/project/app/another.rb", normalized[1].FilteredBacktrace[0].File)
}

func TestShouldExclude(t *testing.T) {
	t.Run("Excludes frames not starting with project path", func(t *testing.T) {
		normalizer := &Normalizer{
			projectPath: "/path/to/project",
		}

		frame := types.StackFrame{File: "/gems/rspec.rb", Line: 20}
		assert.True(t, normalizer.shouldExclude(frame))
	})

	t.Run("Includes frames starting with project path", func(t *testing.T) {
		normalizer := &Normalizer{
			projectPath: "/path/to/project",
		}

		frame := types.StackFrame{File: "/path/to/project/app/test.rb", Line: 10}
		assert.False(t, normalizer.shouldExclude(frame))
	})

	t.Run("Empty project path includes all frames", func(t *testing.T) {
		normalizer := &Normalizer{
			projectPath: "",
		}

		frame := types.StackFrame{File: "/gems/rspec.rb", Line: 20}
		assert.False(t, normalizer.shouldExclude(frame))
	})

	t.Run("Includes frames with exact project path match", func(t *testing.T) {
		normalizer := &Normalizer{
			projectPath: "/path/to/project",
		}

		frame := types.StackFrame{File: "/path/to/project", Line: 10}
		assert.False(t, normalizer.shouldExclude(frame))
	})
}
