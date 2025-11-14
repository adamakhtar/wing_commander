package backtrace

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/adamakhtar/wing_commander/internal/projectfs"
	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/charmbracelet/log"
)

// Backtrace represents a collection of stack frames
type Backtrace struct {
	Frames []types.StackFrame
}

// NewBacktrace creates a new empty Backtrace
func NewBacktrace() Backtrace {
	return Backtrace{
		Frames: []types.StackFrame{},
	}
}

// Append parses a frame string and adds it to the backtrace.
// If the filepath is relative, it converts it to absolute using ProjectFS.
// Invalid frame strings result in a minimal StackFrame with a warning logged.
func (b *Backtrace) Append(frameStr string) {
	if frameStr == "" {
		log.Warn("empty frame string, creating minimal StackFrame")
		b.Frames = append(b.Frames, types.StackFrame{})
		return
	}

	frame := b.parseStackFrame(frameStr)
	b.Frames = append(b.Frames, frame)
}

// parseStackFrame parses a backtrace frame string into a StackFrame.
// Common formats:
// - "app/models/user.rb:42:in `create_user'"
// - "app/models/user.rb:42"
// - "File \"app/models/user.rb\", line 42, in create_user"
func (b *Backtrace) parseStackFrame(frameStr string) types.StackFrame {
	// Handle Python format first
	if strings.HasPrefix(frameStr, "File \"") {
		absPath, _ := types.NewAbsPath(frameStr)
		return types.StackFrame{
			FilePath: absPath,
			Line:     0,
			Function: "",
		}
	}

	parts := strings.Split(frameStr, ":")
	if len(parts) < 2 {
		absPath, err := b.convertPath(frameStr)
		if err != nil {
			log.Warn("failed to parse frame string", "frame", frameStr, "error", err)
			// Return minimal StackFrame with empty path on error
			return types.StackFrame{}
		}
		return types.StackFrame{FilePath: absPath}
	}

	file := parts[0]

	// Try to extract line number
	var line int
	var function string

	if len(parts) >= 2 {
		// Parse line number
		if _, err := fmt.Sscanf(parts[1], "%d", &line); err != nil {
			absPath, err := b.convertPath(file)
			if err != nil {
				log.Warn("failed to parse frame string", "frame", frameStr, "error", err)
			}
			return types.StackFrame{FilePath: absPath}
		}
	}

	// Try to extract function name
	if len(parts) >= 3 {
		funcPart := parts[2]
		// Remove "in `" and "`" wrapper, or "in '" and "'" wrapper
		if strings.HasPrefix(funcPart, "in `") && strings.HasSuffix(funcPart, "'") {
			function = funcPart[4 : len(funcPart)-1]
		} else if strings.HasPrefix(funcPart, "in '") && strings.HasSuffix(funcPart, "'") {
			function = funcPart[4 : len(funcPart)-1]
		} else if strings.HasPrefix(funcPart, "in ") {
			function = funcPart[3:]
		}
	}

	absPath, err := b.convertPath(file)
	if err != nil {
		log.Warn("failed to convert path in frame string", "frame", frameStr, "error", err)
	}

	return types.StackFrame{
		FilePath: absPath,
		Line:     line,
		Function: function,
	}
}

// convertPath converts a file path string to AbsPath.
// If the path is relative, it uses ProjectFS to convert it to absolute.
// If the path is already absolute, it uses NewAbsPath directly.
func (b *Backtrace) convertPath(path string) (types.AbsPath, error) {
	cleaned := filepath.Clean(path)

	if filepath.IsAbs(cleaned) {
		return types.NewAbsPath(cleaned)
	}

	// Relative path - convert using ProjectFS
	fs := projectfs.GetProjectFS()
	relPath, err := types.NewRelPath(cleaned)
	if err != nil {
		return types.AbsPath(""), fmt.Errorf("failed to create RelPath: %w", err)
	}

	return fs.Abs(relPath), nil
}

// FilterProjectStackFramesOnly returns a new Backtrace containing only
// stack frames whose file path is within the project's root path.
func (b Backtrace) FilterProjectStackFramesOnly() Backtrace {
	filtered := NewBacktrace()
	fs := projectfs.GetProjectFS()

	for _, frame := range b.Frames {
		if fs.IsProjectFile(frame.FilePath) {
			filtered.Frames = append(filtered.Frames, frame)
		}
	}

	return filtered
}

// AllStackFrames returns all stack frames in the backtrace
func (b Backtrace) AllStackFrames() []types.StackFrame {
	return b.Frames
}
