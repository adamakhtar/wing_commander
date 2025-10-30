package config

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
	ProjectPath     string
	TestFramework   TestFramework
	TestCommand     string
	TestFilePattern string
	TestResultsPath string
	Debug           bool
	ExcludePatterns []string // TODO deprecated - remove
}

func NewConfig(projectPath string, testCommand string, testFilePattern string, testResultsPath string, debug bool) *Config {
	return &Config{
		ProjectPath: projectPath,
		TestCommand: testCommand,
		TestFilePattern: testFilePattern,
		TestResultsPath: testResultsPath,
		Debug: debug,
	}
}