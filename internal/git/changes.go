package git

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/adamakhtar/wing_commander/internal/types"
)

// ChangeDetector handles detection of line-level changes in git
type ChangeDetector struct {
	hunkRegex *regexp.Regexp
}

// NewChangeDetector creates a new ChangeDetector
func NewChangeDetector() *ChangeDetector {
	return &ChangeDetector{
		hunkRegex: regexp.MustCompile(`@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@`),
	}
}

// DetectChanges analyzes all stack frames and assigns change intensities
// Returns a map of file paths to their changed line numbers for each change type
func (cd *ChangeDetector) DetectChanges(frames []types.StackFrame) map[string]*FileChanges {
	// Get all unique files from frames
	files := make(map[string]bool)
	for _, frame := range frames {
		files[frame.File] = true
	}

	// Detect changes for each file
	fileChanges := make(map[string]*FileChanges)
	for file := range files {
		changes := cd.detectFileChanges(file)
		if changes != nil {
			fileChanges[file] = changes
		}
	}

	return fileChanges
}

// FileChanges represents the changed lines for a specific file
type FileChanges struct {
	UncommittedLines    map[int]bool // Line numbers with uncommitted changes
	LastCommitLines     map[int]bool // Line numbers changed in last commit
	PreviousCommitLines map[int]bool // Line numbers changed in previous commit
}

// detectFileChanges detects changes for a specific file
func (cd *ChangeDetector) detectFileChanges(filePath string) *FileChanges {
	changes := &FileChanges{
		UncommittedLines:    make(map[int]bool),
		LastCommitLines:     make(map[int]bool),
		PreviousCommitLines: make(map[int]bool),
	}

	// Detect uncommitted changes (intensity 3)
	uncommittedLines, err := cd.getUncommittedChanges(filePath)
	if err == nil {
		for _, line := range uncommittedLines {
			changes.UncommittedLines[line] = true
		}
	}

	// Detect last commit changes (intensity 2)
	lastCommitLines, err := cd.getLastCommitChanges(filePath)
	if err == nil {
		for _, line := range lastCommitLines {
			changes.LastCommitLines[line] = true
		}
	}

	// Detect previous commit changes (intensity 1)
	previousCommitLines, err := cd.getPreviousCommitChanges(filePath)
	if err == nil {
		for _, line := range previousCommitLines {
			changes.PreviousCommitLines[line] = true
		}
	}

	return changes
}

// getUncommittedChanges gets line numbers with uncommitted changes
func (cd *ChangeDetector) getUncommittedChanges(filePath string) ([]int, error) {
	cmd := exec.Command("git", "diff", "--unified=0", filePath)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return cd.parseDiffOutput(string(output)), nil
}

// getLastCommitChanges gets line numbers changed in the last commit
func (cd *ChangeDetector) getLastCommitChanges(filePath string) ([]int, error) {
	cmd := exec.Command("git", "diff", "HEAD~1", "--unified=0", filePath)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return cd.parseDiffOutput(string(output)), nil
}

// getPreviousCommitChanges gets line numbers changed in the commit before last
func (cd *ChangeDetector) getPreviousCommitChanges(filePath string) ([]int, error) {
	cmd := exec.Command("git", "diff", "HEAD~2", "HEAD~1", "--unified=0", filePath)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return cd.parseDiffOutput(string(output)), nil
}

// parseDiffOutput parses unified diff output to extract changed line numbers
func (cd *ChangeDetector) parseDiffOutput(diffOutput string) []int {
	// Pre-allocate slice with reasonable capacity
	changedLines := make([]int, 0, 100)

	lines := strings.Split(diffOutput, "\n")
	for _, line := range lines {
		matches := cd.hunkRegex.FindStringSubmatch(line)
		if len(matches) >= 4 {
			// Extract new line start and count
			newStart, _ := strconv.Atoi(matches[3])
			newCount := 1
			if matches[4] != "" {
				newCount, _ = strconv.Atoi(matches[4])
			}

			// Add all lines in the range
			for i := 0; i < newCount; i++ {
				changedLines = append(changedLines, newStart+i)
			}
		}
	}

	return changedLines
}

// AssignChangeIntensities assigns change intensities to stack frames based on detected changes
func (cd *ChangeDetector) AssignChangeIntensities(frames []types.StackFrame, fileChanges map[string]*FileChanges) {
	for i := range frames {
		frame := &frames[i]
		changes, exists := fileChanges[frame.File]
		if !exists {
			continue
		}

		// Assign highest intensity (uncommitted > last commit > previous commit)
		if changes.UncommittedLines[frame.Line] {
			frame.ChangeIntensity = 3
			frame.ChangeReason = "uncommitted"
		} else if changes.LastCommitLines[frame.Line] {
			frame.ChangeIntensity = 2
			frame.ChangeReason = "last_commit"
		} else if changes.PreviousCommitLines[frame.Line] {
			frame.ChangeIntensity = 1
			frame.ChangeReason = "previous_commit"
		}
	}
}
