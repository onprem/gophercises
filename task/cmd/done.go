package cmd

import (
	"fmt"
	"os"

	"github.com/prmsrswt/gophercises/task/store"
	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:     "done",
	Short:   "List completed tasks",
	Long:    "This subcommand lists tasks which are completed",
	Aliases: []string{"completed"},
	Run:     done,
}
var showAll bool

func done(cmd *cobra.Command, args []string) {
	s, err := store.NewStore()
	if err != nil {
		fmt.Println("Error creating store:", err)
		os.Exit(1)
	}
	defer s.Close()

	var tasks []store.Task
	var noTaskMsg, taskMsg string

	if showAll {
		tasks, err = s.GetAllCompletedTasks()
		taskMsg = "You have completed following tasks:"
		noTaskMsg = "You have no completed tasks."
	} else {
		tasks, err = s.GetTasksDoneToday()
		taskMsg = "You have completed following tasks today:"
		noTaskMsg = "You have no completed tasks today."
	}
	if err != nil {
		fmt.Println("Error fetching tasks from db:", err)
		os.Exit(1)
	}

	if len(tasks) == 0 {
		fmt.Println(noTaskMsg)
		return
	}

	fmt.Println(taskMsg)

	for i, v := range tasks {
		fmt.Printf("%d. %s\n", i+1, v.Value)
	}
}

func init() {
	doneCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all completed tasks")
	rootCmd.AddCommand(doneCmd)
}
