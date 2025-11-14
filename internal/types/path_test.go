package types

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAbsPath(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "empty path",
			input:   "",
			wantErr: true,
		},
		{
			name:    "absolute unix path",
			input:   "/tmp/test",
			want:    "/tmp/test",
			wantErr: false,
		},
		{
			name:    "absolute path with cleanup",
			input:   "/tmp/../tmp/test",
			want:    "/tmp/test",
			wantErr: false,
		},
		{
			name:    "relative path",
			input:   "test",
			wantErr: true,
		},
		{
			name:    "current directory",
			input:   ".",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAbsPath(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, AbsPath(""), got)
			} else {
				require.NoError(t, err)
				assert.True(t, filepath.IsAbs(got.String()), "result should be absolute: %s", got.String())
				if tt.want != "" {
					assert.Equal(t, tt.want, got.String())
				}
			}
		})
	}
}

func TestAbsPath_String(t *testing.T) {
	t.Run("returns underlying string", func(t *testing.T) {
		path, err := NewAbsPath("/tmp/test")
		require.NoError(t, err)
		assert.Equal(t, "/tmp/test", path.String())
	})
}

func TestNewRelPath(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "empty path",
			input:   "",
			wantErr: true,
		},
		{
			name:    "relative path",
			input:   "test/file",
			want:    "test/file",
			wantErr: false,
		},
		{
			name:    "relative path with cleanup",
			input:   "test/../test/file",
			want:    "test/file",
			wantErr: false,
		},
		{
			name:    "absolute path",
			input:   "/tmp/test",
			wantErr: true,
		},
		{
			name:    "current directory",
			input:   ".",
			want:    ".",
			wantErr: false,
		},
		{
			name:    "parent directory",
			input:   "..",
			want:    "..",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRelPath(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, RelPath(""), got)
			} else {
				require.NoError(t, err)
				assert.False(t, filepath.IsAbs(got.String()), "result should be relative: %s", got.String())
				if tt.want != "" {
					assert.Equal(t, tt.want, got.String())
				}
			}
		})
	}
}

func TestRelPath_String(t *testing.T) {
	t.Run("returns underlying string", func(t *testing.T) {
		path, err := NewRelPath("test/file")
		require.NoError(t, err)
		assert.Equal(t, "test/file", path.String())
	})
}

func TestAbsPath_Integration(t *testing.T) {
	t.Run("real absolute path conversion", func(t *testing.T) {
		wd, err := os.Getwd()
		require.NoError(t, err)

		absPath, err := NewAbsPath(wd)
		require.NoError(t, err)
		assert.Equal(t, wd, absPath.String())
	})

	t.Run("relative path returns error", func(t *testing.T) {
		relPath := "test"
		absPath, err := NewAbsPath(relPath)
		assert.Error(t, err)
		assert.Equal(t, AbsPath(""), absPath)
	})
}
