package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:     "complete TODOID",
	Short:   "Mark a todo as complete",
	Aliases: []string{"do"},
	Args:    cobra.ExactArgs(1),
	RunE:    completeFunc,
}

func init() {
	rootCmd.AddCommand(completeCmd)
}

func completeFunc(cmd *cobra.Command, args []string) error {
	var err error
	todoManager := getManager()

	todoNum := args[0]
	todoID, err := strconv.Atoi(todoNum)
	if err != nil {
		return err
	}

	err = todoManager.Complete(todoID)

	if err != nil {
		return err
	}

	fmt.Printf("Completed Todo ID %d\n", todoID)

	return nil
}
