package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:     "resume [todo id]",
	Short:   "marks a todo as incomplete",
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

	data := make([][]string, 1)
	data[0] = []string{fmt.Sprintf("Resumed Todo ID %d", todoID)}
	drawTable([]string{}, data)

	return nil
}
