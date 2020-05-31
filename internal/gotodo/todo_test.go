package gotodo

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	now := time.Now()
	displayTime := now.Format(TimeFormat)

	var todo Todo
	var expected string

	todo = Todo{
		Description:  "foo",
		CreationDate: ValidTime(now),
	}
	expected = fmt.Sprintf("%s foo", displayTime)
	assert.Equal(t, expected, todo.String())

	todo = Todo{
		Description:  "foo",
		CreationDate: ValidTime(now),
		Priority:     1,
	}
	expected = fmt.Sprintf("(A) %s foo", displayTime)
	assert.Equal(t, expected, todo.String())

	todo = Todo{
		Description:    "foo",
		Complete:       true,
		CreationDate:   ValidTime(now),
		CompletionDate: ValidTime(now),
	}
	expected = fmt.Sprintf("x %s %s foo", displayTime, displayTime)
	assert.Equal(t, expected, todo.String())
}

func TestFromString(t *testing.T) {
	var todoStr string
	var todo *Todo
	var ok bool

	todoStr = "(B) 2020-04-28 Work on unit tests @codehealth +gotodo"
	todo = FromString(todoStr)
	assert.Equal(t, false, todo.Complete)
	assert.Equal(t, 2, todo.Priority)
	assert.Equal(t, false, todo.CompletionDate.Valid)
	assert.Equal(t, "", todo.CompletionDate.Display())
	assert.Equal(t, true, todo.CreationDate.Valid)
	assert.Equal(t, "2020-04-28", todo.CreationDate.Display())
	assert.Equal(t, "Work on unit tests @codehealth +gotodo", todo.Description)

	assert.Equal(t, 1, len(todo.Projects))
	_, ok = todo.Projects["gotodo"]
	assert.Equal(t, true, ok)

	assert.Equal(t, 1, len(todo.Contexts))
	_, ok = todo.Contexts["codehealth"]
	assert.Equal(t, true, ok)

	todoStr = "x 2020-04-29 2020-04-28 Add parser test +gotodo"
	todo = FromString(todoStr)
	assert.Equal(t, true, todo.Complete)
	assert.Equal(t, 0, todo.Priority)
	assert.Equal(t, true, todo.CompletionDate.Valid)
	assert.Equal(t, "2020-04-29", todo.CompletionDate.Display())
	assert.Equal(t, true, todo.CreationDate.Valid)
	assert.Equal(t, "2020-04-28", todo.CreationDate.Display())
	assert.Equal(t, "Add parser test +gotodo", todo.Description)

	assert.Equal(t, 1, len(todo.Projects))
	_, ok = todo.Projects["gotodo"]
	assert.Equal(t, true, ok)

	assert.Equal(t, 0, len(todo.Contexts))
}

func TestFromStringParts(t *testing.T) {
	var todoStr string
	var todo *Todo

	todoStr = "x"
	todo = FromString(todoStr)
	assert.Equal(t, false, todo.Complete)
	assert.Equal(t, 0, todo.Priority)
	assert.Equal(t, false, todo.CompletionDate.Valid)
	assert.Equal(t, "", todo.CompletionDate.Display())
	assert.Equal(t, false, todo.CreationDate.Valid)
	assert.Equal(t, "", todo.CreationDate.Display())
	assert.Equal(t, "x", todo.Description)

	todoStr = "x x"
	todo = FromString(todoStr)
	assert.Equal(t, true, todo.Complete)
	assert.Equal(t, 0, todo.Priority)
	assert.Equal(t, false, todo.CompletionDate.Valid)
	assert.Equal(t, "", todo.CompletionDate.Display())
	assert.Equal(t, false, todo.CreationDate.Valid)
	assert.Equal(t, "", todo.CreationDate.Display())
	assert.Equal(t, "x", todo.Description)

	todoStr = "2020-04-29"
	todo = FromString(todoStr)
	assert.Equal(t, false, todo.Complete)
	assert.Equal(t, 0, todo.Priority)
	assert.Equal(t, false, todo.CompletionDate.Valid)
	assert.Equal(t, "", todo.CompletionDate.Display())
	assert.Equal(t, false, todo.CreationDate.Valid)
	assert.Equal(t, "", todo.CreationDate.Display())
	assert.Equal(t, "2020-04-29", todo.Description)

	todoStr = "2020-04-29 2020-04-29"
	todo = FromString(todoStr)
	assert.Equal(t, false, todo.Complete)
	assert.Equal(t, 0, todo.Priority)
	assert.Equal(t, false, todo.CompletionDate.Valid)
	assert.Equal(t, "", todo.CompletionDate.Display())
	assert.Equal(t, true, todo.CreationDate.Valid)
	assert.Equal(t, "2020-04-29", todo.CreationDate.Display())
	assert.Equal(t, "2020-04-29", todo.Description)

	todoStr = "2020-04-29 2020-04-29 2020-04-29"
	todo = FromString(todoStr)
	assert.Equal(t, false, todo.Complete)
	assert.Equal(t, 0, todo.Priority)
	assert.Equal(t, false, todo.CompletionDate.Valid)
	assert.Equal(t, "", todo.CompletionDate.Display())
	assert.Equal(t, true, todo.CreationDate.Valid)
	assert.Equal(t, "2020-04-29", todo.CreationDate.Display())
	assert.Equal(t, "2020-04-29 2020-04-29", todo.Description)

	todoStr = " x 2020-04-29 2020-04-29 2020-04-29"
	todo = FromString(todoStr)
	assert.Equal(t, true, todo.Complete)
	assert.Equal(t, 0, todo.Priority)
	assert.Equal(t, true, todo.CompletionDate.Valid)
	assert.Equal(t, "2020-04-29", todo.CompletionDate.Display())
	assert.Equal(t, true, todo.CreationDate.Valid)
	assert.Equal(t, "2020-04-29", todo.CreationDate.Display())
	assert.Equal(t, "2020-04-29", todo.Description)
}

func TestHasProject(t *testing.T) {
	var todoStr string
	var todo *Todo

	todoStr = "Add tests for hasXXX helpers @codehealth +gotodo due:2020-06-01"
	todo = FromString(todoStr)

	assert.Equal(t, false, todo.hasProject("codehealth"))
	assert.Equal(t, true, todo.hasProject("gotodo"))
	assert.Equal(t, false, todo.hasProject("due"))
}

func TestHasContext(t *testing.T) {
	var todoStr string
	var todo *Todo

	todoStr = "Add tests for hasXXX helpers @codehealth +gotodo due:2020-06-01"
	todo = FromString(todoStr)

	assert.Equal(t, true, todo.hasContext("codehealth"))
	assert.Equal(t, false, todo.hasContext("gotodo"))
	assert.Equal(t, false, todo.hasContext("due"))
}

func TestHasAttribute(t *testing.T) {
	var todoStr string
	var todo *Todo

	todoStr = "Add tests for hasXXX helpers @codehealth +gotodo due:2020-06-01"
	todo = FromString(todoStr)

	assert.Equal(t, false, todo.hasAttribute("codehealth"))
	assert.Equal(t, false, todo.hasAttribute("gotodo"))
	assert.Equal(t, true, todo.hasAttribute("due"))
}
