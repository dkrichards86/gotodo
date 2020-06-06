package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var contextsCmd = &cobra.Command{
	Use:   "listcontexts",
	Short: "List your contexts",
	RunE:  contextsFunc,
}

func init() {
	rootCmd.AddCommand(contextsCmd)
}

func contextsFunc(cmd *cobra.Command, args []string) error {
	var err error
	todoManager := getManager()

	items, err := todoManager.ListContexts()
	if err != nil {
		return err
	}

	if len(items) == 0 {
		fmt.Println("No contexts to display.")
		return nil
	}

	header := []string{"Contexts"}
	data := make([][]string, len(items))
	for i, tag := range items {
		data[i] = []string{tag}
	}

	drawTable(header, data)

	return nil
}
