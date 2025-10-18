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

	assert.Equal(t, FrameworkUnknown, config.TestFramework)
	assert.Equal(t, "bundle exec rspec --format json", config.TestCommand)
	assert.NotEmpty(t, config.ExcludePatterns)
	assert.Contains(t, config.ExcludePatterns, "/gems/")
	assert.Contains(t, config.ExcludePatterns, "/lib/ruby/")
}

func TestLoadConfig_MissingFile(t *testing.T) {
	// Test with non-existent config file
	config, err := LoadConfig()

	// Should return default config without error
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, FrameworkUnknown, config.TestFramework)
}

func TestLoadConfig_ValidFile(t *testing.T) {
	// Create a temporary config file
	configDir := ".wing_commander"
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	configPath := filepath.Join(configDir, "config.yml")
	configContent := `test_framework: rspec
test_command: "bundle exec rspec --format json"
exclude_patterns:
  - "/gems/"
  - "/lib/ruby/"
  - "/custom/path/"`

	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Clean up after test
	defer func() {
		os.Remove(configPath)
		os.Remove(configDir)
	}()

	config, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, FrameworkRSpec, config.TestFramework)
	assert.Equal(t, "bundle exec rspec --format json", config.TestCommand)
	assert.Contains(t, config.ExcludePatterns, "/gems/")
	assert.Contains(t, config.ExcludePatterns, "/custom/path/")
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	// Create a temporary config file with invalid YAML
	configDir := ".wing_commander"
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	configPath := filepath.Join(configDir, "config.yml")
	invalidYAML := `test_framework: rspec
test_command: "bundle exec rspec --format json"
exclude_patterns:
  - "/gems/"
  - "/lib/ruby/"
invalid: yaml: content`

	err = os.WriteFile(configPath, []byte(invalidYAML), 0644)
	require.NoError(t, err)

	// Clean up after test
	defer func() {
		os.Remove(configPath)
		os.Remove(configDir)
	}()

	config, err := LoadConfig()
	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestSaveConfig(t *testing.T) {
	config := &Config{
		TestFramework: FrameworkRSpec,
		TestCommand:   "bundle exec rspec --format json",
		ExcludePatterns: []string{
			"/gems/",
			"/custom/path/",
		},
	}

	err := SaveConfig(config)
	require.NoError(t, err)

	// Clean up after test
	defer func() {
		os.Remove(".wing_commander/config.yml")
		os.Remove(".wing_commander")
	}()

	// Verify file was created
	_, err = os.Stat(".wing_commander/config.yml")
	assert.NoError(t, err)

	// Load and verify content
	loadedConfig, err := LoadConfig()
	require.NoError(t, err)
	assert.Equal(t, config.TestFramework, loadedConfig.TestFramework)
	assert.Equal(t, config.TestCommand, loadedConfig.TestCommand)
	assert.Equal(t, config.ExcludePatterns, loadedConfig.ExcludePatterns)
}

func TestValidateFramework(t *testing.T) {
	tests := []struct {
		input    string
		expected TestFramework
		wantErr  bool
	}{
		{"rspec", FrameworkRSpec, false},
		{"minitest", FrameworkMinitest, false},
		{"pytest", FrameworkPytest, false},
		{"jest", FrameworkJest, false},
		{"unknown", FrameworkUnknown, false},
		{"invalid", FrameworkUnknown, true},
		{"", FrameworkUnknown, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ValidateFramework(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetDefaultTestCommand(t *testing.T) {
	tests := []struct {
		framework TestFramework
		expected  string
	}{
		{FrameworkRSpec, "bundle exec rspec --format json"},
		{FrameworkMinitest, "bundle exec ruby -I test test/*_test.rb --format json"},
		{FrameworkPytest, "pytest --json-report --json-report-file=test-results.json"},
		{FrameworkJest, "npx jest --json --outputFile=test-results.json"},
		{FrameworkUnknown, "bundle exec rspec --format json"},
	}

	for _, tt := range tests {
		t.Run(string(tt.framework), func(t *testing.T) {
			result := GetDefaultTestCommand(tt.framework)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfigWithMissingFields(t *testing.T) {
	// Create a config file with missing fields
	configDir := ".wing_commander"
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	configPath := filepath.Join(configDir, "config.yml")
	configContent := `test_framework: rspec
# test_command missing
exclude_patterns:
  - "/gems/"`

	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Clean up after test
	defer func() {
		os.Remove(configPath)
		os.Remove(configDir)
	}()

	config, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, config)

	// Should use defaults for missing fields
	assert.Equal(t, FrameworkRSpec, config.TestFramework)
	assert.Equal(t, "bundle exec rspec --format json", config.TestCommand) // Default
	assert.Contains(t, config.ExcludePatterns, "/gems/")
}
