package gotodo

// ByCreatedDate provides sorting by Todo.CreationDate
type ByCreatedDate TodoList

// Len returns length of the slice
func (s ByCreatedDate) Len() int {
	return len(s)
}

// Swap inverts positions of two elements
func (s ByCreatedDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less compares two elements by CreationDate
func (s ByCreatedDate) Less(i, j int) bool {
	t1 := s[i].CreationDate
	t2 := s[j].CreationDate

	// If we have an invalid time, prioritize valid times
	if !t1.Valid || !t2.Valid {
		if !t1.Valid && t2.Valid {
			return false
		}

		return true
	} else if t1.Time.Equal(t2.Time) {
		return s[i].TodoID < s[j].TodoID
	}

	return t1.Time.Before(t2.Time)
}

// ByDueDate provides sorting by Todo.DueDate
type ByDueDate TodoList

// Len returns length of the slice
func (s ByDueDate) Len() int {
	return len(s)
}

// Swap inverts positions of two elements
func (s ByDueDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less compares two elements by DueDate
func (s ByDueDate) Less(i, j int) bool {
	t1 := s[i].DueDate
	t2 := s[j].DueDate

	// If we have an invalid time, prioritize valid times
	if !t1.Valid || !t2.Valid {
		if !t1.Valid && t2.Valid {
			return false
		}

		return true
	} else if t1.Time.Equal(t2.Time) {
		return s[i].TodoID < s[j].TodoID
	}

	return t1.Time.Before(t2.Time)
}

// ByPriority provides sorting by Todo.Priority
type ByPriority TodoList

// Len returns length of the slice
func (s ByPriority) Len() int {
	return len(s)
}

// Swap inverts positions of two elements
func (s ByPriority) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less compares two elements by Priority
func (s ByPriority) Less(i, j int) bool {
	// 0 priority means the todo is nor prioritized.
	// Consider 0 priority higher than any other value
	if s[i].Priority == 0 && s[j].Priority > 0 {
		return false
	} else if s[i].Priority > 0 && s[j].Priority == 0 {
		return true
	} else if s[i].Priority == 0 && s[j].Priority == 0 {
		return s[i].TodoID < s[j].TodoID
	}

	return s[i].Priority < s[j].Priority
}
