package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task",
	Long:  "This subcommand adds a new task to current tasks",
	Args:  cobra.MinimumNArgs(1),
	Example: `
  Add a task:
  task add learn about bbolt`,
	Run: add,
}

func add(cmd *cobra.Command, args []string) {
	fmt.Println("This is a fake \"add\" command")
	task := strings.Join(args, " ")
	fmt.Println("Adding:", task)
}

func init() {
	rootCmd.AddCommand(addCmd)
}
