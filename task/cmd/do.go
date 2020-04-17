package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Marks a task as complete",
	Long:  "This subcommand marks an existing task as complete",
	Args:  cobra.MinimumNArgs(1),
	Example: `
  Complete a task:
  task do 1`,
	Run: do,
}

func do(cmd *cobra.Command, args []string) {
	fmt.Println("This is a fake \"do\" command")
	fmt.Println("Done:", args)
}

func init() {
	rootCmd.AddCommand(doCmd)
}
