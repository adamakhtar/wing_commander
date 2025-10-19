package editor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Editor handles opening files in external editors
type Editor struct {
	editorCommand string
}

// NewEditor creates a new Editor instance
func NewEditor() *Editor {
	return &Editor{
		editorCommand: getDefaultEditor(),
	}
}

// OpenFile opens a file in the configured editor at the specified line
func (e *Editor) OpenFile(filePath string, line int) error {
	// Validate inputs
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}
	if line < 1 {
		return fmt.Errorf("line number must be positive, got %d", line)
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %w", filePath, err)
	}

	// Check if file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", absPath)
	} else if err != nil {
		return fmt.Errorf("failed to check if file exists %s: %w", absPath, err)
	}

	// Check if editor command exists
	if _, err := exec.LookPath(e.editorCommand); err != nil {
		return fmt.Errorf("editor command not found: %s (error: %w)", e.editorCommand, err)
	}

	// Build command arguments
	var args []string
	
	// Handle different editor commands
	switch {
	case strings.Contains(e.editorCommand, "vim") || strings.Contains(e.editorCommand, "nvim"):
		args = []string{absPath, fmt.Sprintf("+%d", line)}
	case strings.Contains(e.editorCommand, "code"):
		args = []string{"--goto", fmt.Sprintf("%s:%d", absPath, line)}
	case strings.Contains(e.editorCommand, "subl"):
		args = []string{fmt.Sprintf("%s:%d", absPath, line)}
	case strings.Contains(e.editorCommand, "atom"):
		args = []string{fmt.Sprintf("%s:%d", absPath, line)}
	case strings.Contains(e.editorCommand, "emacs"):
		args = []string{fmt.Sprintf("+%d", line), absPath}
	default:
		// Generic fallback - try to pass line number as argument
		args = []string{absPath, fmt.Sprintf("+%d", line)}
	}

	// Execute the editor command
	cmd := exec.Command(e.editorCommand, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open file %s:%d in editor %s: %w", absPath, line, e.editorCommand, err)
	}

	return nil
}

// getDefaultEditor determines the default editor to use
func getDefaultEditor() string {
	// Check EDITOR environment variable first
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}

	// Check VISUAL environment variable
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}

	// Check for common editors in PATH
	commonEditors := []string{
		"code",    // VS Code
		"subl",    // Sublime Text
		"atom",    // Atom
		"nvim",    // Neovim
		"vim",     // Vim
		"emacs",   // Emacs
		"nano",    // Nano
	}

	for _, editor := range commonEditors {
		if _, err := exec.LookPath(editor); err == nil {
			return editor
		}
	}

	// Fallback to vim (usually available on most systems)
	return "vim"
}

// SetEditor sets a custom editor command
func (e *Editor) SetEditor(editorCommand string) {
	e.editorCommand = editorCommand
}

// GetEditor returns the current editor command
func (e *Editor) GetEditor() string {
	return e.editorCommand
}
