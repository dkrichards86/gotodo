package gotodo

import (
	"strings"
	"time"
)

// void is an empty struct for use in sets
type void struct{}

// TodoList is list of Todos
type TodoList []*Todo

// TodoListFilter provides filtering criteria for a TodoList
type TodoListFilter struct {
	Status    int
	Project   string
	Context   string
	Attribute string
}

// TodoManager controls a TodoList
type TodoManager struct {
	Storage               Storage
	DuePrioritizationRate int
}

// TodoManagerOptions provides functional options to TodoManager
type TodoManagerOptions func(*TodoManager)

// WithBoltStorage configures a BoltStorage instance for TodoManager
func WithBoltStorage(bucket string) TodoManagerOptions {
	return func(tm *TodoManager) {
		tm.Storage = &BoltStorage{Bucket: []byte(bucket)}
	}
}

// WithDuePrioritization configures due prioritization
func WithDuePrioritization(rate int) TodoManagerOptions {
	return func(tm *TodoManager) {
		tm.DuePrioritizationRate = rate
	}
}

// NewTodoManager builds a new TodoManager instance with options
func NewTodoManager(opts ...TodoManagerOptions) *TodoManager {
	const (
		defaultDuePrioritizationRate = 0
	)

	tm := &TodoManager{
		Storage:               &BoltStorage{Bucket: []byte("Todos")},
		DuePrioritizationRate: defaultDuePrioritizationRate,
	}

	for _, opt := range opts {
		opt(tm)
	}

	return tm
}

// List returns a slice of Todos as determined by TodoListFilter criteria
func (tm *TodoManager) List(listFilter TodoListFilter) (TodoList, error) {
	var err error

	items, err := tm.Storage.List()
	if err != nil {
		return items, err
	}

	itemsToDisplay := make(TodoList, 0)

	for _, todo := range items {
		if listFilter.Status == ListPending && todo.Complete {
			continue
		} else if listFilter.Status == ListDone && !todo.Complete {
			continue
		}

		if listFilter.Project != "" && !todo.hasProject(listFilter.Project) {
			continue
		}

		if listFilter.Context != "" && !todo.hasContext(listFilter.Context) {
			continue
		}

		if listFilter.Attribute != "" && !todo.hasAttribute(listFilter.Attribute) {
			continue
		}

		itemsToDisplay = append(itemsToDisplay, todo)
	}

	return itemsToDisplay, nil
}

// Add takes a todotxt string and adds it to the list of todos
func (tm *TodoManager) Add(todoStr string) (int, error) {
	todo := FromString(todoStr)
	return todo.TodoID, tm.Storage.Create(todo)
}

// Update takes the ID number of an existing Todo and a parseable todo string and replaces all
// contents of the existing todo with the update.
func (tm *TodoManager) Update(todoID int, todoStr string) error {
	todo, err := tm.Storage.Get(todoID)
	if err != nil {
		return err
	}

	newTodo := FromString(todoStr)

	todo.Complete = newTodo.Complete
	todo.Priority = newTodo.Priority
	todo.Description = newTodo.Description
	todo.CompletionDate = newTodo.CompletionDate
	todo.CreationDate = newTodo.CreationDate
	todo.Projects = newTodo.Projects
	todo.Contexts = newTodo.Contexts
	todo.Attributes = newTodo.Attributes

	return tm.Storage.Update(todoID, todo)
}

// Prepend adds a string message to the front of a todo description
func (tm *TodoManager) Prepend(todoID int, prependStr string) error {
	todo, err := tm.Storage.Get(todoID)
	if err != nil {
		return err
	}

	todo.Description = prependStr + " " + todo.Description

	return tm.Storage.Update(todoID, todo)
}

// Append adds a string message to the end of a todo description
func (tm *TodoManager) Append(todoID int, appendStr string) error {
	todo, err := tm.Storage.Get(todoID)
	if err != nil {
		return err
	}

	todo.Description = todo.Description + " " + appendStr

	return tm.Storage.Update(todoID, todo)
}

