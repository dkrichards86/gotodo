package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove TODOID",
	Short:   "Deletes a todo",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	RunE:    removeFunc,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func removeFunc(cmd *cobra.Command, args []string) error {
	var err error
	todoManager := getManager()

	todoNum := args[0]
	todoID, err := strconv.Atoi(todoNum)
	if err != nil {
		return err
	}

	err = todoManager.Delete(todoID)

	if err != nil {
		return err
	}

	fmt.Printf("Removed Todo ID %d\n", todoID)

	return nil
}
