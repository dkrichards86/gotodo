package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var depriCmd = &cobra.Command{
	Use:     "deprioritize TODOID",
	Short:   "Remove the priority from a todo",
	Aliases: []string{"depri"},
	Args:    cobra.ExactArgs(1),
	RunE:    depriFunc,
}

func init() {
	rootCmd.AddCommand(depriCmd)
}

func depriFunc(cmd *cobra.Command, args []string) error {
	var err error
	todoManager := getManager()

	todoNum := args[0]
	todoID, err := strconv.Atoi(todoNum)
	if err != nil {
		return err
	}

	err = todoManager.Deprioritize(todoID)

	if err != nil {
		return err
	}

	fmt.Printf("Removed priority for Todo ID %d\n", todoID)

	return nil
}
