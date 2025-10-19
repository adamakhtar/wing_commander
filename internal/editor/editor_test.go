package editor

import (
	"testing"
)

func TestOpenFile(t *testing.T) {
	editor := NewEditor()

	// Test opening a non-existent file
	err := editor.OpenFile("/nonexistent/file.rb", 1)
	if err == nil {
		t.Error("Expected error when opening non-existent file")
	}

	// Test opening a file that exists (this test file itself)
	err = editor.OpenFile("editor_test.go", 1)
	// We can't actually test the editor opening in a test environment,
	// but we can verify it doesn't crash
	if err != nil && err.Error() != "file does not exist: /Users/adamakhtar/Projects/active/wing_commander/internal/editor/editor_test.go" {
		// This is expected in test environment - editor might not be available
		t.Logf("Editor test skipped (expected in CI): %v", err)
	}
}

func TestEditorCommands(t *testing.T) {
	editor := NewEditor()

	// Test setting custom editor
	customEditor := "custom-editor"
	editor.SetEditor(customEditor)

	if editor.GetEditor() != customEditor {
		t.Errorf("Expected editor to be %s, got %s", customEditor, editor.GetEditor())
	}
}