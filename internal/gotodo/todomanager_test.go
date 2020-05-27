package gotodo

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getTestTodoManager() *TodoManager {
	items := make(TodoList, 0)
	todoStrs := []string{
		"(B) 2020-04-28 Work on unit tests @codehealth +gotodo",
		"x 2020-04-29 2020-04-28 Add parser test +gotodo due:2020-05-01",
	}

	for i, todoStr := range todoStrs {
		todo := FromString(todoStr)
		todo.TodoID = i + 1
		items = append(items, todo)
	}

	storage := &TestStorage{items}

	return &TodoManager{storage}
}

func sliceContains(needle string, haystack []string) bool {
	for _, elem := range haystack {
		if needle == elem {
			return true
		}
	}

	return false
}

func TestList(t *testing.T) {
	todoManager := getTestTodoManager()

	status := ListPending
	project := ""
	context := ""
	attribute := ""
	listFilter := TodoListFilter{status, project, context, attribute}

	items, _ := todoManager.List(listFilter)
	assert.Equal(t, 1, len(items))
	assert.Equal(t, "Work on unit tests @codehealth +gotodo", items[0].Description)
}

func TestListAll(t *testing.T) {
	todoManager := getTestTodoManager()

	status := ListAll
	project := ""
	context := ""
	attribute := ""
	listFilter := TodoListFilter{status, project, context, attribute}

	items, _ := todoManager.List(listFilter)
	assert.Equal(t, 2, len(items))
	assert.Equal(t, "Work on unit tests @codehealth +gotodo", items[0].Description)
	assert.Equal(t, "Add parser test +gotodo due:2020-05-01", items[1].Description)
}

func TestListDone(t *testing.T) {
	todoManager := getTestTodoManager()

	status := ListDone
	project := ""
	context := ""
	attribute := ""
	listFilter := TodoListFilter{status, project, context, attribute}

	items, _ := todoManager.List(listFilter)
	assert.Equal(t, 1, len(items))
	assert.Equal(t, "Add parser test +gotodo due:2020-05-01", items[0].Description)
}

func TestListProjectFilter(t *testing.T) {
	todoManager := getTestTodoManager()

	status := ListAll
	project := "gotodo"
	context := ""
	attribute := ""
	listFilter := TodoListFilter{status, project, context, attribute}

	items, _ := todoManager.List(listFilter)
	assert.Equal(t, 2, len(items))
	assert.Equal(t, "Work on unit tests @codehealth +gotodo", items[0].Description)
	assert.Equal(t, "Add parser test +gotodo due:2020-05-01", items[1].Description)
}

func TestListContextFilter(t *testing.T) {
	todoManager := getTestTodoManager()

	status := ListAll
	project := ""
	context := "codehealth"
	attribute := ""
	listFilter := TodoListFilter{status, project, context, attribute}

	items, _ := todoManager.List(listFilter)
	assert.Equal(t, 1, len(items), 1)
	assert.Equal(t, "Work on unit tests @codehealth +gotodo", items[0].Description)
}

func TestListAttributeFilter(t *testing.T) {
	todoManager := getTestTodoManager()

	status := ListAll
	project := ""
	context := ""
	attribute := "due"
	listFilter := TodoListFilter{status, project, context, attribute}

	items, _ := todoManager.List(listFilter)
	assert.Equal(t, 1, len(items), 1)
	assert.Equal(t, "Add parser test +gotodo due:2020-05-01", items[0].Description)
}

func TestListCombinedFilter(t *testing.T) {
	todoManager := getTestTodoManager()

	status := ListAll
	project := ""
	context := "codehealth"
	attribute := "due"
	listFilter := TodoListFilter{status, project, context, attribute}

	items, _ := todoManager.List(listFilter)
	assert.Equal(t, 0, len(items))
}

func TestAdd(t *testing.T) {
	todoManager := getTestTodoManager()
	todoStr := "(A) 2020-05-01 Mock TodoManager struct +gotodo @testing"

	items, _ := todoManager.Storage.List()
	assert.Equal(t, 2, len(items))

	_ = todoManager.Add(todoStr)

	items, _ = todoManager.Storage.List()
	assert.Equal(t, 3, len(items))
}

