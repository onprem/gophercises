package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List current tasks",
	Long:  "This subcommand lists all tasks which are currenlty active",
	Run:   list,
}

func list(cmd *cobra.Command, args []string) {
	fmt.Println("This is a fake \"list\" command")
}

func init() {
	rootCmd.AddCommand(listCmd)
}
