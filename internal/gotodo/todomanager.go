package gotodo

import (
	"strings"
	"time"
)

// void is an empty struct for use in sets
type void struct{}

// TodoList is list of Todos
type TodoList []*Todo

// TodoManager controls a TodoList
type TodoManager struct {
	Storage Storage
}

// TodoListFilter provides filtering criteria for a TodoList
type TodoListFilter struct {
	Status    int
	Project   string
	Context   string
	Attribute string
}

// List returns a slice of Todos as determined by TodoListFilter criteria
func (me *TodoManager) List(listFilter TodoListFilter) (TodoList, error) {
	var err error

	items, err := me.Storage.List()
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

		if listFilter.Context != "" {
			if len(todo.Contexts) == 0 {
				continue
			}

			if _, ok := todo.Contexts[listFilter.Context]; !ok {
				continue
			}
		}

		if listFilter.Project != "" {
			if len(todo.Projects) == 0 {
				continue
			}

			if _, ok := todo.Projects[listFilter.Project]; !ok {
				continue
			}
		}

		if listFilter.Attribute != "" {
			if len(todo.CustomAttributes) == 0 {
				continue
			}

			if _, ok := todo.CustomAttributes[listFilter.Attribute]; !ok {
				continue
			}
		}

		itemsToDisplay = append(itemsToDisplay, todo)
	}

	return itemsToDisplay, nil
}

// Add takes a todotxt string and adds it to the list of todos
func (me *TodoManager) Add(todoStr string) (int, error) {
	todo := FromString(todoStr)
	return todo.TodoID, me.Storage.Create(todo)
}

// Update takes the ID number of an existing Todo and a parseable todo string and merges
// contents of the existing todo and the update.
func (me *TodoManager) Update(todoID int, todoStr string) error {
	todo, err := me.Storage.Get(todoID)
	if err != nil {
		return err
	}

	newTodo := FromString(todoStr)

	todo.Priority = newTodo.Priority
	todo.Description = newTodo.Description
	todo.Projects = newTodo.Projects
	todo.Contexts = newTodo.Contexts

	return me.Storage.Update(todoID, todo)
}

// Replace takes the ID number of an existing Todo and a parseable todo string and replaces all
// contents of the existing todo with the update.
func (me *TodoManager) Replace(todoID int, todoStr string) error {
	todo, err := me.Storage.Get(todoID)
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
	todo.CustomAttributes = newTodo.CustomAttributes

	return me.Storage.Update(todoID, todo)
}

// Prioritize changes the priority of a Todo identified by todoID
func (me *TodoManager) Prioritize(todoID int, priorityString string) error {
	todo, err := me.Storage.Get(todoID)
	if err != nil {
		return err
	}

	priority := parsePriority(priorityString)
	todo.Priority = priority

	return me.Storage.Update(todoID, todo)
}

// Deprioritize changes the priority of a Todo identified by todoID
func (me *TodoManager) Deprioritize(todoID int) error {
	todo, err := me.Storage.Get(todoID)
	if err != nil {
		return err
	}

	todo.Priority = 0

	return me.Storage.Update(todoID, todo)
}

// AddProject adds a project tag to a todo
func (me *TodoManager) AddProject(todoID int, project string) error {
	todo, err := me.Storage.Get(todoID)
	if err != nil {
		return err
	}

	if _, ok := todo.Projects[project]; !ok {
		todo.Projects[project] = void{}
		todo.Description = todo.Description + " +" + project
	}

	return me.Storage.Update(todoID, todo)
}

// AddContext adds a context tag to a todo
func (me *TodoManager) AddContext(todoID int, context string) error {
	todo, err := me.Storage.Get(todoID)
	if err != nil {
		return err
	}

	if _, ok := todo.Contexts[context]; !ok {
		todo.Contexts[context] = void{}
		todo.Description = todo.Description + " @" + context
	}

	return me.Storage.Update(todoID, todo)
}

// AddAttribute adds a context tag to a todo
func (me *TodoManager) AddAttribute(todoID int, attr string) error {
	todo, err := me.Storage.Get(todoID)
	if err != nil {
		return err
	}

	if strings.Contains(attr, ":") {
		parts := strings.Split(attr, ":")
		// If the attribute has multiple colons, use the first part as key and concat the rest
		// as value
		key := parts[0]
		value := strings.Join(parts[1:], ":")

		todo.CustomAttributes[key] = value
		todo.Description = todo.Description + " " + attr
	}

	return me.Storage.Update(todoID, todo)
}

// Complete changes the completion status of a Todo to done and adds CompletionDate
func (me *TodoManager) Complete(todoID int) error {
	todo, err := me.Storage.Get(todoID)
	if err != nil {
		return err
	}

	todo.Complete = true
	todo.CompletionDate = ValidTime(time.Now())
	todo.Priority = 0

	return me.Storage.Update(todoID, todo)
}

// Resume changes the completion status of a todo and invalidates CompletionDate
func (me *TodoManager) Resume(todoID int) error {
	todo, err := me.Storage.Get(todoID)
	if err != nil {
		return err
	}
	todo.Complete = false
	todo.CompletionDate = InvalidTime

	return me.Storage.Update(todoID, todo)
}

// Delete drops the item specified by todoId from a TodoManager
func (me *TodoManager) Delete(todoID int) error {
	return me.Storage.Delete(todoID)
}

// ListProjects returns a list of unique projects
func (me *TodoManager) ListProjects() ([]string, error) {
	projs := make([]string, 0)
	set := make(map[string]void)
	var elem void

	items, err := me.Storage.List()
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
func (me *TodoManager) ListContexts() ([]string, error) {
	ctxs := make([]string, 0)
	set := make(map[string]void)
	var elem void

	items, err := me.Storage.List()
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
func (me *TodoManager) ListAttributes() ([]string, error) {
	attrs := make([]string, 0)
	set := make(map[string]void)
	var elem void

	items, err := me.Storage.List()
	if err != nil {
		return attrs, err
	}

	for _, todo := range items {
		for attr := range todo.CustomAttributes {
			set[attr] = elem
		}
	}

	for key := range set {
		attrs = append(attrs, key)
	}

	return attrs, nil
}
