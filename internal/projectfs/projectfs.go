package projectfs

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/gobwas/glob"
)

var (
	instance *ProjectFS
)

// ProjectFS manages the project root path and provides path conversions
type ProjectFS struct {
	RootPath       types.AbsPath
	testFileGlob glob.Glob
}

// InitProjectFS initializes the singleton with the project root path and optional test file pattern
func InitProjectFS(rootPath types.AbsPath, testFileGlob string) error {
	fs := &ProjectFS{
		RootPath: rootPath,
	}

	if testFileGlob != "" {
		// Normalize pattern to forward slashes for cross-platform compatibility.
		// The glob library expects forward slashes, and paths are normalized when matching.
		pattern := filepath.ToSlash(testFileGlob)
		compiled, err := glob.Compile(pattern)
		if err != nil {
			return fmt.Errorf("failed to compile test file pattern %q: %w", testFileGlob, err)
		}
		fs.testFileGlob = compiled
	}

	instance = fs
	return nil
}

// GetProjectFS returns the singleton instance
func GetProjectFS() *ProjectFS {
	if instance == nil {
		panic("ProjectFS not initialized. Call InitProjectFS() first.")
	}
	return instance
}

// Abs converts a relative path to an absolute path using RootPath
func (fs *ProjectFS) Abs(rel types.RelPath) types.AbsPath {
	fullPath := filepath.Join(fs.RootPath.String(), rel.String())
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		panic(fmt.Sprintf("failed to get absolute path: %v", err))
	}
	return types.AbsPath(absPath)
}

// Rel converts an absolute path to a relative path using RootPath.
// Returns an error if the path is outside the project root.
func (fs *ProjectFS) Rel(abs types.AbsPath) (types.RelPath, error) {
	rel, err := filepath.Rel(fs.RootPath.String(), abs.String())
	if err != nil {
		return types.RelPath(""), fmt.Errorf("failed to get relative path: %w", err)
	}

	if strings.HasPrefix(rel, "..") {
		return types.RelPath(""), fmt.Errorf("path is outside project root")
	}

	relPath, err := types.NewRelPath(rel)
	if err != nil {
		return types.RelPath(""), fmt.Errorf("failed to create RelPath: %w", err)
	}

	return relPath, nil
}

// IsProjectFile checks if the given absolute path is within the project root.
// Uses Rel() for cross-platform path comparison.
func (fs *ProjectFS) IsProjectFile(absPath types.AbsPath) bool {
	_, err := fs.Rel(absPath)
	return err == nil
}

// IsTestFile checks if the given absolute path matches the test file pattern.
// Converts the absolute path to a relative path and matches against the compiled glob pattern.
func (fs *ProjectFS) IsTestFile(absPath types.AbsPath) bool {
	if fs.testFileGlob == nil || absPath == "" {
		return false
	}

	rel, err := fs.Rel(absPath)
	if err != nil {
		return false
	}

	return rel.MatchGlob(fs.testFileGlob)
}
