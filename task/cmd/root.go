package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "task",
	Short:   "Task is a command line task manager",
	Version: "v0.1.0",
	Run:     list,
}
var path = "tasks.db"

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
