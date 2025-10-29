/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start wing commander",
	Args:  cobra.MaximumNArgs(1),
	Long: "",
	Run: func(cmd *cobra.Command, args []string) {
		start(cmd, args)
	},
}

var (
	testsMatching string
	runCommand    string
	debug        bool
)

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable logging to debug problems.")
	startCmd.Flags().StringVarP(&testsMatching, "search-pattern", "s", "", "Only files matching this pattern will be searchable to run (e.g. 'test/**/*.rb')")

	startCmd.Flags().StringVarP(&runCommand, "command", "c", "", "Command to run tests (e.g. 'rake test')")
	if err := startCmd.MarkFlagRequired("command"); err != nil { panic(err) }
}

func start(cmd *cobra.Command, args []string) {
	projectPath := processProjectPathArg(args)
	fmt.Printf("Starting Wing Commander for project at %s\n", projectPath)
	fmt.Printf("Debug mode is %t\n", debug)
	fmt.Printf("Test command is %s\n", runCommand)
	fmt.Printf("Test matching is %s\n", testsMatching)
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