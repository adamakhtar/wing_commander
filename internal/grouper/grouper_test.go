package grouper

import (
	"testing"

	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestNewGrouper(t *testing.T) {
	strategy := NewErrorLocationStrategy()
	grouper := NewGrouper(strategy)

	assert.NotNil(t, grouper)
	assert.Equal(t, strategy, grouper.strategy)
}

func TestGrouper_GroupFailures(t *testing.T) {
	strategy := NewErrorLocationStrategy()
	grouper := NewGrouper(strategy)

	t.Run("Empty results", func(t *testing.T) {
		results := []types.TestResult{}
		groups := grouper.GroupFailures(results)
		assert.Empty(t, groups)
	})

	t.Run("No failed tests", func(t *testing.T) {
		results := []types.TestResult{
			{Name: "Test 1", Status: types.StatusPass},
			{Name: "Test 2", Status: types.StatusSkip},
		}
		groups := grouper.GroupFailures(results)
		assert.Empty(t, groups)
	})

	t.Run("Single failure", func(t *testing.T) {
		results := []types.TestResult{
			{
				Name:         "Test 1",
				Status:       types.StatusFail,
				ErrorMessage: "Something went wrong",
				FilteredBacktrace: []types.StackFrame{
					{File: "app/models/user.rb", Line: 42, Function: "create_user"},
				},
			},
		}
		groups := grouper.GroupFailures(results)

		assert.Len(t, groups, 1)
		assert.Equal(t, "app/models/user.rb:42", groups[0].Hash)
		assert.Equal(t, "Something went wrong", groups[0].ErrorMessage)
		assert.Equal(t, 1, groups[0].Count)
		assert.Len(t, groups[0].Tests, 1)
		assert.Equal(t, "Test 1", groups[0].Tests[0].Name)
	})

	t.Run("Multiple failures with same bottom frame", func(t *testing.T) {
		results := []types.TestResult{
			{
				Name:         "Test 1",
				Status:       types.StatusFail,
				ErrorMessage: "Something went wrong",
				FilteredBacktrace: []types.StackFrame{
					{File: "app/controllers/users_controller.rb", Line: 10, Function: "create"},
					{File: "app/models/user.rb", Line: 42, Function: "create_user"},
				},
			},
			{
				Name:         "Test 2",
				Status:       types.StatusFail,
				ErrorMessage: "Another error",
				FilteredBacktrace: []types.StackFrame{
					{File: "app/services/user_service.rb", Line: 25, Function: "process"},
					{File: "app/models/user.rb", Line: 42, Function: "create_user"},
				},
			},
		}
		groups := grouper.GroupFailures(results)

		assert.Len(t, groups, 1)
		assert.Equal(t, "app/models/user.rb:42", groups[0].Hash)
		assert.Equal(t, 2, groups[0].Count)
		assert.Len(t, groups[0].Tests, 2)
	})

	t.Run("Multiple failures with different bottom frames", func(t *testing.T) {
		results := []types.TestResult{
			{
				Name:         "Test 1",
				Status:       types.StatusFail,
				ErrorMessage: "Error in user.rb",
				FilteredBacktrace: []types.StackFrame{
					{File: "app/models/user.rb", Line: 42, Function: "create_user"},
				},
			},
			{
				Name:         "Test 2",
				Status:       types.StatusFail,
				ErrorMessage: "Error in product.rb",
				FilteredBacktrace: []types.StackFrame{
					{File: "app/models/product.rb", Line: 50, Function: "create_product"},
				},
			},
		}
		groups := grouper.GroupFailures(results)

		assert.Len(t, groups, 2)

		// Groups should be sorted by count (descending)
		// Since both have count 1, order is not guaranteed, so check both possibilities
		groupHashes := []string{groups[0].Hash, groups[1].Hash}
		assert.Contains(t, groupHashes, "app/models/user.rb:42")
		assert.Contains(t, groupHashes, "app/models/product.rb:50")

		assert.Equal(t, 1, groups[0].Count)
		assert.Equal(t, 1, groups[1].Count)
	})

	t.Run("Mixed pass/fail tests", func(t *testing.T) {
		results := []types.TestResult{
			{
				Name:         "Test 1",
				Status:       types.StatusPass,
				ErrorMessage: "This should be ignored",
				FilteredBacktrace: []types.StackFrame{
					{File: "app/models/user.rb", Line: 42, Function: "create_user"},
				},
			},
			{
				Name:         "Test 2",
				Status:       types.StatusFail,
				ErrorMessage: "This should be grouped",
				FilteredBacktrace: []types.StackFrame{
					{File: "app/models/user.rb", Line: 42, Function: "create_user"},
				},
			},
		}
		groups := grouper.GroupFailures(results)

		assert.Len(t, groups, 1)
		assert.Equal(t, "app/models/user.rb:42", groups[0].Hash)
		assert.Equal(t, 1, groups[0].Count)
		assert.Len(t, groups[0].Tests, 1)
		assert.Equal(t, "Test 2", groups[0].Tests[0].Name)
	})

	t.Run("Failures with empty backtraces", func(t *testing.T) {
		results := []types.TestResult{
			{
				Name:         "Test 1",
				Status:       types.StatusFail,
				ErrorMessage: "No backtrace",
				FilteredBacktrace: []types.StackFrame{},
			},
		}
		groups := grouper.GroupFailures(results)

		// Should be empty since no valid group key can be generated
		assert.Empty(t, groups)
	})

	t.Run("Groups sorted by count", func(t *testing.T) {
		results := []types.TestResult{
			// Single failure at user.rb:42
			{
				Name:         "Test 1",
				Status:       types.StatusFail,
				ErrorMessage: "Single failure",
				FilteredBacktrace: []types.StackFrame{
					{File: "app/models/user.rb", Line: 42, Function: "create_user"},
				},
			},
			// Two failures at product.rb:50
			{
				Name:         "Test 2",
				Status:       types.StatusFail,
				ErrorMessage: "First product failure",
				FilteredBacktrace: []types.StackFrame{
					{File: "app/models/product.rb", Line: 50, Function: "create_product"},
				},
			},
			{
				Name:         "Test 3",
				Status:       types.StatusFail,
				ErrorMessage: "Second product failure",
				FilteredBacktrace: []types.StackFrame{
					{File: "app/models/product.rb", Line: 50, Function: "create_product"},
				},
			},
		}
		groups := grouper.GroupFailures(results)

		assert.Len(t, groups, 2)

		// First group should have count 2 (product.rb:50)
		assert.Equal(t, 2, groups[0].Count)
		assert.Equal(t, "app/models/product.rb:50", groups[0].Hash)

		// Second group should have count 1 (user.rb:42)
		assert.Equal(t, 1, groups[1].Count)
		assert.Equal(t, "app/models/user.rb:42", groups[1].Hash)
	})
}

func TestGrouper_CollectAllFrames(t *testing.T) {
	strategy := NewErrorLocationStrategy()
	grouper := NewGrouper(strategy)

	failedTests := []types.TestResult{
		{
			Name:         "Test 1",
			Status:       types.StatusFail,
			ErrorMessage: "Error 1",
			FilteredBacktrace: []types.StackFrame{
				{File: "app/models/user.rb", Line: 42, Function: "create_user"},
				{File: "app/services/user_service.rb", Line: 25, Function: "process"},
			},
		},
		{
			Name:         "Test 2",
			Status:       types.StatusFail,
			ErrorMessage: "Error 2",
			FilteredBacktrace: []types.StackFrame{
				{File: "app/models/product.rb", Line: 30, Function: "create_product"},
			},
		},
	}

	allFrames := grouper.collectAllFrames(failedTests)

	assert.Len(t, allFrames, 3)
	assert.Equal(t, "app/models/user.rb", allFrames[0].File)
	assert.Equal(t, "app/services/user_service.rb", allFrames[1].File)
	assert.Equal(t, "app/models/product.rb", allFrames[2].File)
}

func TestGrouper_GroupFailures_WithChangeDetection(t *testing.T) {
	strategy := NewErrorLocationStrategy()
	grouper := NewGrouper(strategy)

	results := []types.TestResult{
		{
			Name:         "Test 1",
			Status:       types.StatusFail,
			ErrorMessage: "Something went wrong",
			FilteredBacktrace: []types.StackFrame{
				{File: "app/models/user.rb", Line: 42, Function: "create_user"},
			},
		},
	}

	groups := grouper.GroupFailures(results)

	assert.Len(t, groups, 1)
	assert.Equal(t, "app/models/user.rb:42", groups[0].Hash)

	// Note: Change detection will be tested in integration tests
	// since it requires actual git commands or mocking
}

func TestFilterFailedTests(t *testing.T) {
	results := []types.TestResult{
		{Name: "Test 1", Status: types.StatusPass},
		{Name: "Test 2", Status: types.StatusFail},
		{Name: "Test 3", Status: types.StatusSkip},
		{Name: "Test 4", Status: types.StatusFail},
	}

	failed := filterFailedTests(results)

	assert.Len(t, failed, 2)
	assert.Equal(t, "Test 2", failed[0].Name)
	assert.Equal(t, "Test 4", failed[1].Name)
}
