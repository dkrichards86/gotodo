package commands

import (
	"errors"
	"fmt"
	"sort"

	"github.com/dkrichards86/gotodo/internal/gotodo"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:     "list",
	Short:   "list your todos",
	Aliases: []string{"ls"},
	RunE:    lsFunc,
}

func init() {
	rootCmd.AddCommand(lsCmd)

	lsCmd.Flags().Bool("done", false, "only show completed todos")
	lsCmd.Flags().Bool("all", false, "show pending and completed todos")

	lsCmd.Flags().String("sort", "pending", "sort todos")
	lsCmd.Flags().String("project", "", "filter todos by project")
	lsCmd.Flags().String("context", "", "filter todos by context")
	lsCmd.Flags().String("attribute", "", "filter todos by attribute")
}

func lsFunc(cmd *cobra.Command, args []string) error {
	var err error
	todoManager := getManager()

	allFlag, err := cmd.Flags().GetBool("all")
	if err != nil {
		return err
	}
	doneFlag, err := cmd.Flags().GetBool("done")
	if err != nil {
		return err
	}

	if allFlag && doneFlag {
		return errors.New("Can't filter by both done and all status")
	}

	status := gotodo.ListPending
	if allFlag {
		status = gotodo.ListAll
	} else if doneFlag {
		status = gotodo.ListDone
	}

	projectFlag, err := cmd.Flags().GetString("project")
	if err != nil {
		return err
	}
	contextFlag, err := cmd.Flags().GetString("context")
	if err != nil {
		return err
	}
	attributeFlag, err := cmd.Flags().GetString("attribute")
	if err != nil {
		return err
	}

	listFilter := gotodo.TodoListFilter{
		Status:    status,
		Project:   projectFlag,
		Context:   contextFlag,
		Attribute: attributeFlag,
	}

	items, err := todoManager.List(listFilter)
	if err != nil {
		return err
	}

	sortFlag, err := cmd.Flags().GetString("sort")
	if err != nil {
		return err
	}

	switch sortFlag {
	case "created":
		sort.Sort(gotodo.ByCreatedDate(items))
	case "due":
		sort.Sort(gotodo.ByDueDate(items))
	default:
		sort.Sort(gotodo.ByPriority(items))
	}

	header := []string{"ID", "Todo"}
	data := make([][]string, len(items))
	for i, todo := range items {
		data[i] = []string{fmt.Sprintf("%d", todo.TodoID), todo.String()}
	}
	drawTable(header, data)

	return nil
}
