package cmd

import (
	"fmt"
	"os"

	"github.com/prmsrswt/gophercises/task/store"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List current tasks",
	Long:    "This subcommand lists all tasks which are currenlty active",
	Aliases: []string{"ls"},
	Run:     list,
}

func list(cmd *cobra.Command, args []string) {
	s, err := store.NewStore()
	if err != nil {
		fmt.Println("Error creating store:", err)
		os.Exit(1)
	}
	defer s.Close()

	tasks, err := s.GetActiveTasks()
	if err != nil {
		fmt.Println("Error fetching tasks from db:", err)
		os.Exit(1)
	}

	if len(tasks) == 0 {
		fmt.Println("You have no tasks to complete! Why not take a vacation? 🏖")
		return
	}

	fmt.Println("You have the following tasks:")

	for i, v := range tasks {
		fmt.Printf("%d. %s\n", i+1, v.Value)
	}

}

func init() {
	rootCmd.AddCommand(listCmd)
}
