package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "Task is a command line task manager",
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