// Prioritize changes the priority of a Todo identified by todoID
func (tm *TodoManager) Prioritize(todoID int, priorityString string) error {
	todo, err := tm.Storage.Get(todoID)
	if err != nil {
		return err
	}

	priority := parsePriority(priorityString)
	todo.Priority = priority

	return tm.Storage.Update(todoID, todo)
}

// Deprioritize changes the priority of a Todo identified by todoID
func (tm *TodoManager) Deprioritize(todoID int) error {
	todo, err := tm.Storage.Get(todoID)
	if err != nil {
		return err
	}

	todo.Priority = 0

	return tm.Storage.Update(todoID, todo)
}

// AddProject adds a project tag to a todo
func (tm *TodoManager) AddProject(todoID int, project string) error {
	todo, err := tm.Storage.Get(todoID)
	if err != nil {
		return err
	}

	if _, ok := todo.Projects[project]; !ok {
		todo.Projects[project] = void{}
		todo.Description = todo.Description + " +" + project
	}

	return tm.Storage.Update(todoID, todo)
}

// AddContext adds a context tag to a todo
func (tm *TodoManager) AddContext(todoID int, context string) error {
	todo, err := tm.Storage.Get(todoID)
	if err != nil {
		return err
	}

	if _, ok := todo.Contexts[context]; !ok {
		todo.Contexts[context] = void{}
		todo.Description = todo.Description + " @" + context
	}

	return tm.Storage.Update(todoID, todo)
}

// AddAttribute adds a context tag to a todo
func (tm *TodoManager) AddAttribute(todoID int, attr string) error {
	todo, err := tm.Storage.Get(todoID)
	if err != nil {
		return err
	}

	if strings.Contains(attr, ":") {
		parts := strings.Split(attr, ":")
		// If the attribute has multiple colons, use the first part as key and concat the rest
		// as value
		key := parts[0]
		value := strings.Join(parts[1:], ":")

		todo.Attributes[key] = value
		todo.Description = todo.Description + " " + attr
	}

	return tm.Storage.Update(todoID, todo)
}

// Complete changes the completion status of a Todo to done and adds CompletionDate
func (tm *TodoManager) Complete(todoID int) error {
	todo, err := tm.Storage.Get(todoID)
	if err != nil {
		return err
	}

	todo.Complete = true
	todo.CompletionDate = ValidTime(time.Now())
	todo.Priority = 0

	return tm.Storage.Update(todoID, todo)
}

// Resume changes the completion status of a todo and invalidates CompletionDate
func (tm *TodoManager) Resume(todoID int) error {
	todo, err := tm.Storage.Get(todoID)
	if err != nil {
		return err
	}
	todo.Complete = false
	todo.CompletionDate = InvalidTime

	return tm.Storage.Update(todoID, todo)
}

// Delete drops the item specified by todoId from a TodoManager
func (tm *TodoManager) Delete(todoID int) error {
	return tm.Storage.Delete(todoID)
}

// ListProjects returns a list of unique projects
func (tm *TodoManager) ListProjects() ([]string, error) {
	projs := make([]string, 0)
	set := make(map[string]void)
	var elem void

	items, err := tm.Storage.List()
	if err != nil {
		return projs, err
	}

	for _, todo := range items {
		for project := range todo.Projects {
			set[project] = elem
		}
	}

	for key := range set {
		projs = append(projs, key)
	}

	return projs, nil
}

// ListContexts returns a list of unique contexts
func (tm *TodoManager) ListContexts() ([]string, error) {
	ctxs := make([]string, 0)
	set := make(map[string]void)
	var elem void

	items, err := tm.Storage.List()
	if err != nil {
		return ctxs, err
	}

	for _, todo := range items {
		for context := range todo.Contexts {
			set[context] = elem
		}
	}

	for key := range set {
		ctxs = append(ctxs, key)
	}

	return ctxs, nil
}

// ListAttributes returns a list of unique attribute keys
func (tm *TodoManager) ListAttributes() ([]string, error) {
	attrs := make([]string, 0)
	set := make(map[string]void)
	var elem void

	items, err := tm.Storage.List()
	if err != nil {
		return attrs, err
	}

	for _, todo := range items {
		for attr := range todo.Attributes {
			set[attr] = elem
		}
	}

	for key := range set {
		attrs = append(attrs, key)
	}

	return attrs, nil
}
