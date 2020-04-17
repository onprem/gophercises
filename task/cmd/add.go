package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/prmsrswt/gophercises/task/store"
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
	data := strings.Join(args, " ")

	s, err := store.NewStore()
	if err != nil {
		fmt.Println("Error creating store:", err)
		os.Exit(1)
	}
	defer s.Close()

	err = s.InsertTask(data)
	if err != nil {
		fmt.Println("Error inserting task to db:", err)
		os.Exit(1)
	}

	fmt.Printf("Added \"%s\" to your task list.\n", data)
}

func init() {
	rootCmd.AddCommand(addCmd)
}
