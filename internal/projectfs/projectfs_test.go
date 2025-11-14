package projectfs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/adamakhtar/wing_commander/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestProjectFS(t *testing.T, rootPath types.AbsPath) {
	t.Helper()
	err := InitProjectFS(rootPath, "")
	if err != nil {
		t.Fatalf("failed to initialize ProjectFS: %v", err)
	}
}

func teardownTestProjectFS(t *testing.T) {
	t.Helper()
	instance = nil
}

func TestInitProjectFS(t *testing.T) {
	t.Run("initializes singleton", func(t *testing.T) {
		defer teardownTestProjectFS(t)

		rootPath, err := types.NewAbsPath("/tmp/test-project")
		require.NoError(t, err)

		err = InitProjectFS(rootPath, "")
		require.NoError(t, err)

		fs := GetProjectFS()
		require.NotNil(t, fs)
		assert.Equal(t, rootPath, fs.RootPath)
	})
}

func TestGetProjectFS(t *testing.T) {
	t.Run("panics when not initialized", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic when ProjectFS not initialized")
			}
			instance = nil
		}()

		instance = nil
		GetProjectFS()
	})

	t.Run("returns instance when initialized", func(t *testing.T) {
		defer teardownTestProjectFS(t)

		rootPath, err := types.NewAbsPath("/tmp/test-project")
		require.NoError(t, err)

		err = InitProjectFS(rootPath, "")
		require.NoError(t, err)
		fs := GetProjectFS()
		assert.NotNil(t, fs)
		assert.Equal(t, rootPath, fs.RootPath)
	})
}

func TestProjectFS_Abs(t *testing.T) {
	tests := []struct {
		name     string
		rootPath string
		relPath  string
		want     string
	}{
		{
			name:     "simple relative path",
			rootPath: "/tmp/project",
			relPath:  "test/file.go",
			want:     "/tmp/project/test/file.go",
		},
		{
			name:     "relative path with cleanup",
			rootPath: "/tmp/project",
			relPath:  "test/../test/file.go",
			want:     "/tmp/project/test/file.go",
		},
		{
			name:     "current directory",
			rootPath: "/tmp/project",
			relPath:  ".",
			want:     "/tmp/project",
		},
		{
			name:     "nested path",
			rootPath: "/tmp/project",
			relPath:  "src/internal/types.go",
			want:     "/tmp/project/src/internal/types.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer teardownTestProjectFS(t)

			rootPath, err := types.NewAbsPath(tt.rootPath)
			require.NoError(t, err)
			setupTestProjectFS(t, rootPath)

			relPath, err := types.NewRelPath(tt.relPath)
			require.NoError(t, err)

			fs := GetProjectFS()
			got := fs.Abs(relPath)

			// Normalize the expected path
			wantAbs, err := filepath.Abs(tt.want)
			require.NoError(t, err)
			wantNormalized := filepath.Clean(wantAbs)

			assert.Equal(t, wantNormalized, got.String())
			assert.True(t, filepath.IsAbs(got.String()), "result should be absolute")
		})
	}
}

func TestProjectFS_Rel(t *testing.T) {
	tests := []struct {
		name     string
		rootPath string
		absPath  string
		want     string
		wantErr  bool
	}{
		{
			name:     "simple absolute path",
			rootPath: "/tmp/project",
			absPath:  "/tmp/project/test/file.go",
			want:     "test/file.go",
			wantErr:  false,
		},
		{
			name:     "file at root",
			rootPath: "/tmp/project",
			absPath:  "/tmp/project/file.go",
			want:     "file.go",
			wantErr:  false,
		},
		{
			name:     "nested path",
			rootPath: "/tmp/project",
			absPath:  "/tmp/project/src/internal/types.go",
			want:     "src/internal/types.go",
			wantErr:  false,
		},
		{
			name:     "path outside root",
			rootPath: "/tmp/project",
			absPath:  "/tmp/other/file.go",
			wantErr:  true,
		},
		{
			name:     "path above root",
			rootPath: "/tmp/project",
			absPath:  "/tmp/file.go",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer teardownTestProjectFS(t)

			rootPath, err := types.NewAbsPath(tt.rootPath)
			require.NoError(t, err)
			setupTestProjectFS(t, rootPath)

			absPath, err := types.NewAbsPath(tt.absPath)
			require.NoError(t, err)

			fs := GetProjectFS()
			got, err := fs.Rel(absPath)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, types.RelPath(""), got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got.String())
				assert.False(t, filepath.IsAbs(got.String()), "result should be relative")
			}
		})
	}
}

