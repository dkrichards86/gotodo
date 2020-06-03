package commands

import (
	"github.com/spf13/cobra"
)

var contextsCmd = &cobra.Command{
	Use:   "listcontexts",
	Short: "list your contexts",
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

	header := []string{"Contexts"}
	data := make([][]string, len(items))
	for i, tag := range items {
		data[i] = []string{tag}
	}

	drawTable(header, data)

	return nil
}
