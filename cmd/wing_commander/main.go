package main

import (
	"fmt"
	"os"

	"github.com/adamakhtar/wing_commander/internal/parser"
	"github.com/adamakhtar/wing_commander/internal/types"
)

func main() {
	fmt.Println("🚀 Wing Commander - Test Failure Analyzer")
	fmt.Println("==========================================")
	fmt.Println()
	
	// Handle basic commands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			fmt.Println("Wing Commander v0.1.0")
			fmt.Println("Built with Go")
			return
		case "run":
			fmt.Println()
			fmt.Println("🚧 Test runner not implemented yet!")
			fmt.Println("This will be added in Step 8 of the implementation plan.")
			return
		default:
			// Check if it's a JSON file path
			if isJSONFile(os.Args[1]) {
				parseAndDisplayJSON(os.Args[1])
				return
			}
			
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			fmt.Println("Available commands: run, version")
			fmt.Println("Or provide a JSON file path to parse test results")
			return
		}
	}
	
	// Default welcome message
	fmt.Println("Welcome to Wing Commander!")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  wing_commander [command]")
	fmt.Println("  wing_commander <json-file>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  run         - Run tests and analyze failures (coming soon)")
	fmt.Println("  version     - Show version information")
	fmt.Println()
	fmt.Println("JSON Files:")
	fmt.Println("  Parse test results from JSON files")
	fmt.Println("  Example: wing_commander testdata/fixtures/rspec_failures.json")
	fmt.Println()
	
	// Demonstrate that our types work
	frame := types.NewStackFrame("app/models/user.rb", 42, "create_user")
	fmt.Printf("Sample StackFrame: %+v\n", frame)
	
	test := types.NewTestResult("User creation test", types.StatusFail)
	fmt.Printf("Sample TestResult: %+v\n", test)
	
	group := types.NewFailureGroup("abc123", "Validation failed")
	fmt.Printf("Sample FailureGroup: %+v\n", group)
	
	fmt.Println()
	fmt.Println("✅ Core types loaded successfully!")
	fmt.Println("📋 Next step: Implement JSON parser")
}

// isJSONFile checks if the argument looks like a JSON file
func isJSONFile(arg string) bool {
	return len(arg) > 5 && arg[len(arg)-5:] == ".json"
}

// parseAndDisplayJSON parses a JSON file and displays the results
func parseAndDisplayJSON(filePath string) {
	fmt.Printf("📄 Parsing JSON file: %s\n", filePath)
	fmt.Println()
	
	result, err := parser.ParseFile(filePath)
	if err != nil {
		fmt.Printf("❌ Error parsing file: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Successfully parsed %d test results\n", len(result.Tests))
	fmt.Printf("🔍 Detected framework: %s\n", result.Framework)
	fmt.Println()
	
	// Show summary
	fmt.Println("📊 Test Summary:")
	fmt.Printf("  Total:   %d\n", result.Summary.Total)
	fmt.Printf("  Passed:  %d\n", result.Summary.Passed)
	fmt.Printf("  Failed:  %d\n", result.Summary.Failed)
	fmt.Printf("  Skipped: %d\n", result.Summary.Skipped)
	fmt.Println()
	
	// Show failed tests
	failedTests := 0
	for _, test := range result.Tests {
		if test.Status == types.StatusFail {
			failedTests++
			fmt.Printf("❌ %s\n", test.Name)
			if test.ErrorMessage != "" {
				fmt.Printf("   Error: %s\n", test.ErrorMessage)
			}
			if len(test.FullBacktrace) > 0 {
				fmt.Printf("   Backtrace: %d frames\n", len(test.FullBacktrace))
			}
			fmt.Println()
		}
	}
	
	if failedTests == 0 {
		fmt.Println("🎉 No failed tests found!")
	} else {
		fmt.Printf("🔍 Found %d failed tests\n", failedTests)
		fmt.Println("📋 Next step: Implement backtrace grouping")
	}
}