func TestProjectFS_Integration(t *testing.T) {
	t.Run("round trip conversion", func(t *testing.T) {
		defer teardownTestProjectFS(t)

		wd, err := os.Getwd()
		require.NoError(t, err)

		rootPath, err := types.NewAbsPath(wd)
		require.NoError(t, err)
		setupTestProjectFS(t, rootPath)

		fs := GetProjectFS()

		// Start with relative path
		relPath, err := types.NewRelPath("test/file.go")
		require.NoError(t, err)

		// Convert to absolute
		absPath := fs.Abs(relPath)
		assert.True(t, filepath.IsAbs(absPath.String()))

		// Convert back to relative
		relPath2, err := fs.Rel(absPath)
		require.NoError(t, err)
		assert.Equal(t, relPath.String(), relPath2.String())
	})

	t.Run("absolute to relative conversion", func(t *testing.T) {
		defer teardownTestProjectFS(t)

		rootPath, err := types.NewAbsPath("/tmp/project")
		require.NoError(t, err)
		setupTestProjectFS(t, rootPath)

		fs := GetProjectFS()

		// Create absolute path within project
		absPath, err := types.NewAbsPath("/tmp/project/src/file.go")
		require.NoError(t, err)

		// Convert to relative
		relPath, err := fs.Rel(absPath)
		require.NoError(t, err)
		assert.Equal(t, "src/file.go", relPath.String())
	})
}

func TestProjectFS_IsProjectFile(t *testing.T) {
	tests := []struct {
		name     string
		rootPath string
		absPath  string
		want     bool
	}{
		{
			name:     "file within project",
			rootPath: "/tmp/project",
			absPath:  "/tmp/project/test/file.go",
			want:     true,
		},
		{
			name:     "file at project root",
			rootPath: "/tmp/project",
			absPath:  "/tmp/project/file.go",
			want:     true,
		},
		{
			name:     "exact project root match",
			rootPath: "/tmp/project",
			absPath:  "/tmp/project",
			want:     true,
		},
		{
			name:     "file outside project",
			rootPath: "/tmp/project",
			absPath:  "/tmp/other/file.go",
			want:     false,
		},
		{
			name:     "file above project root",
			rootPath: "/tmp/project",
			absPath:  "/tmp/file.go",
			want:     false,
		},
		{
			name:     "completely different path",
			rootPath: "/tmp/project",
			absPath:  "/usr/lib/file.go",
			want:     false,
		},
		{
			name:     "nested subdirectory",
			rootPath: "/tmp/project",
			absPath:  "/tmp/project/src/internal/types.go",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer teardownTestProjectFS(t)

			rootPath, err := types.NewAbsPath(tt.rootPath)
			require.NoError(t, err)
			setupTestProjectFS(t, rootPath)

			absPath, err := types.NewAbsPath(tt.absPath)
			require.NoError(t, err)

			fs := GetProjectFS()
			got := fs.IsProjectFile(absPath)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProjectFS_IsTestFile(t *testing.T) {
	tests := []struct {
		name            string
		rootPath        string
		testFilePattern string
		absPath         string
		want            bool
	}{
		{
			name:            "matches test file pattern",
			rootPath:        "/tmp/project",
			testFilePattern: "test/**/*.rb",
			absPath:         "/tmp/project/test/models/user_test.rb",
			want:            true,
		},
		{
			name:            "matches nested test file pattern",
			rootPath:        "/tmp/project",
			testFilePattern: "test/**/*.rb",
			absPath:         "/tmp/project/test/unit/models/user_test.rb",
			want:            true,
		},
		{
			name:            "does not match non-test file",
			rootPath:        "/tmp/project",
			testFilePattern: "test/**/*.rb",
			absPath:         "/tmp/project/app/models/user.rb",
			want:            false,
		},
		{
			name:            "returns false when pattern not set",
			rootPath:        "/tmp/project",
			testFilePattern: "",
			absPath:         "/tmp/project/test/models/user_test.rb",
			want:            false,
		},
		{
			name:            "returns false for path outside project",
			rootPath:        "/tmp/project",
			testFilePattern: "test/**/*.rb",
			absPath:         "/tmp/other/test.rb",
			want:            false,
		},
		{
			name:            "matches custom pattern",
			rootPath:        "/tmp/project",
			testFilePattern: "spec/**/*_spec.rb",
			absPath:         "/tmp/project/spec/models/user_spec.rb",
			want:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer teardownTestProjectFS(t)

			rootPath, err := types.NewAbsPath(tt.rootPath)
			require.NoError(t, err)
			err = InitProjectFS(rootPath, tt.testFilePattern)
			require.NoError(t, err)

			absPath, err := types.NewAbsPath(tt.absPath)
			require.NoError(t, err)

			fs := GetProjectFS()
			got := fs.IsTestFile(absPath)

			assert.Equal(t, tt.want, got)
		})
	}
}
