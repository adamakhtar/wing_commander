package grouper

import (
	"testing"

	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestErrorLocationStrategy_GroupKey(t *testing.T) {
	strategy := NewErrorLocationStrategy()

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
				{File: "app/models/user.rb", Line: 42, Function: "create_user"},
			},
			expected: "app/models/user.rb:42",
		},
		{
			name: "Multiple frames - uses bottom frame",
			frames: []types.StackFrame{
				{File: "app/controllers/users_controller.rb", Line: 10, Function: "create"},
				{File: "app/services/user_service.rb", Line: 25, Function: "process"},
				{File: "app/models/user.rb", Line: 42, Function: "create_user"},
			},
			expected: "app/models/user.rb:42",
		},
		{
			name: "Different line numbers in same file",
			frames: []types.StackFrame{
				{File: "app/models/user.rb", Line: 10, Function: "validate"},
				{File: "app/models/user.rb", Line: 50, Function: "create_user"},
			},
			expected: "app/models/user.rb:50",
		},
		{
			name: "Frame without function name",
			frames: []types.StackFrame{
				{File: "app/models/user.rb", Line: 42},
			},
			expected: "app/models/user.rb:42",
		},
		{
			name: "Different files with same line number",
			frames: []types.StackFrame{
				{File: "app/models/user.rb", Line: 42, Function: "create_user"},
				{File: "app/models/product.rb", Line: 42, Function: "create_product"},
			},
			expected: "app/models/product.rb:42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strategy.GroupKey(tt.frames)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestErrorLocationStrategy_GroupingBehavior(t *testing.T) {
	strategy := NewErrorLocationStrategy()

	// Test that same bottom frame produces same group key
	frames1 := []types.StackFrame{
		{File: "app/controllers/users_controller.rb", Line: 10, Function: "create"},
		{File: "app/models/user.rb", Line: 42, Function: "create_user"},
	}

	frames2 := []types.StackFrame{
		{File: "app/services/user_service.rb", Line: 25, Function: "process"},
		{File: "app/models/user.rb", Line: 42, Function: "create_user"},
	}

	key1 := strategy.GroupKey(frames1)
	key2 := strategy.GroupKey(frames2)

	assert.Equal(t, key1, key2, "Same bottom frame should produce same group key")
	assert.Equal(t, "app/models/user.rb:42", key1)

	// Test that different line numbers produce different group keys
	frames3 := []types.StackFrame{
		{File: "app/models/user.rb", Line: 50, Function: "create_user"},
	}

	key3 := strategy.GroupKey(frames3)
	assert.NotEqual(t, key1, key3, "Different line numbers should produce different group keys")
	assert.Equal(t, "app/models/user.rb:50", key3)
}

func TestNewErrorLocationStrategy(t *testing.T) {
	strategy := NewErrorLocationStrategy()
	assert.NotNil(t, strategy)
	assert.IsType(t, &ErrorLocationStrategy{}, strategy)
}
