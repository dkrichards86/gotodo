package gotodo

import (
	"fmt"
	"strings"
)

// Tags is a key/value store for custom Todo metadata
type Tags map[string]void

// Attributes is a key/value store for custom Todo metadata
type Attributes map[string]string

// Todo contains information about a specific Todo
type Todo struct {
	TodoID         int
	Complete       bool
	Priority       int
	CompletionDate NullTime
	CreationDate   NullTime
	DueDate        NullTime
	Description    string
	Projects       Tags
	Contexts       Tags
	Attributes     Attributes
}

// FromString translates a todotxt string into a *Todo
func FromString(todoStr string) *Todo {
	parts := strings.Split(strings.TrimSpace(todoStr), " ")
	var complete bool
	var priority int
	var completionDate NullTime
	var creationDate NullTime
	var dueDate NullTime
	var description string

	// If we have more than 1 part and the first part is x, consider the Todo complete.
	complete = false
	if len(parts) > 1 && isCompleteToken(parts[0]) {
		complete = true
		parts = parts[1:]
	}

	// If we have more than 1 part and the first part is a valid priority flag, translate
	// the priority.
	priority = 0
	if len(parts) > 1 && isPriorityToken(parts[0]) {
		end := len(parts[0])
		arg := parts[0][1 : end-1]
		priority = parsePriority(arg)
		parts = parts[1:]
	}

	// Check for zero, one or two times.
	// If there are zero, we move along.
	// If there is one time, it is creation time.
	// If there are two times and the todo is complete, the first is completed
	// time and the second is created.
	if len(parts) > 1 {
		firstTime := parseDate(parts[0])
		if firstTime.Valid {
			if len(parts) > 2 {
				secondTime := parseDate(parts[1])
				if secondTime.Valid && complete {
					completionDate = firstTime
					creationDate = secondTime
					parts = parts[2:]
				} else if complete {
					completionDate = firstTime
					parts = parts[1:]
				} else {
					creationDate = firstTime
					parts = parts[1:]
				}
			} else {
				creationDate = firstTime
				parts = parts[1:]
			}
		}
	}

	projects, contexts := parseTags(parts)
	customAttrs := make(Attributes)

	for i := range parts {
		idx := len(parts) - i - 1
		word := parts[idx]
		if strings.Contains(word, ":") {
			attr := strings.Split(word, ":")
			// If the attribute has multiple colons, use the first part as key and concat the rest
			// as value
			key := attr[0]
			value := strings.Join(attr[1:], ":")
			customAttrs[key] = value
		} else {
			break
		}
	}

	// Due date is a special attribute in todo.txt. It's not part of the official spec, but has
	// gained enough traction in the community that it gets special attention.
	if dueAttr, ok := customAttrs["due"]; ok {
		dueDate = parseDate(dueAttr)
	}

	description = strings.Join(parts, " ")

	return &Todo{
		Complete:       complete,
		Priority:       priority,
		CompletionDate: completionDate,
		CreationDate:   creationDate,
		DueDate:        dueDate,
		Description:    description,
		Projects:       projects,
		Contexts:       contexts,
		Attributes:     customAttrs,
	}
}

// String converts a Todo into a todotxt-formatted string
func (t *Todo) String() string {
	parts := make([]string, 0)

	if t.Complete {
		parts = append(parts, "x")
	}

	if t.Priority > 0 {
		priStr := fmt.Sprintf("(%s)", unparsePriority(t.Priority))
		parts = append(parts, priStr)
	}

	if t.Complete && t.CompletionDate.Valid {
		parts = append(parts, t.CompletionDate.Display())
	}

	if t.CreationDate.Valid {
		parts = append(parts, t.CreationDate.Display())
	}

	parts = append(parts, t.Description)

	return strings.Join(parts, " ")
}

// hasProject checks a todo.Projects for a specific project
func (t *Todo) hasProject(project string) bool {
	if len(t.Projects) == 0 {
		return false
	}

	_, exists := t.Projects[project]
	return exists
}

// hasContext checks a todo.Contexts for a specific context
func (t *Todo) hasContext(context string) bool {
	if len(t.Contexts) == 0 {
		return false
	}

	_, exists := t.Contexts[context]
	return exists
}

// hasAttribute checks a todo.Attributes for a specific attribute
func (t *Todo) hasAttribute(attribute string) bool {
	if len(t.Attributes) == 0 {
		return false
	}

	_, exists := t.Attributes[attribute]
	return exists
}
