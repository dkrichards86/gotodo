package gotodo

// TimeFormat is the YYYY-MM-DD format used by todo.txt
const TimeFormat = "2006-01-02"

// ListPending is all active todos
// ListAll is all todos, active and complete
// ListDone is all completed todos
const (
	ListPending = iota
	ListAll
	ListDone
)
