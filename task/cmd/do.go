package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/prmsrswt/gophercises/task/store"
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
	s, err := store.NewStore()
	if err != nil {
		fmt.Println("Error creating store:", err)
		os.Exit(1)
	}
	defer s.Close()

	i, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Error: invalid ID", err)
		os.Exit(1)
	}

	tasks, err := s.GetActiveTasks()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if i > len(tasks) || i <= 0 {
		fmt.Println("Error: Invalid ID")
		os.Exit(1)
	}

	id := tasks[i-1].ID

	data, err := s.CompleteTask(id)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Printf("You have completed the \"%s\" task.\n", data.Value)
}

func init() {
	rootCmd.AddCommand(doCmd)
}
