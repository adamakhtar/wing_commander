package types

import (
	"fmt"
	"path/filepath"

	"github.com/gobwas/glob"
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

	return AbsPath(""), fmt.Errorf("path %q is relative, expected absolute path", path)
}

func (p AbsPath) String() string {
	return string(p)
}

// MatchGlob checks if the AbsPath matches the given glob pattern.
// The path is normalized (forward slashes) before matching to ensure
// cross-platform compatibility. Matching is case-sensitive.
func (p AbsPath) MatchGlob(g glob.Glob) bool {
	if g == nil || p == "" {
		return false
	}
	normalized := filepath.ToSlash(string(p))
	return g.Match(normalized)
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

// MatchGlob checks if the RelPath matches the given glob pattern.
// The path is normalized (forward slashes) before matching to ensure
// cross-platform compatibility. Matching is case-sensitive.
func (p RelPath) MatchGlob(g glob.Glob) bool {
	if g == nil || p == "" {
		return false
	}
	normalized := filepath.ToSlash(string(p))
	return g.Match(normalized)
}
