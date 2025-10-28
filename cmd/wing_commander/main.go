package main

import (
	"fmt"

	"github.com/adamakhtar/wing_commander/internal/config"
	"github.com/adamakhtar/wing_commander/internal/logger"
	"github.com/adamakhtar/wing_commander/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		fmt.Printf("❌ Error loading config: %v\n", err)
		return
	}

	// Logging
	closeLogger, err := logger.SetupLogger(cfg.Debug)
	if err != nil {
		fmt.Printf("❌ Error setting up logger: %v\n", err)
		return
	}
	defer closeLogger()

	model := ui.NewModel(cfg)
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		fmt.Printf("❌ Error running TUI: %v\n", err)
		return
	}
	// fmt.Println("🚀 Wing Commander - Test Failure Analyzer")
	// fmt.Println("==========================================")
	// fmt.Println()

	// // Load configuration
	// cfg, err := config.LoadConfig("")
	// if err != nil {
	// 	fmt.Printf("❌ Error loading config: %v\n", err)
	// 	return
	// }

	// // Handle basic commands
	// if len(os.Args) > 1 {
	// 	switch os.Args[1] {
	// 	case "version":
	// 		fmt.Println("Wing Commander v0.1.0")
	// 		fmt.Println("Built with Go")
	// 		return
	// 	case "run":
	// 		runCommand(os.Args[2:])
	// 		return
	// 	case "config":
	// 		// Load config with CLI options to show the actual project path
	// 		configWithCLI, err := loadConfigWithCLIOptions("", "", "")
	// 		if err != nil {
	// 			fmt.Printf("❌ Error loading config: %v\n", err)
	// 			return
	// 		}
	// 		showConfig(configWithCLI)
	// 		return
	// 	default:
	// 		// Check if it's an XML file path
	// 		if isXMLFile(os.Args[1]) {
	// 			parseAndDisplayXML(os.Args[1], cfg)
	// 			return
	// 		}

	// 		fmt.Printf("Unknown command: %s\n", os.Args[1])
	// 		fmt.Println("Available commands: run, version, config")
	// 		fmt.Println("Or provide an XML file path to parse test results")
	// 		return
	// 	}
	// }

	// Default welcome message
	// fmt.Println("Welcome to Wing Commander!")
	// fmt.Println()
	// fmt.Println("Usage:")
	// fmt.Println("  wing_commander [command]")
	// fmt.Println("  wing_commander <xml-file>")
	// fmt.Println()
	// fmt.Println("Commands:")
	// fmt.Println("  run         - Run tests and analyze failures")
	// fmt.Println("             --config PATH         Specify config file path (default: .wing_commander/config.yml)")
	// fmt.Println("             --project-path PATH   Path to project directory (default: current working directory)")
	// fmt.Println("             --test-command CMD    Test runner command with interpolation (e.g., 'rails test {{.Paths}} --output junit')")
	// fmt.Println("  version     - Show version information")
	// fmt.Println("  config      - Show current configuration")
	// fmt.Println()
	// fmt.Println("XML Files:")
	// fmt.Println("  Parse test results from JUnit XML files")
	// fmt.Println("  Example: wing_commander testdata/fixtures/rspec_failures.xml")
	// fmt.Println()
	// fmt.Printf("📋 Current framework: %s\n", cfg.TestFramework)
	// fmt.Printf("📋 Test command: %s\n", cfg.TestCommand)
	// fmt.Println()

	// // Demonstrate that our types work
	// frame := types.NewStackFrame("app/models/user.rb", 42, "create_user")
	// fmt.Printf("Sample StackFrame: %+v\n", frame)

	// test := types.NewTestResult("User creation test", types.StatusFail)
	// fmt.Printf("Sample TestResult: %+v\n", test)

	// group := types.NewFailureGroup("abc123", "Validation failed")
	// fmt.Printf("Sample FailureGroup: %+v\n", group)

	// fmt.Println()
	// fmt.Println("✅ Core types loaded successfully!")
	// fmt.Println("📋 Next step: Implement JUnit XML parser")
}

// isXMLFile checks if the argument looks like an XML file
// func isXMLFile(arg string) bool {
// 	return len(arg) > 4 && arg[len(arg)-4:] == ".xml"
// }

// // parseAndDisplayXML parses an XML file and displays the results
// func parseAndDisplayXML(filePath string, cfg *config.Config) {
// 	fmt.Printf("📄 Parsing XML file: %s\n", filePath)
// 	fmt.Println()

// 	result, err := parser.ParseFile(filePath)
// 	if err != nil {
// 		fmt.Printf("❌ Error parsing file: %v\n", err)
// 		return
// 	}

// 	fmt.Printf("✅ Successfully parsed %d test results\n", len(result.Tests))
// 	fmt.Printf("🔍 Framework: %s\n", cfg.TestFramework)
// 	fmt.Println()

// 	// Normalize backtraces using config exclude patterns
// 	normalizer := grouper.NewNormalizer(cfg)
// 	normalizedResults := normalizer.NormalizeTestResults(result.Tests)

// 	// Show filtering statistics
// 	totalFrames, filteredFrames := grouper.CountFilteredFrames(normalizedResults)
// 	fmt.Printf("🔧 Backtrace Filtering:\n")
// 	fmt.Printf("  Total frames:   %d\n", totalFrames)
// 	fmt.Printf("  Project frames: %d\n", filteredFrames)
// 	fmt.Printf("  Filtered out:   %d\n", totalFrames-filteredFrames)
// 	fmt.Println()

