package filesnippet

import (
	"fmt"
	"os"
	"strings"
)

// Line represents a single line from a file
type Line struct {
	Number   int
	Content  string
	IsCenter bool
}

// FileSnippet represents a range of lines extracted from a file
type FileSnippet struct {
	FilePath string
	Lines    []Line
}

// ExtractLines extracts a range of lines from a file centered on the given line number
// centerLine is 1-indexed. size specifies how many lines before and after to include.
// Returns a FileSnippet containing the file path and extracted lines.
func ExtractLines(filePath string, centerLine int, size int) (*FileSnippet, error) {
	if centerLine < 1 {
		return nil, fmt.Errorf("line number must be >= 1, got %d", centerLine)
	}
	if size < 0 {
		return nil, fmt.Errorf("size must be >= 0, got %d", size)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	if len(data) == 0 {
		if centerLine > 0 {
			return nil, fmt.Errorf("file is empty but centerLine is %d", centerLine)
		}
		return &FileSnippet{
			FilePath: filePath,
			Lines:    []Line{},
		}, nil
	}

	lines := strings.Split(string(data), "\n")
	totalLines := len(lines)

	if centerLine > totalLines {
		return nil, fmt.Errorf("line number %d exceeds file length %d", centerLine, totalLines)
	}

	startLine := centerLine - size
	if startLine < 1 {
		startLine = 1
	}

	endLine := centerLine + size
	if endLine > totalLines {
		endLine = totalLines
	}

	result := make([]Line, 0, endLine-startLine+1)
	for i := startLine; i <= endLine; i++ {
		lineContent := lines[i-1]
		result = append(result, Line{
			Number:   i,
			Content:  lineContent,
			IsCenter: i == centerLine,
		})
	}

	return &FileSnippet{
		FilePath: filePath,
		Lines:    result,
	}, nil
}
