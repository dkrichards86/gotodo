package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var addAttributeCmd = &cobra.Command{
	Use:   "addattribute [todo id] [attribute:value]",
	Short: "add a new attribute to a todo",
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

	data := make([][]string, 1)
	data[0] = []string{fmt.Sprintf("Added attribute \"%s\" to Todo ID %d", attribute, todoID)}
	drawTable([]string{}, data)

	return nil
}
