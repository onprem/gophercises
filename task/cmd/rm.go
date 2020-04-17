package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/prmsrswt/gophercises/task/store"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:     "rm",
	Short:   "Permanently removes a task",
	Long:    "This subcommand removes an existing task permanently",
	Aliases: []string{"remove", "delete"},
	Example: `  Delete an active task:
  task rm 1

  Delete a completed tast
  task rm 1 --done

  Delete all tasks:
  task rm -a`,
	Run: rm,
}
var deleteAll, isCompleted bool

func rm(cmd *cobra.Command, args []string) {
	s, err := store.NewStore()
	if err != nil {
		fmt.Println("Error creating store:", err)
		os.Exit(1)
	}
	defer s.Close()

	if deleteAll {
		err = s.DeleteAllTasks()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		return
	}

	if len(args) == 0 {
		fmt.Println("Error: no ID provided to delete")
	}
	i, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Error: invalid ID", err)
		os.Exit(1)
	}

	tasks, err := s.GetActiveTasks()
	if isCompleted {
		tasks, err = s.GetAllCompletedTasks()
	}

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if i > len(tasks) || i <= 0 {
		fmt.Println("Error: Invalid ID")
		os.Exit(1)
	}

	id := tasks[i-1].ID

	data, err := s.DeleteTask(id)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Printf("You have removed the \"%s\" task.\n", data.Value)
}

func init() {
	rmCmd.Flags().BoolVarP(&deleteAll, "all", "a", false, "Delete all tasks")
	rmCmd.Flags().BoolVarP(&isCompleted, "done", "d", false, "Delete a completed task")
	rootCmd.AddCommand(rmCmd)
}
