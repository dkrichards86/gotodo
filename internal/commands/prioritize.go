package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var priCmd = &cobra.Command{
	Use:     "prioritize [todo id] [priority]",
	Short:   "updates the priority of a todo",
	Aliases: []string{"pri"},
	Args:    cobra.ExactArgs(2),
	RunE:    priFunc,
}

func init() {
	rootCmd.AddCommand(priCmd)
}

func priFunc(cmd *cobra.Command, args []string) error {
	var err error
	todoManager := getManager()

	todoNum := args[0]
	todoID, err := strconv.Atoi(todoNum)
	if err != nil {
		return err
	}

	priorityArg := args[1]
	err = todoManager.Prioritize(todoID, priorityArg)

	if err != nil {
		return err
	}

	data := make([][]string, 1)
	data[0] = []string{fmt.Sprintf("Updated priority for Todo ID %d to %s", todoID, priorityArg)}
	drawTable([]string{}, data)

	return nil
}
