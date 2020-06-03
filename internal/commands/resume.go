package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:     "resume [TODO ID]",
	Short:   "Mark a todo as incomplete",
	Aliases: []string{"undo"},
	Args:    cobra.ExactArgs(1),
	RunE:    resumeFunc,
}

func init() {
	rootCmd.AddCommand(resumeCmd)
}

func resumeFunc(cmd *cobra.Command, args []string) error {
	var err error
	todoManager := getManager()

	todoNum := args[0]
	todoID, err := strconv.Atoi(todoNum)
	if err != nil {
		return err
	}

	err = todoManager.Resume(todoID)

	if err != nil {
		return err
	}

	fmt.Printf("Resumed Todo ID %d\n", todoID)

	return nil
}
