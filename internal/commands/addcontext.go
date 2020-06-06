package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var addContextCmd = &cobra.Command{
	Use:   "addcontext [TODO ID] [CONTEXT]",
	Short: "Add a new context to a todo",
	Args:  cobra.ExactArgs(2),
	RunE:  addContextFunc,
}

func init() {
	rootCmd.AddCommand(addContextCmd)
}

func addContextFunc(cmd *cobra.Command, args []string) error {
	var err error
	todoManager := getManager()

	todoNum := args[0]
	todoID, err := strconv.Atoi(todoNum)
	if err != nil {
		return err
	}

	context := args[1]
	err = todoManager.AddContext(todoID, context)

	if err != nil {
		return err
	}

	fmt.Printf("Added context \"%s\" to Todo ID %d\n", context, todoID)

	return nil
}
