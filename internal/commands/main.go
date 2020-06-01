package commands

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"

	cli "github.com/urfave/cli/v2"

	"github.com/dkrichards86/gotodo/internal/gotodo"
	"github.com/olekukonko/tablewriter"
)

func getManager() *gotodo.TodoManager {
	storage := &gotodo.BoltStorage{}
	return &gotodo.TodoManager{Storage: storage}
}

func drawTable(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetTablePadding("\t")
	table.AppendBulk(data)
	table.Render()
}

func lsAction(c *cli.Context) error {
	var err error
	todoManager := getManager()

	if c.Bool("all") && c.Bool("done") {
		return errors.New("Can't filter by both done and all status")
	}

	status := gotodo.ListPending
	if c.Bool("all") {
		status = gotodo.ListAll
	} else if c.Bool("done") {
		status = gotodo.ListDone
	}

	listFilter := gotodo.TodoListFilter{
		Status:    status,
		Project:   c.String("project"),
		Context:   c.String("context"),
		Attribute: c.String("attribute"),
	}

	items, err := todoManager.List(listFilter)
	if err != nil {
		return err
	}

	switch c.String("sort") {
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

func addAction(c *cli.Context) error {
	var err error
	todoManager := getManager()

	if c.NArg() != 1 {
		return errors.New("No todo message provided")
	}

	todoStr := c.Args().Get(0)
	todoID, err := todoManager.Add(todoStr)
	if err != nil {
		return err
	}

	data := make([][]string, 1)
	data[0] = []string{fmt.Sprintf("Addded Todo ID %d", todoID)}
	drawTable([]string{}, data)

	return nil
}

func editAction(c *cli.Context) error {
	var err error
	todoManager := getManager()

	if c.NArg() != 2 {
		return errors.New("Missing todo ID or message")
	}

	if c.Bool("append") && c.Bool("prepend") {
		return errors.New("Can't append and prepend at the same time")
	}

	todoNum := c.Args().Get(0)
	todoID, err := strconv.Atoi(todoNum)
	if err != nil {
		return err
	}

	todoStr := c.Args().Get(1)
	fn := todoManager.Update

	if c.Bool("append") {
		fn = todoManager.Append
	} else if c.Bool("prepend") {
		fn = todoManager.Prepend
	}

	err = fn(todoID, todoStr)

	if err != nil {
		return err
	}

	data := make([][]string, 1)
	data[0] = []string{fmt.Sprintf("Updated Todo ID %d", todoID)}
	drawTable([]string{}, data)

	return nil
}

func priAction(c *cli.Context) error {
	var err error
	todoManager := getManager()

	if c.NArg() != 2 {
		return errors.New("Missing todo ID or priority")
	}
	todoNum := c.Args().Get(0)
	todoID, err := strconv.Atoi(todoNum)
	if err != nil {
		return err
	}

	priorityArg := c.Args().Get(1)

	err = todoManager.Prioritize(todoID, priorityArg)
	if err != nil {
		return err
	}

	data := make([][]string, 1)
	data[0] = []string{fmt.Sprintf("Updated priority for Todo ID %d to %s", todoID, priorityArg)}
	drawTable([]string{}, data)

	return nil
}

func depriAction(c *cli.Context) error {
	var err error
	todoManager := getManager()

	if c.NArg() != 1 {
		return errors.New("No todo ID provided")
	}
	todoNum := c.Args().Get(0)
	todoID, err := strconv.Atoi(todoNum)
	if err != nil {
		return err
	}

	err = todoManager.Deprioritize(todoID)
	if err != nil {
		return err
	}

	data := make([][]string, 1)
	data[0] = []string{fmt.Sprintf("Removed priority for Todo ID %d", todoID)}
	drawTable([]string{}, data)

	return nil
}

func addprojectAction(c *cli.Context) error {
	var err error
	todoManager := getManager()

	if c.NArg() != 2 {
		return errors.New("Missing todo ID or project")
	}
	todoNum := c.Args().Get(0)
	todoID, err := strconv.Atoi(todoNum)
	if err != nil {
		return err
	}

	project := c.Args().Get(1)
	err = todoManager.AddProject(todoID, project)
	if err != nil {
		return err
	}

	data := make([][]string, 1)
	data[0] = []string{fmt.Sprintf("Added project \"%s\" to Todo ID %d", project, todoID)}
	drawTable([]string{}, data)

	return nil
}

func addcontextAction(c *cli.Context) error {
	var err error
	todoManager := getManager()

	if c.NArg() != 2 {
		return errors.New("Missing todo ID or context")
	}
	todoNum := c.Args().Get(0)
	todoID, err := strconv.Atoi(todoNum)
	if err != nil {
		return err
	}

	context := c.Args().Get(1)
	err = todoManager.AddContext(todoID, context)
	if err != nil {
		return err
	}

	data := make([][]string, 1)
	data[0] = []string{fmt.Sprintf("Added context \"%s\" to Todo ID %d", context, todoID)}
	drawTable([]string{}, data)

	return nil
}

func addattributeAction(c *cli.Context) error {
	var err error
	todoManager := getManager()

	if c.NArg() != 2 {
		return errors.New("Missing todo ID or attribute")
	}
	todoNum := c.Args().Get(0)
	todoID, err := strconv.Atoi(todoNum)
	if err != nil {
		return err
	}

	attr := c.Args().Get(1)
	err = todoManager.AddAttribute(todoID, attr)
	if err != nil {
		return err
	}

	data := make([][]string, 1)
	data[0] = []string{fmt.Sprintf("Added attribute \"%s\" to Todo ID %d", attr, todoID)}
	drawTable([]string{}, data)

	return nil
}

func completeAction(c *cli.Context) error {
	var err error
	todoManager := getManager()

	if c.NArg() != 1 {
		return errors.New("No todo ID provided")
	}
	todoNum := c.Args().Get(0)
	todoID, err := strconv.Atoi(todoNum)
	if err != nil {
		return err
	}

	err = todoManager.Complete(todoID)
	if err != nil {
		return err
	}

	data := make([][]string, 1)
	data[0] = []string{fmt.Sprintf("Completed Todo ID %d", todoID)}
	drawTable([]string{}, data)

	return nil
}

func resumeAction(c *cli.Context) error {
	var err error
	todoManager := getManager()

	if c.NArg() != 1 {
		return errors.New("No todo ID provided")
	}
	todoNum := c.Args().Get(0)
	todoID, err := strconv.Atoi(todoNum)
	if err != nil {
		return err
	}

	err = todoManager.Resume(todoID)
	if err != nil {
		return err
	}

	data := make([][]string, 1)
	data[0] = []string{fmt.Sprintf("Resumed Todo ID %d", todoID)}
	drawTable([]string{}, data)

	return nil
}

func rmAction(c *cli.Context) error {
	var err error
	todoManager := getManager()
	if c.NArg() != 1 {
		return errors.New("No todo ID provided")
	}
	todoNum := c.Args().Get(0)
	todoID, err := strconv.Atoi(todoNum)
	if err != nil {
		return err
	}

	err = todoManager.Delete(todoID)
	if err != nil {
		return err
	}

	data := make([][]string, 1)
	data[0] = []string{fmt.Sprintf("Removed Todo ID %d", todoID)}
	drawTable([]string{}, data)

	return nil
}

func projectsAction(c *cli.Context) error {
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

func contextsAction(c *cli.Context) error {
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

func attributesAction(c *cli.Context) error {
	var err error
	todoManager := getManager()
	items, err := todoManager.ListAttributes()
	if err != nil {
		return err
	}

	header := []string{"Attributes"}
	data := make([][]string, len(items))
	for i, tag := range items {
		data[i] = []string{tag}
	}

	drawTable(header, data)

	return nil
}

// RunCli manages execution of all CLI commands
func RunCli(args []string) {
	app := &cli.App{
		Name:  "gotodo",
		Usage: "CLI client to manage your todos",
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "Shows a lists of your todos",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "done",
						Value: false,
						Usage: "only show completed todos",
					},
					&cli.BoolFlag{
						Name:  "all",
						Value: false,
						Usage: "show pending and completed todos",
					},
					&cli.StringFlag{
						Name:  "sort",
						Value: "priority",
						Usage: "sort todos.",
					},
					&cli.StringFlag{
						Name:  "project",
						Value: "",
						Usage: "filter todos by project",
					},
					&cli.StringFlag{
						Name:  "context",
						Value: "",
						Usage: "filter todos by context",
					},
					&cli.StringFlag{
						Name:  "attribute",
						Value: "",
						Usage: "filter todos by attribute",
					},
				},
				Action: lsAction,
			},
			{
				Name:      "add",
				Usage:     "Creates a new todo",
				ArgsUsage: "[todo text]",
				Action:    addAction,
			},
			{
				Name:  "edit",
				Usage: "Edits an existing todo",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "append",
						Value: false,
						Usage: "add message to the end of a todo",
					},
					&cli.BoolFlag{
						Name:  "prepend",
						Value: false,
						Usage: "add message to the front of a todo",
					},
				},
				ArgsUsage: "[line number] [todo text]",
				Action:    editAction,
			},
			{
				Name:      "pri",
				Usage:     "Updates the priority of a todo",
				ArgsUsage: "[line number] [priority]",
				Action:    priAction,
			},
			{
				Name:      "depri",
				Usage:     "Removes the priority from a todo",
				ArgsUsage: "[line number]",
				Action:    depriAction,
			},
			{
				Name:      "addproject",
				Usage:     "Adds a new project to a todo",
				ArgsUsage: "[line number] [project]",
				Action:    addprojectAction,
			},
			{
				Name:      "addcontext",
				Usage:     "Adds a new context to a todo",
				ArgsUsage: "[line number] [context]",
				Action:    addcontextAction,
			},
			{
				Name:      "addattribute",
				Usage:     "Adds a new attribute to a todo",
				ArgsUsage: "[line number] [attribute]",
				Action:    addattributeAction,
			},
			{
				Name:      "complete",
				Aliases:   []string{"do"},
				Usage:     "Marks a todo as complete",
				ArgsUsage: "[line number]",
				Action:    completeAction,
			},
			{
				Name:      "resume",
				Usage:     "Marks a todo as incomplete",
				ArgsUsage: "[line number]",
				Action:    resumeAction,
			},
			{
				Name:      "remove",
				Aliases:   []string{"rm"},
				Usage:     "Removes a todo",
				ArgsUsage: "[line number]",
				Action:    rmAction,
			},
			{
				Name:   "projects",
				Usage:  "Shows a list of projects",
				Action: projectsAction,
			},
			{
				Name:   "contexts",
				Usage:  "Shows a list of contexts",
				Action: contextsAction,
			},
			{
				Name:   "attributes",
				Usage:  "Shows a list of custom attributes",
				Action: attributesAction,
			},
		},
	}

	err := app.Run(args)
	if err != nil {
		header := []string{"Error"}
		data := make([][]string, 1)
		data[0] = []string{fmt.Sprintf("%s", err)}
		drawTable(header, data)
	}
}
