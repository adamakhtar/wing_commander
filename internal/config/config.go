package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// TestFramework represents the supported test framework.
type TestFramework string

const (
	FrameworkMinitest TestFramework = "minitest"
)

const (
	defaultConfigDir       = ".wing_commander"
	defaultConfigFile      = "config.yml"
	defaultResultsPath     = ".wing_commander/test_results/summary.yml"
	defaultMinitestCommand = "bundle exec rake test {{.Paths}}"
)

var defaultExcludePatterns = []string{
	"/gems/",
	"/lib/ruby/",
	"/vendor/bundle/",
}

// Config represents the Wing Commander configuration.
type Config struct {
	ProjectPath        string        `yaml:"project_path"`
	TestFramework      TestFramework `yaml:"test_framework"`
	TestCommand        string        `yaml:"test_command"`
	RunTestCaseCommand string        `yaml:"run_test_case_command"`
	TestFilePattern    string        `yaml:"test_file_pattern"`
	TestResultsPath    string        `yaml:"test_results_path"`
	Debug              bool          `yaml:"debug"`
	ExcludePatterns    []string      `yaml:"exclude_patterns"`
}

// NewConfig creates a new configuration instance, applying sensible defaults for
// WingCommanderReporter-based YAML summaries when values are omitted.
func NewConfig(projectPath string, testCommand string, testFilePattern string, testResultsPath string, runTestCaseCommand string, debug bool) *Config {
	cfg := DefaultConfig()

	if projectPath != "" {
		cfg.ProjectPath = projectPath
	}
	if testCommand != "" {
		cfg.TestCommand = testCommand
	}
	if runTestCaseCommand != "" {
		cfg.RunTestCaseCommand = runTestCaseCommand
	}
	if testFilePattern != "" {
		cfg.TestFilePattern = testFilePattern
	}
	if testResultsPath != "" {
		cfg.TestResultsPath = testResultsPath
	}
	cfg.ensureRunTestCaseCommand()
	cfg.Debug = debug

	return cfg
}

// DefaultConfig returns the baseline configuration for the Wing Commander CLI.
func DefaultConfig() *Config {
	cfg := &Config{
		ProjectPath:        "",
		TestFramework:      FrameworkMinitest,
		TestCommand:        defaultMinitestCommand,
		RunTestCaseCommand: "",
		TestFilePattern:    "",
		TestResultsPath:    defaultResultsPath,
		Debug:              false,
		ExcludePatterns:    append([]string{}, defaultExcludePatterns...),
	}

	cfg.ensureRunTestCaseCommand()
	return cfg
}

// LoadConfig reads configuration from disk, falling back to defaults if the file
// does not exist. Missing fields are backfilled with defaults.
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	configPath := path
	if configPath == "" {
		configPath = filepath.Join(defaultConfigDir, defaultConfigFile)
	}

	if !filepath.IsAbs(configPath) {
		abs, err := filepath.Abs(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve config path %s: %w", configPath, err)
		}
		configPath = abs
	}

	data, err := os.ReadFile(configPath)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var loaded Config
	if err := yaml.Unmarshal(data, &loaded); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	if loaded.ProjectPath != "" {
		cfg.ProjectPath = loaded.ProjectPath
	}
	if loaded.TestFramework != "" {
		cfg.TestFramework = loaded.TestFramework
	}
	if loaded.TestCommand != "" {
		cfg.TestCommand = loaded.TestCommand
	}
	if loaded.RunTestCaseCommand != "" {
		cfg.RunTestCaseCommand = loaded.RunTestCaseCommand
	}
	if loaded.TestFilePattern != "" {
		cfg.TestFilePattern = loaded.TestFilePattern
	}
	if loaded.TestResultsPath != "" {
		cfg.TestResultsPath = loaded.TestResultsPath
	}
	if len(loaded.ExcludePatterns) > 0 {
		cfg.ExcludePatterns = loaded.ExcludePatterns
	}
	if loaded.Debug {
		cfg.Debug = true
	}

	cfg.ensureRunTestCaseCommand()
	return cfg, nil
}

// SaveConfig writes the current configuration to the default location on disk.
func SaveConfig(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("cannot save nil config")
	}

	cfg.ensureRunTestCaseCommand()

	if err := os.MkdirAll(defaultConfigDir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", defaultConfigDir, err)
	}

	configPath := filepath.Join(defaultConfigDir, defaultConfigFile)

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", configPath, err)
	}

	return nil
}

// ValidateFramework ensures that the provided framework is supported.
func ValidateFramework(value string) (TestFramework, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", string(FrameworkMinitest):
		return FrameworkMinitest, nil
	default:
		return FrameworkMinitest, fmt.Errorf("unsupported test framework: %s (Wing Commander currently requires WingCommanderReporter YAML output)", value)
	}
}

// GetDefaultTestCommand returns the WingCommanderReporter-aware test command.
func GetDefaultTestCommand(framework TestFramework) string {
	return defaultMinitestCommand
}

func (cfg *Config) ensureRunTestCaseCommand() {
	if cfg.RunTestCaseCommand == "" {
		cfg.RunTestCaseCommand = cfg.TestCommand
	}
}
