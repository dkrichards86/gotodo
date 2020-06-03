package commands

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [TODO ID] [TODO]",
	Short: "Edit a new todo",
	Args:  cobra.ExactArgs(2),
	RunE:  editFunc,
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().Bool("append", false, "add message to the end of a tod")
	editCmd.Flags().Bool("prepend", false, "add message to the front of a todo")
}

func editFunc(cmd *cobra.Command, args []string) error {
	var err error
	todoManager := getManager()

	appendFlag, err := cmd.Flags().GetBool("append")
	if err != nil {
		return err
	}
	prependFlag, err := cmd.Flags().GetBool("prepend")
	if err != nil {
		return err
	}
	if appendFlag && prependFlag {
		return errors.New("Can't append and prepend at the same time")
	}

	todoNum := args[0]
	todoID, err := strconv.Atoi(todoNum)
	if err != nil {
		return err
	}

	fn := todoManager.Update
	if appendFlag {
		fn = todoManager.Append
	} else if prependFlag {
		fn = todoManager.Prepend
	}

	todoStr := args[1]
	err = fn(todoID, todoStr)

	if err != nil {
		return err
	}

	fmt.Printf("Updated Todo ID %d\n", todoID)

	return nil
}
