package git

import (
	"testing"

	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestNewChangeDetector(t *testing.T) {
	detector := NewChangeDetector()
	assert.NotNil(t, detector)
}

func TestChangeDetector_ParseDiffOutput(t *testing.T) {
	detector := NewChangeDetector()

	tests := []struct {
		name     string
		diffOutput string
		expected []int
	}{
		{
			name: "Single hunk",
			diffOutput: `@@ -10,3 +10,4 @@
+new line
 unchanged line
 unchanged line
 unchanged line`,
			expected: []int{10, 11, 12, 13},
		},
		{
			name: "Multiple hunks",
			diffOutput: `@@ -5,2 +5,3 @@
+added line
 unchanged
 unchanged
@@ -20,1 +21,2 @@
 unchanged
+another added line`,
			expected: []int{5, 6, 7, 21, 22},
		},
		{
			name: "No count specified",
			diffOutput: `@@ -10 +10,2 @@
+new line
+another new line`,
			expected: []int{10, 11},
		},
		{
			name: "Empty diff",
			diffOutput: "",
			expected: []int{},
		},
		{
			name: "Complex diff with context",
			diffOutput: `diff --git a/test.go b/test.go
index 1234567..abcdefg 100644
--- a/test.go
+++ b/test.go
@@ -15,3 +15,4 @@ func test() {
     return true
+    // new comment
 }
@@ -25,1 +26,2 @@ func another() {
     return false
+    // another comment
+    return true`,
			expected: []int{15, 16, 17, 18, 26, 27},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.parseDiffOutput(tt.diffOutput)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestChangeDetector_AssignChangeIntensities(t *testing.T) {
	detector := NewChangeDetector()

	userPath, _ := types.NewAbsPath("/app/models/user.rb")
	productPath, _ := types.NewAbsPath("/app/models/product.rb")
	frames := []types.StackFrame{
		{FilePath: userPath, Line: 42, Function: "create_user"},
		{FilePath: userPath, Line: 50, Function: "validate"},
		{FilePath: productPath, Line: 30, Function: "create_product"},
		{FilePath: userPath, Line: 60, Function: "save"},
	}

	fileChanges := map[string]*FileChanges{
		userPath.String(): {
			UncommittedLines:    map[int]bool{42: true},
			LastCommitLines:     map[int]bool{50: true},
			PreviousCommitLines: map[int]bool{60: true},
		},
		productPath.String(): {
			UncommittedLines:    map[int]bool{},
			LastCommitLines:     map[int]bool{},
			PreviousCommitLines: map[int]bool{},
		},
	}

	detector.AssignChangeIntensities(frames, fileChanges)

	// Check intensities and reasons
	assert.Equal(t, 3, frames[0].ChangeIntensity)
	assert.Equal(t, "uncommitted", frames[0].ChangeReason)

	assert.Equal(t, 2, frames[1].ChangeIntensity)
	assert.Equal(t, "last_commit", frames[1].ChangeReason)

	assert.Equal(t, 0, frames[2].ChangeIntensity)
	assert.Equal(t, "", frames[2].ChangeReason)

	assert.Equal(t, 1, frames[3].ChangeIntensity)
	assert.Equal(t, "previous_commit", frames[3].ChangeReason)
}

func TestChangeDetector_AssignChangeIntensities_Priority(t *testing.T) {
	detector := NewChangeDetector()

	// Frame with line that has changes in multiple commits
	userPath, _ := types.NewAbsPath("/app/models/user.rb")
	frames := []types.StackFrame{
		{FilePath: userPath, Line: 42, Function: "create_user"},
	}

	fileChanges := map[string]*FileChanges{
		userPath.String(): {
			UncommittedLines:    map[int]bool{42: true},
			LastCommitLines:     map[int]bool{42: true},
			PreviousCommitLines: map[int]bool{42: true},
		},
	}

	detector.AssignChangeIntensities(frames, fileChanges)

	// Should get highest priority (uncommitted)
	assert.Equal(t, 3, frames[0].ChangeIntensity)
	assert.Equal(t, "uncommitted", frames[0].ChangeReason)
}

func TestChangeDetector_DetectChanges(t *testing.T) {
	detector := NewChangeDetector()

	userPath, _ := types.NewAbsPath("/app/models/user.rb")
	productPath, _ := types.NewAbsPath("/app/models/product.rb")
	frames := []types.StackFrame{
		{FilePath: userPath, Line: 42, Function: "create_user"},
		{FilePath: productPath, Line: 30, Function: "create_product"},
		{FilePath: userPath, Line: 50, Function: "validate"},
	}

	// This test would require actual git commands, so we'll test the structure
	fileChanges := detector.DetectChanges(frames)

	// Should have entries for both files
	assert.Contains(t, fileChanges, userPath.String())
	assert.Contains(t, fileChanges, productPath.String())

	// Each file should have the three change types
	userChanges := fileChanges[userPath.String()]
	assert.NotNil(t, userChanges)
	assert.NotNil(t, userChanges.UncommittedLines)
	assert.NotNil(t, userChanges.LastCommitLines)
	assert.NotNil(t, userChanges.PreviousCommitLines)
}

func TestFileChanges_Structure(t *testing.T) {
	changes := &FileChanges{
		UncommittedLines:    map[int]bool{42: true, 50: true},
		LastCommitLines:     map[int]bool{30: true},
		PreviousCommitLines: map[int]bool{60: true},
	}

	// Test uncommitted lines
	assert.True(t, changes.UncommittedLines[42])
	assert.True(t, changes.UncommittedLines[50])
	assert.False(t, changes.UncommittedLines[30])

	// Test last commit lines
	assert.True(t, changes.LastCommitLines[30])
	assert.False(t, changes.LastCommitLines[42])

	// Test previous commit lines
	assert.True(t, changes.PreviousCommitLines[60])
	assert.False(t, changes.PreviousCommitLines[42])
}

func TestStackFrame_NewFields(t *testing.T) {
	absPath, _ := types.NewAbsPath("/test.rb")
	frame := types.NewStackFrame(absPath, 42, "test_function")

	assert.Equal(t, absPath, frame.FilePath)
	assert.Equal(t, 42, frame.Line)
	assert.Equal(t, "test_function", frame.Function)
	assert.Equal(t, 0, frame.ChangeIntensity)
	assert.Equal(t, "", frame.ChangeReason)
}

func TestStackFrame_ManualCreation(t *testing.T) {
	absPath, _ := types.NewAbsPath("/app/models/user.rb")
	frame := types.StackFrame{
		FilePath:        absPath,
		Line:            42,
		Function:        "create_user",
		ChangeIntensity: 3,
		ChangeReason:    "uncommitted",
	}

	assert.Equal(t, absPath, frame.FilePath)
	assert.Equal(t, 42, frame.Line)
	assert.Equal(t, "create_user", frame.Function)
	assert.Equal(t, 3, frame.ChangeIntensity)
	assert.Equal(t, "uncommitted", frame.ChangeReason)
}
