package commands

import (
	"os"

	"github.com/dkrichards86/gotodo/internal/gotodo"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var cfgFile string
var todoList string

var rootCmd = &cobra.Command{
	Use:          "gotodo",
	Short:        "A CLI client to manage your todos",
	SilenceUsage: true,
}

// Execute runs the specified command.
func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&todoList, "bucket", "todos", "todo bucket to use")
}

func getManager() *gotodo.TodoManager {
	return gotodo.NewTodoManager(
		gotodo.WithBoltStorage(todoList),
	)
}

func drawTable(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetAutoWrapText(false)
	table.SetBorder(false)
	table.SetAutoFormatHeaders(true)
	table.AppendBulk(data)
	table.Render()
}
