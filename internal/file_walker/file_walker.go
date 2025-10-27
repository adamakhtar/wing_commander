package filewalker

import (
	"os"
	"path/filepath"

	"github.com/gobwas/glob"
)


func FileEntriesRecursive(path string, includePatterns []string, excludePatterns []string) []string {
	files := []string{}
	filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Convert to relative path
		relPath, err := filepath.Rel(path, filePath)
		if err != nil {
			return err
		}

		// Add trailing slash for directories to aid user in differentiating between files and directories
		if info.IsDir() {
			relPath = relPath + string(filepath.Separator)
		}

		if shouldExclude(relPath, excludePatterns) {
			return nil
		}

		if shouldInclude(relPath, includePatterns) {
			files = append(files, relPath)
		}

		return nil
	})

	return files
}

func shouldExclude(filePath string, excludePatterns []string) bool {
	if len(excludePatterns) == 0 {
		return false
	}

	for _, pattern := range excludePatterns {
		g, err := glob.Compile(pattern)
		if err != nil {
			continue
		}
		if g.Match(filePath) {
			return true
		}
		if g.Match(filepath.Base(filePath)) {
			return true
		}
	}
	return false
}

func shouldInclude(filePath string, includePatterns []string) bool {
	if len(includePatterns) == 0 {
		return true
	}

	for _, pattern := range includePatterns {
		g, err := glob.Compile(pattern)
		if err != nil {
			continue
		}
		if g.Match(filePath) {
			return true
		}
		if g.Match(filepath.Base(filePath)) {
			return true
		}
	}
	return false
}