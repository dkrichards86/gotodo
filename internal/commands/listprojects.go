package commands

import (
	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "listprojects",
	Short: "list your projects",
	RunE:  projectsFunc,
}

func init() {
	rootCmd.AddCommand(projectsCmd)
}

func projectsFunc(cmd *cobra.Command, args []string) error {
	var err error
	todoManager := getManager()

	items, err := todoManager.ListProjects()
	if err != nil {
		return err
	}

	header := []string{"Projects"}
	data := make([][]string, len(items))
	for i, tag := range items {
		data[i] = []string{tag}
	}

	drawTable(header, data)

	return nil
}
