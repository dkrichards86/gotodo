package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dkrichards86/gotodo/internal/gotodo"
	"github.com/spf13/cobra"
)

var priCmd = &cobra.Command{
	Use:     "prioritize TODOID PRIORITY",
	Short:   "Updates the priority of a todo",
	Aliases: []string{"pri"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("accepts 2 arg(s), received %d", len(args))
		}

		if gotodo.IsPriorityString(args[1]) {
			return nil
		}

		return fmt.Errorf("invalid priority value: %s", args[1])
	},
	RunE: priFunc,
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

	fmt.Printf("Updated priority for Todo ID %d to (%s)\n", todoID, strings.ToUpper(priorityArg))

	return nil
}
