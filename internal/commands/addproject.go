package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var addProjectCmd = &cobra.Command{
	Use:   "addproject [TODO ID] [PROJECT]",
	Short: "Add a new project to a todo",
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

	fmt.Printf("Added project \"%s\" to Todo ID %d\n", project, todoID)

	return nil
}
