package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add TODO",
	Short: "Create a new todo",
	Args:  cobra.ExactArgs(1),
	RunE:  addFunc,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func addFunc(cmd *cobra.Command, args []string) error {
	var err error
	todoManager := getManager()

	todoStr := args[0]
	todoID, err := todoManager.Add(todoStr)
	if err != nil {
		return err
	}

	fmt.Printf("Addded Todo ID %d\n", todoID)

	return nil
}