func TestUpdate(t *testing.T) {
	todoManager := getTestTodoManager()
	todoStr := "(A) 2020-05-01 Mock TodoManager struct +gotodo @testing"
	ts, _ := time.Parse(TimeFormat, "2020-04-28")

	items, _ := todoManager.Storage.List()
	assert.Equal(t, 2, len(items))

	todo, _ := todoManager.Storage.Get(0)
	assert.Equal(t, "Work on unit tests @codehealth +gotodo", todo.Description)
	assert.Equal(t, ValidTime(ts), todo.CreationDate)
	assert.Equal(t, 2, todo.Priority)

	todoManager.Update(0, todoStr)

	items, _ = todoManager.Storage.List()
	assert.Equal(t, 2, len(items))

	todo, _ = todoManager.Storage.Get(0)
	assert.Equal(t, "Mock TodoManager struct +gotodo @testing", todo.Description)
	assert.Equal(t, ValidTime(ts), todo.CreationDate)
	assert.Equal(t, 1, todo.Priority)
}

func TestReplace(t *testing.T) {
	todoManager := getTestTodoManager()
	todoStr := "(A) 2020-05-01 Mock TodoManager struct +gotodo @testing"
	ts, _ := time.Parse(TimeFormat, "2020-04-28")

	items, _ := todoManager.Storage.List()
	assert.Equal(t, 2, len(items))

	todo, _ := todoManager.Storage.Get(0)
	assert.Equal(t, "Work on unit tests @codehealth +gotodo", todo.Description)
	assert.Equal(t, ValidTime(ts), todo.CreationDate)
	assert.Equal(t, 2, todo.Priority)

	todoManager.Replace(0, todoStr)

	items, _ = todoManager.Storage.List()
	assert.Equal(t, 2, len(items))

	todo, _ = todoManager.Storage.Get(0)

	assert.Equal(t, "Mock TodoManager struct +gotodo @testing", todo.Description)
	ts2, _ := time.Parse(TimeFormat, "2020-05-01")
	assert.Equal(t, ValidTime(ts2), todo.CreationDate)
	assert.Equal(t, 1, todo.Priority)
}

func TestPrioritize(t *testing.T) {
	todoManager := getTestTodoManager()

	todo, _ := todoManager.Storage.Get(0)
	assert.Equal(t, 2, todo.Priority)

	todoManager.Prioritize(0, "A")

	todo, _ = todoManager.Storage.Get(0)
	assert.Equal(t, 1, todo.Priority)
}

func TestDeprioritize(t *testing.T) {
	todoManager := getTestTodoManager()

	todo, _ := todoManager.Storage.Get(0)
	assert.Equal(t, 2, todo.Priority)

	todoManager.Deprioritize(0)

	todo, _ = todoManager.Storage.Get(0)
	assert.Equal(t, 0, todo.Priority)
}

func TestAddProject(t *testing.T) {
	todoManager := getTestTodoManager()

	todo, _ := todoManager.Storage.Get(0)
	assert.Equal(t, 1, len(todo.Projects))

	todoManager.AddProject(0, "testing")

	todo, _ = todoManager.Storage.Get(0)
	assert.Equal(t, 2, len(todo.Projects))

	_, ok := todo.Projects["testing"]
	assert.Equal(t, true, ok)
	assert.Equal(t, true, strings.Contains(todo.Description, "+testing"))
}

func TestAddContext(t *testing.T) {
	todoManager := getTestTodoManager()

	todo, _ := todoManager.Storage.Get(0)
	assert.Equal(t, 1, len(todo.Contexts))

	todoManager.AddContext(0, "testing")

	todo, _ = todoManager.Storage.Get(0)
	assert.Equal(t, 2, len(todo.Contexts))

	_, ok := todo.Contexts["testing"]
	assert.Equal(t, true, ok)
	assert.Equal(t, true, strings.Contains(todo.Description, "@testing"))
}

func TestComplete(t *testing.T) {
	todoManager := getTestTodoManager()
	now := ValidTime(time.Now())

	todo, _ := todoManager.Storage.Get(0)
	assert.Equal(t, false, todo.Complete)
	assert.Equal(t, false, todo.CompletionDate.Valid)
	assert.Equal(t, "", todo.CompletionDate.Display())
	assert.Equal(t, 2, todo.Priority)

	todoManager.Complete(0)

	todo, _ = todoManager.Storage.Get(0)
	assert.Equal(t, true, todo.Complete)
	assert.Equal(t, true, todo.CompletionDate.Valid)
	assert.Equal(t, now.Display(), todo.CompletionDate.Display())
	assert.Equal(t, 0, todo.Priority)
}

