package projectfs

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/adamakhtar/wing_commander/internal/types"
)

var (
	instance *ProjectFS
)

// ProjectFS manages the project root path and provides path conversions
type ProjectFS struct {
	RootPath types.AbsPath
}

// InitProjectFS initializes the singleton with the project root path
func InitProjectFS(rootPath types.AbsPath) {
	instance = &ProjectFS{
		RootPath: rootPath,
	}
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

// Rel converts an absolute path to a relative path using RootPath
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
