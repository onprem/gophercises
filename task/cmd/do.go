package cmd

import (
	"fmt"
	"os"
	"strconv"

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
	s, err := newStore()
	if err != nil {
		fmt.Println("Error creating store:", err)
		os.Exit(1)
	}
	defer s.db.Close()

	i, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Error: invalid ID", err)
		os.Exit(1)
	}

	tasks, err := s.getActiveTasks()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	id := tasks[i-1].ID

	data, err := s.completeTask(id)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Printf("You have completed the \"%s\" task.\n", data.Task)
}

func init() {
	rootCmd.AddCommand(doCmd)
}