func TestResume(t *testing.T) {
	todoManager := getTestTodoManager()

	todo, _ := todoManager.Storage.Get(1)
	assert.Equal(t, true, todo.Complete)
	assert.Equal(t, true, todo.CompletionDate.Valid)
	assert.Equal(t, "2020-04-29", todo.CompletionDate.Display())

	todoManager.Resume(1)

	todo, _ = todoManager.Storage.Get(1)
	assert.Equal(t, false, todo.Complete)
	assert.Equal(t, false, todo.CompletionDate.Valid)
	assert.Equal(t, "", todo.CompletionDate.Display())
}

func TestDelete(t *testing.T) {
	todoManager := getTestTodoManager()
	todoStr := "(A) 2020-05-01 Mock TodoManager struct +gotodo @testing"
	todoManager.Add(todoStr)

	items, _ := todoManager.Storage.List()
	assert.Equal(t, 3, len(items))

	t1, _ := todoManager.Storage.Get(0)
	assert.Equal(t, "Work on unit tests @codehealth +gotodo", t1.Description)
	t2, _ := todoManager.Storage.Get(1)
	assert.Equal(t, "Add parser test +gotodo due:2020-05-01", t2.Description)
	t3, _ := todoManager.Storage.Get(2)
	assert.Equal(t, "Mock TodoManager struct +gotodo @testing", t3.Description)

	todoManager.Delete(1)

	items, _ = todoManager.Storage.List()
	assert.Equal(t, 2, len(items))

	t1, _ = todoManager.Storage.Get(0)
	assert.Equal(t, t1.Description, "Work on unit tests @codehealth +gotodo")
	t2, _ = todoManager.Storage.Get(1)
	assert.Equal(t, t2.Description, "Mock TodoManager struct +gotodo @testing")
}

func TestDeleteFirst(t *testing.T) {
	todoManager := getTestTodoManager()
	items, _ := todoManager.Storage.List()
	assert.Equal(t, 2, len(items))

	t1, _ := todoManager.Storage.Get(0)
	assert.Equal(t, t1.Description, "Work on unit tests @codehealth +gotodo")
	t2, _ := todoManager.Storage.Get(1)
	assert.Equal(t, t2.Description, "Add parser test +gotodo due:2020-05-01")

	todoManager.Delete(0)

	items, _ = todoManager.Storage.List()
	assert.Equal(t, 1, len(items))

	t1, _ = todoManager.Storage.Get(0)
	assert.Equal(t, t1.Description, "Add parser test +gotodo due:2020-05-01")
}

func TestDeleteLast(t *testing.T) {
	todoManager := getTestTodoManager()
	items, _ := todoManager.Storage.List()
	assert.Equal(t, 2, len(items))

	t1, _ := todoManager.Storage.Get(0)
	assert.Equal(t, t1.Description, "Work on unit tests @codehealth +gotodo")
	t2, _ := todoManager.Storage.Get(1)
	assert.Equal(t, t2.Description, "Add parser test +gotodo due:2020-05-01")

	todoManager.Delete(1)

	items, _ = todoManager.Storage.List()
	assert.Equal(t, 1, len(items))

	t1, _ = todoManager.Storage.Get(0)
	assert.Equal(t, t1.Description, "Work on unit tests @codehealth +gotodo")
}

func TestListProjecs(t *testing.T) {
	todoManager := getTestTodoManager()
	todoStr := "(A) 2020-05-01 Mock TodoManager struct +gotodo @testing"
	todoManager.Add(todoStr)

	items, _ := todoManager.ListProjects()
	assert.Equal(t, len(items), 1)
	assert.Equal(t, sliceContains("gotodo", items), true)
}

func TestListContexts(t *testing.T) {
	todoManager := getTestTodoManager()
	todoStr := "(A) 2020-05-01 Mock TodoManager struct +gotodo @testing"
	todoManager.Add(todoStr)

	items, _ := todoManager.ListContexts()
	assert.Equal(t, len(items), 2)
	assert.Equal(t, sliceContains("codehealth", items), true)
	assert.Equal(t, sliceContains("testing", items), true)
}

func TestListAttributes(t *testing.T) {
	todoManager := getTestTodoManager()
	todoStr := "(A) 2020-05-01 Mock TodoManager struct +gotodo @testing"
	todoManager.Add(todoStr)

	items, _ := todoManager.ListAttributes()
	assert.Equal(t, len(items), 1)
	assert.Equal(t, sliceContains("due", items), true)
}
