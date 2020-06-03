package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var addProjectCmd = &cobra.Command{
	Use:   "addproject [todo id] [project]",
	Short: "add a new project to a todo",
	Args:  cobra.ExactArgs(2),
	RunE:  addProjectFunc,
}

func init() {
	rootCmd.AddCommand(addProjectCmd)
}

func addProjectFunc(cmd *cobra.Command, args []string) error {
	var err error
	todoManager := getManager()

	todoNum := args[0]
	todoID, err := strconv.Atoi(todoNum)
	if err != nil {
		return err
	}

	project := args[1]
	err = todoManager.AddProject(todoID, project)

	if err != nil {
		return err
	}

	data := make([][]string, 1)
	data[0] = []string{fmt.Sprintf("Added project \"%s\" to Todo ID %d", project, todoID)}
	drawTable([]string{}, data)

	return nil
}
