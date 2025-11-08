/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adamakhtar/wing_commander/internal/config"
	"github.com/adamakhtar/wing_commander/internal/logger"
	"github.com/adamakhtar/wing_commander/internal/ui/styles"

	"github.com/adamakhtar/wing_commander/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start wing commander",
	Args:  cobra.MaximumNArgs(1),
	Long: "",
	Run: func(cmd *cobra.Command, args []string) {
		start(args)
	},
}

var (
	testFilePattern string
	runCommand string
	testResultsPath string
	debug bool
)

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable logging to debug problems.")
	startCmd.Flags().StringVarP(&testFilePattern, "test-file-pattern", "p", "", "A pattern to use to match test files in the project (e.g. 'test/**/*.rb')")

	startCmd.Flags().StringVarP(&runCommand, "run-command", "r", "", "The command to execute tests on the command line (e.g. 'rake test')")
	if err := startCmd.MarkFlagRequired("run-command"); err != nil { panic(err) }

	startCmd.Flags().StringVarP(&testResultsPath, "test-results-path", "t", "", "path to the directory where the test results are written (e.g. '.wing_commander/test_results')")
	if err := startCmd.MarkFlagRequired("test-results-path"); err != nil { panic(err) }
}

func start(args []string) {
	closeLogger := setupLogger()
	defer closeLogger()

	projectPath := processProjectPathArg(args)
	testResultsPath := processTestResultsPathOption(projectPath, testResultsPath)

	config := config.NewConfig(projectPath, runCommand, testFilePattern, testResultsPath, debug)
	styles := styles.BuildStyles(styles.DefaultTheme)
	model := ui.NewModel(config, styles)

	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		fmt.Printf("❌ Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

func processProjectPathArg(args []string) (string) {
	var projectPath string
	var err error

	if len(args) == 0 {
		projectPath, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current working directory: %v\n", err)
			os.Exit(1)
		}
	} else {
		projectPathArg := args[0]
		projectPath, err = filepath.Abs(projectPathArg)
		if err != nil {
			fmt.Printf("Error getting absolute path for your project %s: %v\n", projectPathArg, err)
			os.Exit(1)
		}
		info, err := os.Stat(projectPath)
		if os.IsNotExist(err) {
			fmt.Printf("❌ Error: The path to your project %s does not exist.\n", projectPathArg)
			os.Exit(1)
		}

		if !info.IsDir() {
			fmt.Printf("❌ Error: The path to your project %s is not a directory.\n", projectPath)
			os.Exit(1)
		}
	}

	return projectPath
}

func processTestResultsPathOption(projectPath string, providedPath string) string {
    var absPath string
    if filepath.IsAbs(providedPath) {
        absPath = providedPath
    } else {
        absPath = filepath.Join(projectPath, providedPath)
    }

    info, err := os.Stat(absPath)
    if err != nil || info.IsDir() {
        fmt.Printf("❌ Error: testResultsPath %s must be an existing file\n", absPath)
        os.Exit(1)
    }

    return absPath
}

func setupLogger() (func() error) {
	closeLogger, err := logger.SetupLogger(debug)
	if err != nil {
		fmt.Printf("Error setting up logger: %v\n", err)
		os.Exit(1)
	}
	return closeLogger
}