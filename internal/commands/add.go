package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [todo text]",
	Short: "create a new todo",
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

	data := make([][]string, 1)
	data[0] = []string{fmt.Sprintf("Addded Todo ID %d", todoID)}
	drawTable([]string{}, data)

	return nil
}
