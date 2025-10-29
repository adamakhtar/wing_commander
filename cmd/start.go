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

	"github.com/adamakhtar/wing_commander/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
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
	runCommand    string
	debug        bool
)

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable logging to debug problems.")
	startCmd.Flags().StringVarP(&testFilePattern, "test-file-pattern", "t", "", "Only files matching this pattern will be appear as test files to run (e.g. 'test/**/*.rb')")

	startCmd.Flags().StringVarP(&runCommand, "run-command", "c", "", "Command to execute tests (e.g. 'rake test')")
	if err := startCmd.MarkFlagRequired("run-command"); err != nil { panic(err) }
}

func start(args []string) {
	closeLogger := setupLogger()
	defer closeLogger()

	projectPath := processProjectPathArg(args)

	log.Debug("Starting Wing Commander", "projectPath", projectPath, "-debug", debug, "-command", runCommand, "-test-file-pattern", testFilePattern)

	config := config.NewConfig(projectPath, runCommand, testFilePattern, debug)
	model := ui.NewModel(config)

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

func setupLogger() (func() error) {
	closeLogger, err := logger.SetupLogger(debug)
	if err != nil {
		fmt.Printf("Error setting up logger: %v\n", err)
		os.Exit(1)
	}
	return closeLogger
}