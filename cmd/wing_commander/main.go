package main

import (
	"fmt"
	"os"

	"github.com/adamakhtar/wing_commander/internal/types"
)

func main() {
	fmt.Println("🚀 Wing Commander - Test Failure Analyzer")
	fmt.Println("==========================================")
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
	fmt.Println()
	fmt.Println("Usage: wing_commander [command]")
	fmt.Println("Commands:")
	fmt.Println("  (no args)  - Show this welcome message")
	fmt.Println("  run         - Run tests and analyze failures (coming soon)")
	fmt.Println("  version     - Show version information")

	// Handle basic commands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			fmt.Println()
			fmt.Println("Wing Commander v0.1.0")
			fmt.Println("Built with Go")
		case "run":
			fmt.Println()
			fmt.Println("🚧 Test runner not implemented yet!")
			fmt.Println("This will be added in Step 8 of the implementation plan.")
		default:
			fmt.Println()
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			fmt.Println("Available commands: run, version")
		}
	}
}
