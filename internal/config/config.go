package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// TestFramework represents different test frameworks
type TestFramework string

const (
	FrameworkRSpec    TestFramework = "rspec"
	FrameworkMinitest TestFramework = "minitest"
	FrameworkPytest   TestFramework = "pytest"
	FrameworkJest     TestFramework = "jest"
	FrameworkUnknown  TestFramework = "unknown"
)

// Config represents the Wing Commander configuration
type Config struct {
	TestFramework   TestFramework `yaml:"test_framework"`
	TestCommand     string        `yaml:"test_command"`
	ExcludePatterns []string      `yaml:"exclude_patterns"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		TestFramework: FrameworkUnknown,
		TestCommand:   "bundle exec rspec --format RspecJunitFormatter --out results.xml",
		ExcludePatterns: []string{
			"/gems/",
			"/lib/ruby/",
			"/vendor/bundle/",
			"/.rbenv/",
			"/.rvm/",
			"/rspec-core/",
			"/minitest/",
			"/node_modules/",
			"/.venv/",
			"/site-packages/",
		},
	}
}

// LoadConfig loads configuration from the specified path or .wing_commander/config.yml if empty
func LoadConfig(configPath string) (*Config, error) {
	// Use default path if none provided
	if configPath == "" {
		configPath = ".wing_commander/config.yml"
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config with warning
		fmt.Printf("⚠️  Config file not found: %s\n", configPath)
		fmt.Println("Using default configuration. Create config file to customize settings.")
		return DefaultConfig(), nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate and set defaults for missing fields
	if config.TestFramework == "" {
		config.TestFramework = FrameworkUnknown
	}
	if config.TestCommand == "" {
		config.TestCommand = "bundle exec rspec --format RspecJunitFormatter --out results.xml"
	}
	if len(config.ExcludePatterns) == 0 {
		config.ExcludePatterns = DefaultConfig().ExcludePatterns
	}

	return &config, nil
}

// SaveConfig saves configuration to .wing_commander/config.yml
func SaveConfig(config *Config) error {
	// Create .wing_commander directory if it doesn't exist
	configDir := ".wing_commander"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.yml")

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ValidateFramework checks if the framework string is valid
func ValidateFramework(framework string) (TestFramework, error) {
	switch framework {
	case "rspec":
		return FrameworkRSpec, nil
	case "minitest":
		return FrameworkMinitest, nil
	case "pytest":
		return FrameworkPytest, nil
	case "jest":
		return FrameworkJest, nil
	case "unknown":
		return FrameworkUnknown, nil
	default:
		return FrameworkUnknown, fmt.Errorf("unknown test framework: %s", framework)
	}
}

// GetDefaultTestCommand returns the default test command for a framework
func GetDefaultTestCommand(framework TestFramework) string {
	switch framework {
	case FrameworkRSpec:
		return "bundle exec rspec --format RspecJunitFormatter --out results.xml"
	case FrameworkMinitest:
		return "bundle exec rake test TESTOPTS='--junit --junit-filename=results.xml'"
	case FrameworkPytest:
		return "pytest --junit-xml=results.xml"
	case FrameworkJest:
		return "npx jest --reporters=jest-junit"
	default:
		return "bundle exec rspec --format RspecJunitFormatter --out results.xml"
	}
}
