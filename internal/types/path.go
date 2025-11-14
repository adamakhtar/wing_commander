package types

import (
	"fmt"
	"path/filepath"
)

// AbsPath represents an absolute file path
type AbsPath string

// NewAbsPath creates a new AbsPath from a string, validating it's absolute
func NewAbsPath(path string) (AbsPath, error) {
	if path == "" {
		return AbsPath(""), fmt.Errorf("path cannot be empty")
	}

	cleaned := filepath.Clean(path)

	if filepath.IsAbs(cleaned) {
		return AbsPath(cleaned), nil
	}

	absPath, err := filepath.Abs(cleaned)
	if err != nil {
		return AbsPath(""), fmt.Errorf("failed to get absolute path for %q: %w", path, err)
	}

	return AbsPath(absPath), nil
}

func (p AbsPath) String() string {
	return string(p)
}

// RelPath represents a relative file path
type RelPath string

// NewRelPath creates a new RelPath from a string, validating it's relative
func NewRelPath(path string) (RelPath, error) {
	if path == "" {
		return RelPath(""), fmt.Errorf("path cannot be empty")
	}

	cleaned := filepath.Clean(path)

	if filepath.IsAbs(cleaned) {
		return RelPath(""), fmt.Errorf("path %q is absolute, expected relative path", path)
	}

	return RelPath(cleaned), nil
}

func (p RelPath) String() string {
	return string(p)
}
