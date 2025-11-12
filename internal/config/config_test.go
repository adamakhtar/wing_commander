package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, FrameworkMinitest, config.TestFramework)
	assert.Equal(t, defaultMinitestCommand, config.TestCommand)
	assert.Equal(t, defaultResultsPath, config.TestResultsPath)
	assert.NotEmpty(t, config.ExcludePatterns)
	assert.Contains(t, config.ExcludePatterns, "/gems/")
	assert.Contains(t, config.ExcludePatterns, "/lib/ruby/")
}

func TestLoadConfig_MissingFileReturnsDefaults(t *testing.T) {
	config, err := LoadConfig("")
	assert.NoError(t, err)

	assert.Equal(t, FrameworkMinitest, config.TestFramework)
	assert.Equal(t, defaultMinitestCommand, config.TestCommand)
	assert.Equal(t, defaultResultsPath, config.TestResultsPath)
}

func TestLoadConfig_ValidFileOverridesDefaults(t *testing.T) {
	configDir := ".wing_commander"
	err := os.MkdirAll(configDir, 0o755)
	require.NoError(t, err)

	configPath := filepath.Join(configDir, "config.yml")
	configContent := `project_path: /tmp/project
test_framework: minitest
test_command: "bundle exec rake test test/workers/user_worker_test.rb"
test_results_path: "/tmp/project/.wing_commander/test_results/summary.yml"
exclude_patterns:
  - "/gems/"
  - "/custom/"`

	err = os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err)

	defer func() {
		os.Remove(configPath)
		os.Remove(configDir)
	}()

	config, err := LoadConfig("")
	require.NoError(t, err)

	assert.Equal(t, "/tmp/project", config.ProjectPath)
	assert.Equal(t, FrameworkMinitest, config.TestFramework)
	assert.Equal(t, "bundle exec rake test test/workers/user_worker_test.rb", config.TestCommand)
	assert.Equal(t, "/tmp/project/.wing_commander/test_results/summary.yml", config.TestResultsPath)
	assert.Equal(t, []string{"/gems/", "/custom/"}, config.ExcludePatterns)
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	configDir := ".wing_commander"
	err := os.MkdirAll(configDir, 0o755)
	require.NoError(t, err)

	configPath := filepath.Join(configDir, "config.yml")
	invalidYAML := `test_framework: minitest
invalid: yaml: content`

	err = os.WriteFile(configPath, []byte(invalidYAML), 0o644)
	require.NoError(t, err)

	defer func() {
		os.Remove(configPath)
		os.Remove(configDir)
	}()

	config, err := LoadConfig("")
	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestSaveConfig(t *testing.T) {
	cfg := &Config{
		ProjectPath:     "/tmp/project",
		TestFramework:   FrameworkMinitest,
		TestCommand:     "bundle exec rake test test/models/user_test.rb",
		TestResultsPath: "/tmp/project/.wing_commander/test_results/summary.yml",
		ExcludePatterns: []string{"/vendor/bundle/"},
	}

	err := SaveConfig(cfg)
	require.NoError(t, err)

	defer func() {
		os.Remove(".wing_commander/config.yml")
		os.Remove(".wing_commander")
	}()

	loaded, err := LoadConfig("")
	require.NoError(t, err)
	assert.Equal(t, cfg.ProjectPath, loaded.ProjectPath)
	assert.Equal(t, cfg.TestCommand, loaded.TestCommand)
	assert.Equal(t, cfg.TestResultsPath, loaded.TestResultsPath)
	assert.Equal(t, cfg.ExcludePatterns, loaded.ExcludePatterns)
}

func TestValidateFramework(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    TestFramework
		wantErr bool
	}{
		{
			name:  "empty string defaults to minitest",
			input: "",
			want:  FrameworkMinitest,
		},
		{
			name:  "minitest supported",
			input: "minitest",
			want:  FrameworkMinitest,
		},
		{
			name:    "other frameworks rejected",
			input:   "rspec",
			want:    FrameworkMinitest,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateFramework(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetDefaultTestCommand(t *testing.T) {
	assert.Equal(t, defaultMinitestCommand, GetDefaultTestCommand(FrameworkMinitest))
	assert.Equal(t, defaultMinitestCommand, GetDefaultTestCommand(TestFramework("rspec")))
}

func TestConfigWithMissingFieldsUsesDefaults(t *testing.T) {
	configDir := ".wing_commander"
	err := os.MkdirAll(configDir, 0o755)
	require.NoError(t, err)

	configPath := filepath.Join(configDir, "config.yml")
	configContent := `project_path: /tmp/project`

	err = os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err)

	defer func() {
		os.Remove(configPath)
		os.Remove(configDir)
	}()

	config, err := LoadConfig("")
	require.NoError(t, err)

	assert.Equal(t, FrameworkMinitest, config.TestFramework)
	assert.Equal(t, defaultMinitestCommand, config.TestCommand)
	assert.Equal(t, defaultResultsPath, config.TestResultsPath)
	assert.Equal(t, defaultExcludePatterns, config.ExcludePatterns)
}