// 	// Show summary
// 	fmt.Println("📊 Test Summary:")
// 	fmt.Printf("  Total:   %d\n", result.Summary.Total)
// 	fmt.Printf("  Passed:  %d\n", result.Summary.Passed)
// 	fmt.Printf("  Failed:  %d\n", result.Summary.Failed)
// 	fmt.Printf("  Skipped: %d\n", result.Summary.Skipped)
// 	fmt.Println()

// 	// Show failed tests
// 	failedTests := 0
// 	for _, test := range normalizedResults {
// 		if test.Status == types.StatusFail {
// 			failedTests++
// 			fmt.Printf("❌ %s\n", test.Name)
// 			if test.ErrorMessage != "" {
// 				fmt.Printf("   Error: %s\n", test.ErrorMessage)
// 			}
// 			if len(test.FullBacktrace) > 0 {
// 				fmt.Printf("   Full backtrace: %d frames\n", len(test.FullBacktrace))
// 			}
// 			if len(test.FilteredBacktrace) > 0 {
// 				fmt.Printf("   Project frames: %d frames\n", len(test.FilteredBacktrace))
// 			}
// 			fmt.Println()
// 		}
// 	}

// 	if failedTests == 0 {
// 		fmt.Println("🎉 No failed tests found!")
// 	} else {
// 		fmt.Printf("🔍 Found %d failed tests\n", failedTests)
// 		fmt.Println("📋 Next step: Implement failure grouping")
// 	}
// }

// // showConfig displays the current configuration
// func showConfig(cfg *config.Config) {
// 	fmt.Println("📋 Wing Commander Configuration")
// 	fmt.Println("===============================")
// 	fmt.Println()

// 	fmt.Printf("Project Path:   %s\n", cfg.ProjectPath)
// 	fmt.Printf("Test Framework: %s\n", cfg.TestFramework)
// 	fmt.Printf("Test Command:   %s\n", cfg.TestCommand)
// 	fmt.Println()

// 	fmt.Println("Exclude Patterns:")
// 	for _, pattern := range cfg.ExcludePatterns {
// 		fmt.Printf("  - %s\n", pattern)
// 	}
// 	fmt.Println()

// 	fmt.Println("Configuration file: .wing_commander/config.yml")
// 	fmt.Println("Create this file to customize settings.")
// }

// // runCommand handles the "run" command with flag parsing
// func runCommand(args []string) {
// 	// Create flag set for run command
// 	runFlags := flag.NewFlagSet("run", flag.ExitOnError)
// 	configPath := runFlags.String("config", ".wing_commander/config.yml", "Path to config file")
// 	projectPath := runFlags.String("project-path", "", "Path to the project whose tests are being observed (default: current working directory)")
// 	testCommand := runFlags.String("test-command", "", "Test runner command with interpolation support (e.g., 'rails test {{.Paths}} --output junit')")

// 	// Parse flags
// 	if err := runFlags.Parse(args); err != nil {
// 		fmt.Printf("❌ Error parsing flags: %v\n", err)
// 		return
// 	}

// 	// Load configuration with CLI options taking precedence
// 	cfg, err := loadConfigWithCLIOptions(*configPath, *projectPath, *testCommand)
// 	if err != nil {
// 		fmt.Printf("❌ Error loading config: %v\n", err)
// 		return
// 	}

// 	// Execute tests
// 	runTests(cfg)
// }

// // loadConfigWithCLIOptions loads configuration with CLI options taking precedence over config file
// func loadConfigWithCLIOptions(configPath, projectPath, testCommand string) (*config.Config, error) {
// 	// Load base configuration from file
// 	cfg, err := config.LoadConfig(configPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to load config: %w", err)
// 	}

// 	// Override with CLI options if provided
// 	if projectPath != "" {
// 		// Convert to absolute path if relative
// 		if !filepath.IsAbs(projectPath) {
// 			absPath, err := filepath.Abs(projectPath)
// 			if err != nil {
// 				return nil, fmt.Errorf("failed to resolve project path: %w", err)
// 			}
// 			projectPath = absPath
// 		}
// 		cfg.ProjectPath = projectPath
// 	} else {
// 		// Default to current working directory
// 		cwd, err := os.Getwd()
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to get current working directory: %w", err)
// 		}
// 		cfg.ProjectPath = cwd
// 	}

// 	if testCommand != "" {
// 		cfg.TestCommand = testCommand
// 	}

// 	return cfg, nil
// }

// // runTests executes tests using the TestRunner service and launches TUI
// func runTests(cfg *config.Config) {
// 	fmt.Println("🚀 Wing Commander - Test Failure Analyzer")
// 	fmt.Println("==========================================")
// 	fmt.Println()

// 	// Create test runner
// 	testRunner := runner.NewTestRunner(cfg)

// 	// Validate configuration
// 	if err := testRunner.ValidateConfig(); err != nil {
// 		fmt.Printf("❌ Configuration error: %v\n", err)
// 		return
// 	}

// 	// Check if test command exists
// 	if err := testRunner.CheckTestCommandExists(); err != nil {
// 		fmt.Printf("❌ Test command error: %v\n", err)
// 		return
// 	}

// 	// Display configuration
// 	fmt.Printf("📋 Framework: %s\n", cfg.TestFramework)
// 	fmt.Printf("📋 Command: %s\n", cfg.TestCommand)
// 	fmt.Println()

// 	// Launch TUI in empty state (no tests executed yet)
// 	fmt.Println("🚀 Launching interactive UI...")
// 	fmt.Println()

// 	// Create UI model with nil result (empty state)
// 	model := ui.NewModel(nil, testRunner)

// 	// Create and run the TUI program
// 	program := tea.NewProgram(model, tea.WithAltScreen())
// 	if _, err := program.Run(); err != nil {
// 		fmt.Printf("❌ Error running TUI: %v\n", err)
// 		return
// 	}
// }
