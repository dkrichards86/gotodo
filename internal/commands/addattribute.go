package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var addAttributeCmd = &cobra.Command{
	Use:   "addattribute [TODO ID] [KEY:VALUE]",
	Short: "Add a new attribute to a todo",
	Args:  cobra.ExactArgs(2),
	RunE:  addAttributeFunc,
}

func init() {
	rootCmd.AddCommand(addAttributeCmd)
}

func addAttributeFunc(cmd *cobra.Command, args []string) error {
	var err error
	todoManager := getManager()

	todoNum := args[0]
	todoID, err := strconv.Atoi(todoNum)
	if err != nil {
		return err
	}

	attribute := args[1]
	err = todoManager.AddAttribute(todoID, attribute)

	if err != nil {
		return err
	}

	fmt.Printf("Added attribute \"%s\" to Todo ID %d\n", attribute, todoID)

	return nil
}
